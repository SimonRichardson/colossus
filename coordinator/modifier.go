package coordinator

import (
	"strings"
	"sync"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm/persistence"
	"github.com/SimonRichardson/colossus/farm/store"
	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/schemas/schema"
	"github.com/SimonRichardson/colossus/selectors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// Match and Replace define possible operations for patching queries
const (
	Match   selectors.Op = "match"
	Replace selectors.Op = "replace"
)

type modifier struct {
	s.LifeCycleManager

	co          *Coordinator
	store       *store.Farm
	persistence *persistence.Farm
	accessor    s.Accessor
}

func newModifier(co *Coordinator,
	s *store.Farm,
	p *persistence.Farm,
	a s.Accessor,
) *modifier {
	return &modifier{
		LifeCycleManager: newLifeCycleService(),

		co:          co,
		store:       s,
		persistence: p,
		accessor:    a,
	}
}

func (m *modifier) Modify(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	var (
		buckets      = s.KeyFieldScoreTxnValues(members).Bucketize()
		changes, err = persist(m.persistence, buckets, maxSize)
	)
	if err != nil {
		return changes, err
	}

	instr := m.co.instrumentation
	if err := removeFromStore(instr, m.store, buckets, maxSize); err != nil {
		return 0, err
	}
	return changes, nil
}

type tuple struct {
	total, changes int
	err            error
}

func newTuple(total, changes int, err error) tuple {
	return tuple{total, changes, err}
}

func persist(persistence *persistence.Farm,
	buckets map[bs.Key][]s.KeyFieldScoreTxnValue,
	maxSize s.KeySizeExpiry,
) (int, error) {
	var (
		numOfBuckets = len(buckets)
		resp         = make(chan tuple, numOfBuckets)
		wg           = &sync.WaitGroup{}
	)

	wg.Add(numOfBuckets)
	go func() { wg.Wait(); close(resp) }()

	for _, v := range buckets {
		go func(elements []s.KeyFieldScoreTxnValue) {
			defer wg.Done()

			changes, err := persistence.Insert(elements, maxSize)
			resp <- newTuple(len(elements), changes, err)
		}(v)
	}

	repair, changes, errs := repairsRequired(resp)
	if repair {
		var errs []error
		for _, members := range buckets {
			values := s.KeyFieldScoreTxnValues(members).KeyFieldTxnValues()
			if err := persistence.Repair(values, maxSize); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return 0, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Error Repairing (%s)", common.SumErrors(errs).Error())
		}
	}

	if len(errs) > 0 {
		return 0, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Error Inserting (%s)", common.SumErrors(errs).Error())
	}

	return changes, nil
}

func removeFromStore(instr instrumentation.Instrumentation,
	store *store.Farm,
	buckets map[bs.Key][]s.KeyFieldScoreTxnValue,
	maxSize s.KeySizeExpiry,
) error {
	var (
		numOfBuckets = len(buckets)
		resp         = make(chan tuple, numOfBuckets)
		wg           = &sync.WaitGroup{}
	)

	wg.Add(numOfBuckets)
	go func() { wg.Wait(); close(resp) }()

	for _, v := range buckets {
		go func(elements []s.KeyFieldScoreTxnValue) {
			defer wg.Done()

			changes, err := store.Delete(elements, maxSize)
			resp <- newTuple(len(elements), changes, err)
		}(v)
	}

	repair, _, repairErrs := repairsRequired(resp)
	if len(repairErrs) > 0 {
		return typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Error Deleting (%s)", common.SumErrors(repairErrs).Error())
	}

	if repair {
		go func() {
			for _, members := range buckets {
				values := s.KeyFieldScoreTxnValues(members).KeyFieldTxnValues()
				if err := store.Repair(values, maxSize); err != nil {
					instr.RepairError(1)
				}
			}
		}()
	}

	return nil
}

func repairsRequired(values <-chan tuple) (bool, int, []error) {
	var (
		repair  bool
		errors  []error
		changes int
	)
	for t := range values {
		if t.err != nil {
			errors = append(errors, t.err)
			repair = true
			continue
		}

		if t.changes < t.total {
			repair = true
		}

		changes += t.changes
	}

	return repair, changes, errors
}

func (m *modifier) ModifyWithOperations(key, id bs.Key,
	operations []s.Operation,
	score float64,
	maxSize s.SizeExpiry,
) (int, error) {
	targets, err := getTargets(operations)
	if err != nil {
		return -1, err
	}

	res, err := m.store.Select(key, id)
	if err != nil {
		return -1, err
	}

	if header, err := records.ReadType(res.Value); err != nil {
		return -1, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Unable to read header type")
	} else if header != schema.TypePost {
		return -1, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Unable to check type.")
	}

	var (
		num        = len(operations)
		accessor   = m.accessor
		record     = &records.PostRecord{}
		body, err1 = records.ReadBody(res.Value)
	)
	if err1 != nil {
		return -1, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Unable to read value.")
	}
	if err := record.Read(body); err != nil {
		return -1, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Unable to parse record.")
	}

	fb := pool.Get()
	defer pool.Put(fb)

loop:
	for k, op := range operations {
		var (
			t = targets[k]
			a = strings.Join(t.accessors, ".")
		)
		switch op.Op {
		case Match:
			if res, err := accessor.GetFieldValue(record, a); err == nil && res == op.Value {
				if k == num-1 {
					return k, nil
				}
				continue loop
			} else {
				return k, typex.Errorf(errors.Source, errors.UnexpectedResults,
					"Nothing matched for match.")
			}
		case Replace:
			if _, err := accessor.GetFieldValue(record, a); err == nil {
				if err := accessor.SetFieldValue(record, a, op.Value); err != nil {
					return k, err
				}

				fb.Reset()

				value, err := record.Write(fb)
				if err != nil {
					return k, err
				}

				// Insert should cause an over write of existing data
				if _, err := m.store.Insert([]s.KeyFieldScoreTxnValue{
					s.KeyFieldScoreTxnValue{
						Key:   key,
						Field: id,
						Score: score,
						Txn:   bs.Key(record.TransactionId.Hex()),
						Value: records.PackagePostRecord(value),
					},
				}, s.KeySizeExpiry{
					key: maxSize,
				}); err != nil {
					return k, err
				}

				if k == num-1 {
					return k, nil
				}
				continue loop
			} else {
				return k, typex.Errorf(errors.Source, errors.UnexpectedResults,
					"Nothing matched for replace.")
			}
		default:
			return k, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Missing implementation.")
		}
	}
	return -1, typex.Errorf(errors.Source, errors.UnexpectedResults,
		"Nothing matched.")
}

type target struct {
	accessors []string
	path      selectors.Path
}

func getTargets(ops []selectors.Operation) ([]target, error) {
	targets := make([]target, 0, len(ops))
	for k, v := range ops {
		t, err := parsePath(v.Path)
		if err != nil {
			return nil, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Error at index %d with: %s", k, err.Error())
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func parsePath(p selectors.Path) (target, error) {
	parts := p.Parts()
	if len(parts) < 2 {
		return target{}, typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid path")
	}

	return target{
		accessors: parts[1:],
		path:      p,
	}, nil
}
