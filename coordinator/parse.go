package coordinator

import (
	"io"

	a "github.com/SimonRichardson/colossus/alertmanager"
	ap "github.com/SimonRichardson/colossus/alertmanager/parse"
	blist "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/consul"
	"github.com/SimonRichardson/colossus/cluster/counter"
	"github.com/SimonRichardson/colossus/cluster/notifier"
	"github.com/SimonRichardson/colossus/cluster/persistence"
	"github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/env"
	c "github.com/SimonRichardson/colossus/farm/counter"
	n "github.com/SimonRichardson/colossus/farm/notifier"
	p "github.com/SimonRichardson/colossus/farm/persistence"
	s "github.com/SimonRichardson/colossus/farm/store"
	r "github.com/SimonRichardson/colossus/fusion/redis"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
	i "github.com/SimonRichardson/colossus/instrumentation"
	ip "github.com/SimonRichardson/colossus/instrumentation/parse"
	"github.com/SimonRichardson/colossus/selectors"
)

func newInstrumentation(e *env.Env, writer io.Writer) (i.Instrumentation, error) {
	return ip.ParseString(e.Instrumentation,
		ip.InstrumentationOptions{
			e.StatsdAddress,
			e.StatsdSampleRate,
			writer,
			e.LogsInstance,
			e.LogsBufferDuration,
			e.LogsTimeout,
		},
	)
}

func newAlertManager(e *env.Env) (a.AlertManager, error) {
	return ap.ParseString(e.AlertManager,
		ap.AlertManagerOptions{e.StatsdAddress, e.StatsdSampleRate},
	)
}

func newCounterClusters(e *env.Env) ([]counter.Cluster, error) {
	clusters, err := c.ParseString(
		e.CounterInstances,
		e.CounterConnectTimeout, e.CounterReadTimeout, e.CounterWriteTimeout,
		e.CounterPoolRoutingStrategy,
		e.CounterMaxSize,
		e.RedisCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newCounterFarm(e *env.Env,
	instr i.Instrumentation,
) (*c.Farm, error) {
	var (
		err         error
		clusters    []counter.Cluster
		insStrategy c.InsertCreator
		delStrategy c.DeleteCreator
		repStrategy c.RepairCreator
		scaStrategy c.ScanCreator
	)

	if clusters, err = newCounterClusters(e); err != nil {
		return nil, err
	}

	if insStrategy, err = c.ParseInsertStrategy(e.GetInsertOptions(env.Counter)); err != nil {
		return nil, err
	}

	if delStrategy, err = c.ParseDeleteStrategy(e.GetDeleteOptions(env.Counter)); err != nil {
		return nil, err
	}

	if repStrategy, err = c.ParseRepairStrategy(e.GetRepairOptions(env.Counter)); err != nil {
		return nil, err
	}

	if scaStrategy, err = c.ParseScanStrategy(e.GetScanOptions(env.Counter)); err != nil {
		return nil, err
	}

	return c.New(clusters,
		insStrategy,
		delStrategy,
		scaStrategy,
		repStrategy,
		instr,
	), nil
}

func newStoreOptions(e *env.Env, kvs blist.KeyStore) *s.Options {
	switch common.Normalise(e.StorePoolRoutingStrategy) {
	case "keystore":
		return &s.Options{
			KeyStorePrefix: fs.Prefix(e.StoreKeyStorePrefix),
			KeyStoreTicker: common.Tick(e.StoreKeyStoreDelay),
			KeyStore:       kvs,
		}
	}
	return nil
}

func newRedisOptions(opts *s.Options) *r.RedisOptions {
	if opts != nil {
		return &r.RedisOptions{
			KeyStorePrefix: opts.KeyStorePrefix,
			KeyStoreTicker: opts.KeyStoreTicker,
			KeyStore:       opts.KeyStore,
		}
	}
	return nil
}

func newStoreClusters(e *env.Env, opts *s.Options) ([]store.Cluster, error) {
	clusters, err := s.ParseString(
		e.StoreInstances,
		e.StoreConnectTimeout, e.StoreReadTimeout, e.StoreWriteTimeout,
		e.StorePoolRoutingStrategy,
		e.StoreMaxSize,
		newRedisOptions(opts),
		e.RedisCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newStoreFarm(e *env.Env,
	instr i.Instrumentation,
	opts *s.Options,
) (*s.Farm, error) {
	var (
		err         error
		clusters    []store.Cluster
		selStrategy s.SelectCreator
		insStrategy s.InsertCreator
		delStrategy s.DeleteCreator
		repStrategy s.RepairCreator
		scaStrategy s.ScanCreator
	)

	if clusters, err = newStoreClusters(e, opts); err != nil {
		return nil, err
	}

	if selStrategy, err = s.ParseSelectStrategy(e.GetSelectOptions(env.Store)); err != nil {
		return nil, err
	}

	if insStrategy, err = s.ParseInsertStrategy(e.GetInsertOptions(env.Store)); err != nil {
		return nil, err
	}

	if delStrategy, err = s.ParseDeleteStrategy(e.GetDeleteOptions(env.Store)); err != nil {
		return nil, err
	}

	if repStrategy, err = s.ParseRepairStrategy(e.GetRepairOptions(env.Store)); err != nil {
		return nil, err
	}

	if scaStrategy, err = s.ParseScanStrategy(e.GetScanOptions(env.Store)); err != nil {
		return nil, err
	}

	return s.New(clusters,
		selStrategy,
		insStrategy,
		delStrategy,
		scaStrategy,
		repStrategy,
		instr,
	), nil
}

func newPersistenceClusters(e *env.Env, transformer selectors.Transformer) ([]persistence.Cluster, error) {
	clusters, err := p.ParseString(
		e.MongoInstances,
		e.MongoConnectTimeout,
		e.PersistencePoolRoutingStrategy,
		e.PersistenceDbName, e.PersistenceKeyPrefix,
		transformer,
		e.PersistenceMaxSize,
		e.MongoCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newPersistenceFarm(e *env.Env,
	instr i.Instrumentation,
	transformer selectors.Transformer,
) (*p.Farm, error) {
	var (
		err         error
		clusters    []persistence.Cluster
		insStrategy p.InsertCreator
		delStrategy p.DeleteCreator
		repStrategy p.RepairCreator
	)

	if clusters, err = newPersistenceClusters(e, transformer); err != nil {
		return nil, err
	}

	if insStrategy, err = p.ParseInsertStrategy(e.GetInsertOptions(env.Persistence)); err != nil {
		return nil, err
	}

	if delStrategy, err = p.ParseDeleteStrategy(e.GetDeleteOptions(env.Persistence)); err != nil {
		return nil, err
	}

	if repStrategy, err = p.ParseRepairStrategy(e.GetRepairOptions(env.Persistence)); err != nil {
		return nil, err
	}

	return p.New(clusters,
		insStrategy,
		delStrategy,
		repStrategy,
		instr,
	), nil
}

func newNotifierClusters(e *env.Env) ([]notifier.Cluster, error) {
	clusters, err := n.ParseString(
		e.NotifierInstances,
		e.NotifierConnectTimeout, e.NotifierReadTimeout, e.NotifierWriteTimeout,
		e.NotifierPoolRoutingStrategy,
		e.NotifierMaxSize,
		e.RedisCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newNotifierFarm(e *env.Env,
	instr i.Instrumentation,
) (*n.Farm, error) {
	var (
		err         error
		clusters    []notifier.Cluster
		notStrategy n.NotifyCreator
	)

	if clusters, err = newNotifierClusters(e); err != nil {
		return nil, err
	}

	if notStrategy, err = n.ParseNotifyStrategy(e.GetNotifyOptions(env.Notifier)); err != nil {
		return nil, err
	}

	return n.New(clusters,
		notStrategy,
		instr,
	), nil
}

func newConsulClusters(e *env.Env) ([]consul.Cluster, error) {
	clusters, err := consul.ParseString(
		e.ConsulInstances,
		e.ConsulCheckId,
		e.ConsulOutput,
		e.ConsulMaxSize,
		e.ConsulClientCreator,
	)

	if err != nil {
		return nil, err
	}

	return clusters, err
}

func newConsulService(e *env.Env, instr i.Instrumentation) (*consul.Service, error) {
	var (
		clusters    []consul.Cluster
		semStrategy consul.SemaphoreCreator
		hrtStrategy consul.HeartbeatCreator
		kvsStrategy consul.KeyStoreCreator

		opts consul.StrategyOptions

		err error
	)

	if clusters, err = newConsulClusters(e); err != nil {
		return nil, err
	}

	opts = consulStrategyOptions(e.GetSemaphoreOptions(env.Consul))
	if semStrategy, err = consul.ParseSemaphoreStrategy(opts); err != nil {
		return nil, err
	}

	opts = consulStrategyOptions(e.GetHeartbeatOptions(env.Consul))
	if hrtStrategy, err = consul.ParseHeartbeatStrategy(opts); err != nil {
		return nil, err
	}

	opts = consulStrategyOptions(e.GetKeyStoreOptions(env.Consul))
	if kvsStrategy, err = consul.ParseKeyStoreStrategy(opts); err != nil {
		return nil, err
	}

	return consul.New(clusters,
		semStrategy,
		hrtStrategy,
		kvsStrategy,
		instr,
	), nil
}

func consulStrategyOptions(opts env.StrategyOptions) consul.StrategyOptions {
	return consul.StrategyOptions{
		Strategy:            opts.Strategy,
		Tactic:              opts.Tactic,
		RequestsDuration:    opts.RequestsDuration,
		RequestsPerDuration: opts.RequestsPerDuration,
		Quorum:              opts.Quorum,
	}
}
