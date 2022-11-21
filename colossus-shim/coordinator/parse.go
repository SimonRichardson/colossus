package coordinator

import (
	"io"

	a "github.com/SimonRichardson/colossus/alertmanager"
	ap "github.com/SimonRichardson/colossus/alertmanager/parse"
	"github.com/SimonRichardson/colossus/colossus-shim/cluster/score"
	"github.com/SimonRichardson/colossus/colossus-shim/env"
	t "github.com/SimonRichardson/colossus/colossus-shim/farm/score"
	i "github.com/SimonRichardson/colossus/instrumentation"
	ip "github.com/SimonRichardson/colossus/instrumentation/parse"
)

func newInstrumentation(e *env.Env, writer io.Writer) (i.Instrumentation, error) {
	return ip.ParseString(e.C.Instrumentation,
		ip.InstrumentationOptions{
			e.C.StatsdAddress,
			e.C.StatsdSampleRate,
			writer,
			e.C.LogsInstance,
			e.C.LogsBufferDuration,
			e.C.LogsTimeout,
		},
	)
}

func newAlertManager(e *env.Env) (a.AlertManager, error) {
	return ap.ParseString(e.C.AlertManager,
		ap.AlertManagerOptions{e.C.StatsdAddress, e.C.StatsdSampleRate},
	)
}

func newScoreClusters(e *env.Env) ([]score.Cluster, error) {
	clusters, err := t.ParseString(
		e.ShimRedisInstances,
		e.ShimRedisConnectTimeout, e.ShimRedisReadTimeout, e.ShimRedisWriteTimeout,
		e.ScorePoolRoutingStrategy,
		e.ScoreMaxSize,
		e.C.RedisCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newScoreFarm(e *env.Env, instr i.Instrumentation) (*t.Farm, error) {
	var (
		err         error
		clusters    []score.Cluster
		incStrategy t.IncrementCreator
	)

	if clusters, err = newScoreClusters(e); err != nil {
		return nil, err
	}

	if incStrategy, err = t.ParseIncrementStrategy(e.GetIncrementOptions(env.Score)); err != nil {
		return nil, err
	}

	return t.New(clusters,
		incStrategy,
		instr,
	), nil
}
