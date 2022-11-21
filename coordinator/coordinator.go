package coordinator

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
	b "github.com/SimonRichardson/colossus/blist/selectors"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/consul"
	"github.com/SimonRichardson/colossus/coordinator/strategies"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	c "github.com/SimonRichardson/colossus/farm/counter"
	n "github.com/SimonRichardson/colossus/farm/notifier"
	p "github.com/SimonRichardson/colossus/farm/persistence"
	r "github.com/SimonRichardson/colossus/farm/store"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	defaultDebugExeceptions = false

	defaultInsertChannel = s.Channel("insert")

	defaultQuitTicker  = time.Millisecond * 10
	defaultQuitTimeout = time.Second * 30
)

var (
	ErrPartialInsertionFailure = typex.Errorf(errors.Source, errors.Partial, "Partial Insertion Failure")
	ErrPartialDeletionFailure  = typex.Errorf(errors.Source, errors.Partial, "Partial Deletion Failure")
	ErrPartialRollbackFailure  = typex.Errorf(errors.Source, errors.Partial, "Partial Rollback Failure")
)

// Requester defines a standardised way of requesting for items found within the
// data store.
type Requester interface {
	Request(elements map[string][]string,
		amount int,
		info map[string]interface{},
	) (map[string][]interface{}, error)
}

// Coordinator defines a single point for accessing the data store.
type Coordinator struct {
	mutex *sync.Mutex
	cond  *sync.Cond

	paused  bool
	running bool

	consul *consul.Service

	counter     *c.Farm
	store       *r.Farm
	persistence *p.Farm
	notifier    *n.Farm

	storeOpts *r.Options

	selector  s.Selector
	inserter  s.Inserter
	modifier  s.Modifier
	deleter   s.Deleter
	repairer  s.Repairer
	scanner   s.Scanner
	inspector s.Inspector
	manager   s.Manager
	service   s.Manager

	accessor    s.Accessor
	transformer s.Transformer

	managers []s.LifeCycleManager

	instrumentation instrumentation.Instrumentation
	alertmanager    alertmanager.AlertManager
}

// New defines a way to create a new Coordinator. It knits all the farms
// together to make sure the Coordinator can work, if something is missing the
// Coordinator will log out the error then exit.
func New(e *env.Env, transformer s.Transformer, accessor s.Accessor) *Coordinator {
	var (
		instr instrumentation.Instrumentation
		alert alertmanager.AlertManager

		err error
	)

	if instr, err = newInstrumentation(e, teleprinter.L.Instr()); err != nil {
		typex.Fatal(err)
	}

	if alert, err = newAlertManager(e); err != nil {
		typex.Fatal(err)
	}

	mutex := &sync.Mutex{}

	co := &Coordinator{
		mutex: mutex,
		cond:  sync.NewCond(mutex),

		paused:  false,
		running: false,

		instrumentation: instr,
		alertmanager:    alert,

		accessor:    accessor,
		transformer: transformer,
	}

	if err := co.init(e); err != nil {
		typex.Fatal(err)
	}

	return co
}

func (co *Coordinator) init(e *env.Env) error {
	var (
		consul *consul.Service

		counter     *c.Farm
		store       *r.Farm
		persistence *p.Farm
		notifier    *n.Farm

		storeOpts *r.Options

		insertStrategy  strategies.InsertStrategy
		repairStrategy  strategies.RepairStrategy
		managerStrategy strategies.ManagerStrategyCreator

		err error
	)

	if consul, err = newConsulService(e, co.instrumentation); err != nil {
		return err
	}

	if counter, err = newCounterFarm(e, co.instrumentation); err != nil {
		return err
	}

	storeOpts = newStoreOptions(e, consul)

	if store, err = newStoreFarm(e, co.instrumentation, storeOpts); err != nil {
		return err
	}

	if notifier, err = newNotifierFarm(e, co.instrumentation); err != nil {
		return err
	}

	if persistence, err = newPersistenceFarm(e, co.instrumentation, co.transformer); err != nil {
		return err
	}

	if insertStrategy, err = strategies.NewInsertStrategy(e); err != nil {
		return err
	}

	if repairStrategy, err = strategies.NewRepairStrategy(e); err != nil {
		return err
	}

	if managerStrategy, err = strategies.NewManagerStrategy(e); err != nil {
		return err
	}

	co.consul = consul

	co.counter = counter
	co.store = store
	co.persistence = persistence
	co.notifier = notifier

	co.storeOpts = storeOpts

	var (
		selector  = newSelector(co, store)
		inserter  = newInserter(co, counter, store, notifier, insertStrategy)
		modifier  = newModifier(co, store, persistence, co.accessor)
		deleter   = newDeleter(co, counter, store)
		repairer  = newRepairer(co, store, repairStrategy)
		scanner   = newScanner(co, counter)
		inspector = newInspector(co, store, co.transformer)

		manager = newManager(co, counter, store, notifier, managerStrategy)

		service = newService(co, consul, e.ConsulHeartbeatFrequency)
	)

	co.selector = selector
	co.inserter = inserter
	co.modifier = modifier
	co.deleter = deleter
	co.repairer = repairer
	co.scanner = counter
	co.inspector = inspector

	co.manager = manager
	go co.manager.Start()

	co.service = service
	go co.service.Start()

	co.managers = []s.LifeCycleManager{
		selector,
		inserter,
		modifier,
		deleter,
		repairer,
		scanner,
		inspector,
	}

	co.paused = false
	co.running = true

	return nil
}

func handle(co *Coordinator, cycle interface{}, f func()) (err error) {
	// Make sure we check what it is before we action it!
	if cyc, ok := cycle.(s.LifeCycleManager); ok {
		cyc.In()
		defer cyc.Out()
	}

	defer func() {
		switch e := recover().(type) {
		case nil:
			return
		case error:
			if defaultDebugExeceptions {
				typex.PrintStack(false)
			}

			err = e
		default:
			co.alertmanager.CoordinatorPanic()

			typex.PrintStack(false)
			panic(e)
		}
	}()

	co.mutex.Lock()
	if co.paused {
		co.cond.Wait()
	}
	co.mutex.Unlock()

	f()

	return
}

// Insert represents a way to insert various values into the store.
func (co *Coordinator) Insert(values []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (res int, err error) {
	if e := handle(co, co.inserter, func() {
		began := time.Now()
		go co.instrumentation.AInsertCall()
		defer func() { go co.instrumentation.AInsertDuration(time.Since(began)) }()

		res, err = co.inserter.Insert(values, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Modify represents a way to modify various values into the store.
func (co *Coordinator) Modify(values []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (res int, err error) {
	if e := handle(co, co.modifier, func() {
		began := time.Now()
		go co.instrumentation.AModifyCall()
		defer func() { go co.instrumentation.AModifyDuration(time.Since(began)) }()

		res, err = co.modifier.Modify(values, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// ModifyWithOperations represents a way to modify various values into the
// store.
func (co *Coordinator) ModifyWithOperations(key, id bs.Key, ops []s.Operation, score float64, maxSize s.SizeExpiry) (res int, err error) {
	if e := handle(co, co.modifier, func() {
		began := time.Now()
		go co.instrumentation.AModifyWithOperationsCall()
		defer func() { go co.instrumentation.AModifyWithOperationsDuration(time.Since(began)) }()

		res, err = co.modifier.ModifyWithOperations(key, id, ops, score, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Delete represents a way to delete various values into the store.
func (co *Coordinator) Delete(values []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (res int, err error) {
	if e := handle(co, co.deleter, func() {
		began := time.Now()
		go co.instrumentation.ADeleteCall()
		defer func() { go co.instrumentation.ADeleteDuration(time.Since(began)) }()

		res, err = co.deleter.Delete(values, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Rollback represents a way to rollback various values into the store.
func (co *Coordinator) Rollback(values []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (err error) {
	if e := handle(co, co.deleter, func() {
		began := time.Now()
		go co.instrumentation.ARollbackCall()
		defer func() { go co.instrumentation.ARollbackDuration(time.Since(began)) }()

		err = co.deleter.Rollback(values, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Select represents a way to request and select a member from the store.
func (co *Coordinator) Select(key, field bs.Key) (res s.KeyFieldScoreTxnValue, err error) {
	if e := handle(co, co.selector, func() {
		began := time.Now()
		go co.instrumentation.ASelectCall()
		defer func() { go co.instrumentation.ASelectDuration(time.Since(began)) }()

		res, err = co.selector.Select(key, field)
	}); e != nil {
		err = e
	}
	return
}

// SelectRange represents a way to request and select a range of members from
// the store that are under a certain limit.
func (co *Coordinator) SelectRange(key bs.Key, limit int, maxSize s.KeySizeExpiry) (res []s.KeyFieldScoreTxnValue, err error) {
	if e := handle(co, co.selector, func() {
		began := time.Now()
		go co.instrumentation.ASelectRangeCall()
		defer func() { go co.instrumentation.ASelectRangeDuration(time.Since(began)) }()

		res, err = co.selector.SelectRange(key, limit, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Keys defines a way to query the store for all the keys with in it.
func (co *Coordinator) Keys() (res []bs.Key, err error) {
	if e := handle(co, co.scanner, func() {
		began := time.Now()
		go co.instrumentation.AKeysCall()
		defer func() { go co.instrumentation.AKeysDuration(time.Since(began)) }()

		res, err = co.scanner.Keys()
	}); e != nil {
		err = e
	}
	return
}

// Size returns the size of the collection with in the store.
func (co *Coordinator) Size(key bs.Key) (res int, err error) {
	if e := handle(co, co.scanner, func() {
		began := time.Now()
		go co.instrumentation.ASizeCall()
		defer func() { go co.instrumentation.ASizeDuration(time.Since(began)) }()

		res, err = co.scanner.Size(key)
	}); e != nil {
		err = e
	}
	return
}

// Members represents all the items with in the store for a particular key.
func (co *Coordinator) Members(key bs.Key) (res []bs.Key, err error) {
	if e := handle(co, co.scanner, func() {
		began := time.Now()
		go co.instrumentation.AMembersCall()
		defer func() { go co.instrumentation.AMembersDuration(time.Since(began)) }()

		res, err = co.scanner.Members(key)
	}); e != nil {
		err = e
	}
	return
}

// Repair defines a way to request a possible repair of the store of a
// particular key.
func (co *Coordinator) Repair(elements []s.KeyFieldTxnValue, maxSize s.KeySizeExpiry) (err error) {
	if e := handle(co, co.repairer, func() {
		began := time.Now()
		go co.instrumentation.ARepairCall()
		defer func() { go co.instrumentation.ARepairDuration(time.Since(began)) }()

		err = co.repairer.Repair(elements, maxSize)
	}); e != nil {
		err = e
	}
	return
}

// Query defines a way to request a possible query of the store of a
// particular key.
func (co *Coordinator) Query(key bs.Key,
	options s.QueryOptions,
	maxSize s.SizeExpiry,
) (res []s.QueryRecord, err error) {
	if e := handle(co, co.inspector, func() {
		began := time.Now()
		go co.instrumentation.AQueryCall()
		defer func() { go co.instrumentation.AQueryDuration(time.Since(began)) }()

		res, err = co.inspector.Query(key, options, maxSize)
	}); e != nil {
		err = e
	}
	return
}

func (co *Coordinator) Lock(ns b.Namespace) (b.SemaphoreUnlock, error) {
	return co.consul.Lock(ns)
}

// Pause the coordinator
func (co *Coordinator) Pause() {
	go co.instrumentation.APauseCall()

	co.mutex.Lock()
	defer co.mutex.Unlock()

	if !co.running {
		return
	}

	if !co.paused {
		co.paused = true
		co.manager.Stop()
	}
}

// Resume the coordinator
func (co *Coordinator) Resume() {
	go co.instrumentation.AResumeCall()

	co.mutex.Lock()
	defer co.mutex.Unlock()

	if !co.running {
		return
	}

	if co.paused {
		co.paused = false
		co.cond.Signal()
	}
}

// Topology reloads the coordinator with all the various new strategies.
func (co *Coordinator) Topology(e *env.Env) error {
	began := time.Now()
	go co.instrumentation.ATopologyCall()
	defer func() { go co.instrumentation.ATopologyDuration(time.Since(began)) }()

	if clusters, err := newStoreClusters(e, co.storeOpts); err == nil {
		if err := co.store.Topology(clusters); err != nil {
			return err
		}
	} else {
		return err
	}

	if clusters, err := newCounterClusters(e); err == nil {
		if err := co.counter.Topology(clusters); err != nil {
			return err
		}
	} else {
		return err
	}

	if clusters, err := newNotifierClusters(e); err == nil {
		if err := co.notifier.Topology(clusters); err != nil {
			return err
		}
	} else {
		return err
	}

	if clusters, err := newPersistenceClusters(e, co.transformer); err == nil {
		if err := co.persistence.Topology(clusters); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (co *Coordinator) Quit() {
	co.mutex.Lock()
	defer co.mutex.Unlock()

	// We don't need to quit, as we're not running!
	if !co.running {
		return
	}

	// Firstly make sure we set the coordinator as stopped.
	co.running = false
}

type coordinatorAccessor struct {
	co *Coordinator
}

func NewCoordinatorAccessor(co *Coordinator) *coordinatorAccessor {
	return &coordinatorAccessor{co}
}

func (c *coordinatorAccessor) AlertManager() alertmanager.AlertManager {
	return c.co.alertmanager
}

func (c *coordinatorAccessor) StoreOpts() *r.Options {
	return c.co.storeOpts
}

// LifeCycleService is aims to manage the basic shutdown of a service. With the
// idea that it knows about how many users are currently active with in it.
type LifeCycleService struct {
	mutex   *sync.Mutex
	counter int
}

func newLifeCycleService() *LifeCycleService {
	return &LifeCycleService{mutex: &sync.Mutex{}, counter: 0}
}

func (c *LifeCycleService) In() {
	c.mutex.Lock()
	c.counter++
	c.mutex.Unlock()
}

func (c *LifeCycleService) Out() {
	c.mutex.Lock()
	c.counter--
	c.mutex.Unlock()
}

func (c *LifeCycleService) Empty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.counter < 1
}

func (c *LifeCycleService) Quit() (err error) {
	var (
		ticker = time.NewTicker(defaultQuitTicker).C
		timer  = time.NewTimer(defaultQuitTimeout).C
	)

loop:
	for {
		select {
		case <-ticker:
			if c.Empty() {
				break loop
			}
		case <-timer:
			err = typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Shutdown failed due to quit timeout.")
			break loop
		}
	}
	return
}
