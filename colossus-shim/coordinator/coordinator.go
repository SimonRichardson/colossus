package coordinator

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-shim/env"
	"github.com/SimonRichardson/colossus/colossus-shim/farm/score"
	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	defaultDebugExeceptions = true
)

// Coordinator defines a single point for accessing the various services.
type Coordinator struct {
	mutex *sync.Mutex
	cond  *sync.Cond

	paused bool

	score *score.Farm

	instrumentation instrumentation.Instrumentation
	alertmanager    alertmanager.AlertManager
}

func New(e *env.Env) *Coordinator {
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

		paused: false,

		instrumentation: instr,
		alertmanager:    alert,
	}

	if err := co.init(e); err != nil {
		typex.Fatal(err)
	}

	return co
}

func (co *Coordinator) init(e *env.Env) error {
	var (
		score *score.Farm

		err error
	)

	if score, err = newScoreFarm(e, co.instrumentation); err != nil {
		return err
	}

	co.score = score

	return nil
}

func handle(co *Coordinator, f func()) (err error) {
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

func (co *Coordinator) Increment(key bs.Key, t time.Time) (res int, err error) {
	if e := handle(co, func() {
		res, err = co.score.Increment(key, t)
	}); e != nil {
		err = e
	}
	return
}

func (co *Coordinator) Quit() {
}
