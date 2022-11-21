package farm

import (
	"sort"

	s "github.com/SimonRichardson/colossus/selectors"
)

type TupleSet map[s.KeyFieldScoreTxnValue]struct{}

func MakeSet(a []s.KeyFieldScoreTxnValue) TupleSet {
	m := make(TupleSet, len(a))
	for _, tuple := range a {
		m.Add(tuple)
	}
	return m
}

func (m TupleSet) Add(tuple s.KeyFieldScoreTxnValue) {
	m[tuple] = struct{}{}
}

func (m TupleSet) Has(tuple s.KeyFieldScoreTxnValue) bool {
	_, ok := m[tuple]
	return ok
}

func (m TupleSet) Slice() []s.KeyFieldScoreTxnValue {
	a := make([]s.KeyFieldScoreTxnValue, 0, len(m))
	for tuple := range m {
		a = append(a, tuple)
	}
	return a
}

func (m TupleSet) OrderedLimitedSlice(limit int) []s.KeyFieldScoreTxnValue {
	a := m.Slice()
	sort.Sort(KeyFieldScoreTxnValues(a))
	if len(a) > limit {
		a = a[:limit]
	}
	return a
}

type KeyFieldScoreTxnValues []s.KeyFieldScoreTxnValue

func (a KeyFieldScoreTxnValues) Len() int           { return len(a) }
func (a KeyFieldScoreTxnValues) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a KeyFieldScoreTxnValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type KeyFieldTxnValueSet map[s.KeyFieldTxnValue]struct{}

func (m KeyFieldTxnValueSet) Add(keyFieldTxnValue s.KeyFieldTxnValue) {
	m[keyFieldTxnValue] = struct{}{}
}

func (m KeyFieldTxnValueSet) AddMany(other KeyFieldTxnValueSet) {
	for keyFieldTxnValue := range other {
		m.Add(keyFieldTxnValue)
	}
}

func (m KeyFieldTxnValueSet) Slice() []s.KeyFieldTxnValue {
	a := make([]s.KeyFieldTxnValue, 0, len(m))
	for keyFieldTxnValue := range m {
		a = append(a, keyFieldTxnValue)
	}
	return a
}

func UnionDifference(sets []TupleSet) (TupleSet, KeyFieldTxnValueSet) {
	var (
		expectedCount = len(sets)
		scores        = map[s.KeyFieldTxnValue]float64{}
		counts        = map[s.KeyFieldScoreTxnValue]int{}
	)

	for _, set := range sets {
		for tuple := range set {
			// union
			keyFieldTxnValue := s.KeyFieldTxnValue{
				Key:   tuple.Key,
				Field: tuple.Field,
				Txn:   tuple.Txn,
				Value: tuple.Value,
			}
			if score, ok := scores[keyFieldTxnValue]; !ok || tuple.Score > score {
				scores[keyFieldTxnValue] = tuple.Score
			}

			// difference
			counts[tuple]++
		}
	}

	var (
		union      = make(TupleSet, len(scores))
		difference = make(KeyFieldTxnValueSet, len(counts))
	)

	for keyFieldTxnValue, bestScore := range scores {
		union.Add(s.KeyFieldScoreTxnValue{
			Key:   keyFieldTxnValue.Key,
			Field: keyFieldTxnValue.Field,
			Score: bestScore,
			Txn:   keyFieldTxnValue.Txn,
			Value: keyFieldTxnValue.Value,
		})
	}

	for keyFieldScoreTxnValue, count := range counts {
		if count < expectedCount {
			difference.Add(s.KeyFieldTxnValue{
				Key:   keyFieldScoreTxnValue.Key,
				Field: keyFieldScoreTxnValue.Field,
				Value: keyFieldScoreTxnValue.Value,
			})
		}
	}

	return union, difference
}
