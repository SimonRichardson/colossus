package env

import (
	"time"

	consul "github.com/SimonRichardson/colossus/blist/services/consul/client"
	"github.com/SimonRichardson/colossus/fusion/mongo"
	"github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// TODO :This is a bit long winded! Need to move to a structured file setup.

// Env is a struct that contains all the environmental settings.
type Env struct {
	source *viper.Viper

	HttpAddress      string
	HttpReadTimeout  time.Duration
	HttpWriteTimeout time.Duration

	Version string

	Instrumentation string
	AlertManager    string

	// Logs
	Logs               string
	LogsInstance       string
	LogsTimeout        string
	LogsBufferDuration time.Duration

	StatsdAddress     string
	StatsdSampleRate  float32
	PrometheusMetrics bool

	// General

	InsertStrategyTactic string

	RepairStrategyTactic      string
	RepairStrategyPerDuration int
	RepairStrategyDuration    string

	// Redis

	RedisCreator redis.RedisCreator

	// Mongo

	MongoInstances      string
	MongoConnectTimeout string

	MongoCreator mongo.SessionCreator

	// Counter

	CounterInstances      string
	CounterConnectTimeout string
	CounterReadTimeout    string
	CounterWriteTimeout   string

	CounterMaxSize             int
	CounterPoolRoutingStrategy string

	CounterInsertStrategy    string
	CounterInsertTactic      string
	CounterInsertPerDuration int
	CounterInsertDuration    string
	CounterInsertQuorum      float64

	CounterDeleteStrategy    string
	CounterDeleteTactic      string
	CounterDeletePerDuration int
	CounterDeleteDuration    string
	CounterDeleteQuorum      float64

	CounterScanStrategy    string
	CounterScanTactic      string
	CounterScanPerDuration int
	CounterScanDuration    string
	CounterScanQuorum      float64

	CounterRepairStrategy    string
	CounterRepairTactic      string
	CounterRepairPerDuration int
	CounterRepairDuration    string
	CounterRepairQuorum      float64

	// Store

	StoreInstances      string
	StoreConnectTimeout string
	StoreReadTimeout    string
	StoreWriteTimeout   string

	StoreMaxSize             int
	StorePoolRoutingStrategy string
	StoreKeyStorePrefix      string
	StoreKeyStoreDelay       time.Duration

	StoreSelectStrategy    string
	StoreSelectTactic      string
	StoreSelectPerDuration int
	StoreSelectDuration    string
	StoreSelectQuorum      float64

	StoreInsertStrategy    string
	StoreInsertTactic      string
	StoreInsertPerDuration int
	StoreInsertDuration    string
	StoreInsertQuorum      float64

	StoreDeleteStrategy    string
	StoreDeleteTactic      string
	StoreDeletePerDuration int
	StoreDeleteDuration    string
	StoreDeleteQuorum      float64

	StoreScanStrategy    string
	StoreScanTactic      string
	StoreScanPerDuration int
	StoreScanDuration    string
	StoreScanQuorum      float64

	StoreRepairStrategy    string
	StoreRepairTactic      string
	StoreRepairPerDuration int
	StoreRepairDuration    string
	StoreRepairQuorum      float64

	// Notifier

	NotifierInstances      string
	NotifierConnectTimeout string
	NotifierReadTimeout    string
	NotifierWriteTimeout   string

	NotifierMaxSize             int
	NotifierPoolRoutingStrategy string

	NotifierNotifyStrategy    string
	NotifierNotifyTactic      string
	NotifierNotifyPerDuration int
	NotifierNotifyDuration    string
	NotifierNotifyQuorum      float64

	// Persistence

	PersistenceDbName              string
	PersistenceKeyPrefix           string
	PersistenceMaxSize             int
	PersistencePoolRoutingStrategy string

	PersistenceInsertStrategy    string
	PersistenceInsertTactic      string
	PersistenceInsertPerDuration int
	PersistenceInsertDuration    string
	PersistenceInsertQuorum      float64

	PersistenceDeleteStrategy    string
	PersistenceDeleteTactic      string
	PersistenceDeletePerDuration int
	PersistenceDeleteDuration    string
	PersistenceDeleteQuorum      float64

	PersistenceRepairStrategy    string
	PersistenceRepairTactic      string
	PersistenceRepairPerDuration int
	PersistenceRepairDuration    string
	PersistenceRepairQuorum      float64

	// Manager

	ManagerRepairStrategy    string
	ManagerRepairTactic      string
	ManagerRepairPerDuration int
	ManagerRepairDuration    string

	// Consul

	ConsulInstances string

	ConsulMaxSize int
	ConsulCheckId string
	ConsulOutput  string

	ConsulClientCreator consul.ClientCreator

	ConsulSemaphoreStrategy    string
	ConsulSemaphoreTactic      string
	ConsulSemaphorePerDuration int
	ConsulSemaphoreDuration    string

	ConsulHeartbeatStrategy    string
	ConsulHeartbeatTactic      string
	ConsulHeartbeatPerDuration int
	ConsulHeartbeatDuration    string
	ConsulHeartbeatFrequency   time.Duration

	ConsulKeyStoreStrategy    string
	ConsulKeyStoreTactic      string
	ConsulKeyStorePerDuration int
	ConsulKeyStoreDuration    string
}

// Type describes what stragegy options are available
type Type int

const (
	Counter = iota
	Store
	Persistence
	Notifier
	Coordinator
	Manager
	Consul
)

// StrategyOptions defines what options are available when creating a and using
// a stragegy.
type StrategyOptions struct {
	Strategy            string
	Tactic              string
	RequestsPerDuration int
	RequestsDuration    string
	Quorum              float64
}

// New returns a new Env object which contains all the environmental variables
// in the object.
func New(paths []string) *Env {
	v := viper.New()
	v.SetConfigName("config")
	v.AutomaticEnv()

	if len(paths) != 0 {
		for _, path := range paths {
			v.AddConfigPath(path)
		}

		if err := v.ReadInConfig(); err != nil {
			typex.Fatal(err)
		}
	}

	v.SetDefault("http_address", ":9002")
	v.SetDefault("http_read_timeout", "10s")
	v.SetDefault("http_write_timeout", "30s")

	v.SetDefault("version", "0.0.1")

	v.SetDefault("instrumentation", "PlainText")
	v.SetDefault("alert_manager", "PlainText")

	v.SetDefault("logs", "PlainText-Buffered")
	v.SetDefault("logs_instance", "tcp://logs:6379")
	v.SetDefault("logs_timeout", "1m")
	v.SetDefault("logs_buffer_duration", "100ms")

	v.SetDefault("statsd_address", "")
	v.SetDefault("statsd_sample_rate", 0.1)
	v.SetDefault("prometheus_metrics", false)

	v.SetDefault("insert_strategy", "Counter")

	v.SetDefault("repair_strategy", "NonBlocking")
	v.SetDefault("repair_strategy_per_duration", 100)
	v.SetDefault("repair_strategy_duration", "1m")

	v.SetDefault("mongo_instances", "mongo:27017")
	v.SetDefault("mongo_connect_timeout", "1m")

	v.SetDefault("counter_instances", "tcp://counter1:6379;tcp://counter2:6379;tcp://counter3:6379")
	v.SetDefault("counter_connect_timeout", "1m")
	v.SetDefault("counter_read_timeout", "30s")
	v.SetDefault("counter_write_timeout", "30s")

	v.SetDefault("counter_max_size", 1000)
	v.SetDefault("counter_pool_routing_strategy", "Hash")

	v.SetDefault("counter_insert_strategy", "InsertAllReadAll")
	v.SetDefault("counter_insert_tactic", "NonBlocking")
	v.SetDefault("counter_insert_per_duration", 0)
	v.SetDefault("counter_insert_duration", 0)
	v.SetDefault("counter_insert_quorum", 0.51)

	v.SetDefault("counter_delete_strategy", "DeleteAllReadAll")
	v.SetDefault("counter_delete_tactic", "NonBlocking")
	v.SetDefault("counter_delete_per_duration", 0)
	v.SetDefault("counter_delete_duration", 0)
	v.SetDefault("counter_delete_quorum", 0.51)

	v.SetDefault("counter_scan_strategy", "ScanAllReadAll")
	v.SetDefault("counter_scan_tactic", "NonBlocking")
	v.SetDefault("counter_scan_per_duration", 1)
	v.SetDefault("counter_scan_duration", "1m")
	v.SetDefault("counter_scan_quorum", 0.51)

	v.SetDefault("counter_repair_strategy", "RepairAll")
	v.SetDefault("counter_repair_tactic", "NonBlocking")
	v.SetDefault("counter_repair_per_duration", 1)
	v.SetDefault("counter_repair_duration", "1m")
	v.SetDefault("counter_repair_quorum", 0.51)

	v.SetDefault("store_instances", "tcp://store1_a:6379,tcp://store1_b:6379;tcp://store2_a:6379,tcp://store2_b:6379;tcp://store3_a:6379,tcp://store3_b:6379")
	v.SetDefault("store_connect_timeout", "1m")
	v.SetDefault("store_read_timeout", "30s")
	v.SetDefault("store_write_timeout", "30s")

	v.SetDefault("store_max_size", 1000)
	v.SetDefault("store_pool_routing_strategy", "Hash")
	v.SetDefault("store_key_store_prefix", "colossus_store_kv")
	v.SetDefault("store_key_store_delay", time.Minute)

	v.SetDefault("store_select_strategy", "SelectQuorumReadAll")
	v.SetDefault("store_select_tactic", "NonBlocking")
	v.SetDefault("store_select_per_duration", 0)
	v.SetDefault("store_select_duration", 0)
	v.SetDefault("store_select_quorum", 0.51)

	v.SetDefault("store_insert_strategy", "InsertAllReadAll")
	v.SetDefault("store_insert_tactic", "NonBlocking")
	v.SetDefault("store_insert_per_duration", 0)
	v.SetDefault("store_insert_duration", 0)
	v.SetDefault("store_insert_quorum", 0.51)

	v.SetDefault("store_delete_strategy", "DeleteAllReadAll")
	v.SetDefault("store_delete_tactic", "NonBlocking")
	v.SetDefault("store_delete_per_duration", 0)
	v.SetDefault("store_delete_duration", 0)
	v.SetDefault("store_delete_quorum", 0.51)

	v.SetDefault("store_scan_strategy", "ScanAllReadAll")
	v.SetDefault("store_scan_tactic", "NonBlocking")
	v.SetDefault("store_scan_per_duration", 1)
	v.SetDefault("store_scan_duration", "1m")
	v.SetDefault("store_scan_quorum", 0.51)

	v.SetDefault("store_repair_strategy", "RepairAll")
	v.SetDefault("store_repair_tactic", "NonBlocking")
	v.SetDefault("store_repair_per_duration", 1)
	v.SetDefault("store_repair_duration", "1m")
	v.SetDefault("store_repair_quorum", 0.51)

	v.SetDefault("notifier_instances", "tcp://notifier:6379")
	v.SetDefault("notifier_connect_timeout", "1m")
	v.SetDefault("notifier_read_timeout", "30s")
	v.SetDefault("notifier_write_timeout", "30s")

	v.SetDefault("notifier_max_size", 100)
	v.SetDefault("notifier_pool_routing_strategy", "Hash")

	v.SetDefault("notifier_notify_strategy", "Bulk")
	v.SetDefault("notifier_notify_tactic", "NonBlocking")
	v.SetDefault("notifier_notify_per_duration", 0)
	v.SetDefault("notifier_notify_duration", 0)
	v.SetDefault("notifier_notify_quorum", 0.51)

	v.SetDefault("persistence_db_name", "db")
	v.SetDefault("persistence_key_prefix", "tickets_")
	v.SetDefault("persistence_max_size", 100)
	v.SetDefault("persistence_pool_routing_strategy", "Hash")

	v.SetDefault("persistence_insert_strategy", "InsertAllReadAll")
	v.SetDefault("persistence_insert_tactic", "NonBlocking")
	v.SetDefault("persistence_insert_per_duration", 0)
	v.SetDefault("persistence_insert_duration", 0)
	v.SetDefault("persistence_insert_quorum", 0.51)

	v.SetDefault("persistence_delete_strategy", "DeleteAllReadAll")
	v.SetDefault("persistence_delete_tactic", "NonBlocking")
	v.SetDefault("persistence_delete_per_duration", 0)
	v.SetDefault("persistence_delete_duration", 0)
	v.SetDefault("persistence_delete_quorum", 0.51)

	v.SetDefault("persistence_repair_strategy", "RepairAll")
	v.SetDefault("persistence_repair_tactic", "NonBlocking")
	v.SetDefault("persistence_repair_per_duration", 1)
	v.SetDefault("persistence_repair_duration", "1m")
	v.SetDefault("persistence_repair_quorum", 0.51)

	v.SetDefault("manager_repair_strategy", "collect")
	v.SetDefault("manager_repair_tactic", "NonBlocking")
	v.SetDefault("manager_repair_per_duration", 1)
	v.SetDefault("manager_repair_duration", "1m")

	v.SetDefault("consul_instances", "consul:8500")

	v.SetDefault("consul_max_size", 1)
	v.SetDefault("consul_check_id", "service:colossus")
	v.SetDefault("consul_output", "colossus")

	v.SetDefault("consul_semaphore_strategy", "Semaphore")
	v.SetDefault("consul_semaphore_tactic", "NonBlocking")
	v.SetDefault("consul_semaphore_per_duration", 0)
	v.SetDefault("consul_semaphore_duration", 0)

	v.SetDefault("consul_heartbeat_strategy", "Heartbeat")
	v.SetDefault("consul_heartbeat_tactic", "RateLimited")
	v.SetDefault("consul_heartbeat_per_duration", 1)
	v.SetDefault("consul_heartbeat_duration", "10s")
	v.SetDefault("consul_heartbeat_frequency", "60s")

	v.SetDefault("consul_keystore_strategy", "KeyStore")
	v.SetDefault("consul_keystore_tactic", "NonBlocking")
	v.SetDefault("consul_keystore_per_duration", 1)
	v.SetDefault("consul_keystore_duration", "10s")

	e := &Env{
		source: v,
	}

	e.read()

	return e
}

func (e *Env) read() {
	e.HttpAddress = e.source.GetString("http_address")
	e.HttpReadTimeout = e.source.GetDuration("http_read_timeout")
	e.HttpWriteTimeout = e.source.GetDuration("http_write_timeout")

	e.Version = e.source.GetString("version")

	e.Instrumentation = e.source.GetString("instrumentation")
	e.AlertManager = e.source.GetString("alert_manager")

	e.Logs = e.source.GetString("logs")
	e.LogsInstance = e.source.GetString("logs_instance")
	e.LogsTimeout = e.source.GetString("logs_timeout")
	e.LogsBufferDuration = e.source.GetDuration("logs_buffer_duration")

	e.StatsdAddress = e.source.GetString("statsd_address")
	e.StatsdSampleRate = float32(e.source.GetFloat64("statsd_sample_rate"))
	e.PrometheusMetrics = e.source.GetBool("prometheus_metrics")

	e.InsertStrategyTactic = e.source.GetString("insert_strategy")

	e.RepairStrategyTactic = e.source.GetString("repair_strategy")
	e.RepairStrategyPerDuration = e.source.GetInt("repair_strategy_per_duration")
	e.RepairStrategyDuration = e.source.GetString("repair_strategy_duration")

	e.MongoInstances = e.source.GetString("mongo_instances")
	e.MongoConnectTimeout = e.source.GetString("mongo_connect_timeout")

	e.CounterInstances = e.source.GetString("counter_instances")
	e.CounterConnectTimeout = e.source.GetString("counter_connect_timeout")
	e.CounterReadTimeout = e.source.GetString("counter_read_timeout")
	e.CounterWriteTimeout = e.source.GetString("counter_write_timeout")

	e.CounterMaxSize = e.source.GetInt("counter_max_size")
	e.CounterPoolRoutingStrategy = e.source.GetString("counter_pool_routing_strategy")

	e.CounterInsertStrategy = e.source.GetString("counter_insert_strategy")
	e.CounterInsertTactic = e.source.GetString("counter_insert_tactic")
	e.CounterInsertPerDuration = e.source.GetInt("counter_insert_per_duration")
	e.CounterInsertDuration = e.source.GetString("counter_insert_duration")
	e.CounterInsertQuorum = e.source.GetFloat64("counter_insert_quorum")

	e.CounterDeleteStrategy = e.source.GetString("counter_delete_strategy")
	e.CounterDeleteTactic = e.source.GetString("counter_delete_tactic")
	e.CounterDeletePerDuration = e.source.GetInt("counter_delete_per_duration")
	e.CounterDeleteDuration = e.source.GetString("counter_delete_duration")
	e.CounterDeleteQuorum = e.source.GetFloat64("counter_delete_quorum")

	e.CounterScanStrategy = e.source.GetString("counter_scan_strategy")
	e.CounterScanTactic = e.source.GetString("counter_scan_tactic")
	e.CounterScanPerDuration = e.source.GetInt("counter_scan_per_duration")
	e.CounterScanDuration = e.source.GetString("counter_scan_duration")
	e.CounterScanQuorum = e.source.GetFloat64("counter_scan_quorum")

	e.CounterRepairStrategy = e.source.GetString("counter_repair_strategy")
	e.CounterRepairTactic = e.source.GetString("counter_repair_tactic")
	e.CounterRepairPerDuration = e.source.GetInt("counter_repair_per_duration")
	e.CounterRepairDuration = e.source.GetString("counter_repair_duration")
	e.CounterRepairQuorum = e.source.GetFloat64("counter_repair_quorum")

	e.StoreInstances = e.source.GetString("store_instances")
	e.StoreConnectTimeout = e.source.GetString("store_connect_timeout")
	e.StoreReadTimeout = e.source.GetString("store_read_timeout")
	e.StoreWriteTimeout = e.source.GetString("store_write_timeout")

	e.StoreMaxSize = e.source.GetInt("store_max_size")
	e.StorePoolRoutingStrategy = e.source.GetString("store_pool_routing_strategy")
	e.StoreKeyStorePrefix = e.source.GetString("store_key_store_prefix")
	e.StoreKeyStoreDelay = e.source.GetDuration("store_key_store_delay")

	e.StoreSelectStrategy = e.source.GetString("store_select_strategy")
	e.StoreSelectTactic = e.source.GetString("store_select_tactic")
	e.StoreSelectPerDuration = e.source.GetInt("store_select_per_duration")
	e.StoreSelectDuration = e.source.GetString("store_select_duration")
	e.StoreSelectQuorum = e.source.GetFloat64("store_select_quorum")

	e.StoreInsertStrategy = e.source.GetString("store_insert_strategy")
	e.StoreInsertTactic = e.source.GetString("store_insert_tactic")
	e.StoreInsertPerDuration = e.source.GetInt("store_insert_per_duration")
	e.StoreInsertDuration = e.source.GetString("store_insert_duration")
	e.StoreInsertQuorum = e.source.GetFloat64("store_insert_quorum")

	e.StoreDeleteStrategy = e.source.GetString("store_delete_strategy")
	e.StoreDeleteTactic = e.source.GetString("store_delete_tactic")
	e.StoreDeletePerDuration = e.source.GetInt("store_delete_per_duration")
	e.StoreDeleteDuration = e.source.GetString("store_delete_duration")
	e.StoreDeleteQuorum = e.source.GetFloat64("store_delete_quorum")

	e.StoreScanStrategy = e.source.GetString("store_scan_strategy")
	e.StoreScanTactic = e.source.GetString("store_scan_tactic")
	e.StoreScanPerDuration = e.source.GetInt("store_scan_per_duration")
	e.StoreScanDuration = e.source.GetString("store_scan_duration")
	e.StoreScanQuorum = e.source.GetFloat64("store_scan_quorum")

	e.StoreRepairStrategy = e.source.GetString("store_repair_strategy")
	e.StoreRepairTactic = e.source.GetString("store_repair_tactic")
	e.StoreRepairPerDuration = e.source.GetInt("store_repair_per_duration")
	e.StoreRepairDuration = e.source.GetString("store_repair_duration")
	e.StoreRepairQuorum = e.source.GetFloat64("store_repair_quorum")

	e.NotifierInstances = e.source.GetString("notifier_instances")
	e.NotifierConnectTimeout = e.source.GetString("notifier_connect_timeout")
	e.NotifierReadTimeout = e.source.GetString("notifier_read_timeout")
	e.NotifierWriteTimeout = e.source.GetString("notifier_write_timeout")

	e.NotifierMaxSize = e.source.GetInt("notifier_max_size")
	e.NotifierPoolRoutingStrategy = e.source.GetString("notifier_pool_routing_strategy")

	e.NotifierNotifyStrategy = e.source.GetString("notifier_notify_strategy")
	e.NotifierNotifyTactic = e.source.GetString("notifier_notify_tactic")
	e.NotifierNotifyPerDuration = e.source.GetInt("notifier_notify_per_duration")
	e.NotifierNotifyDuration = e.source.GetString("notifier_notify_duration")
	e.NotifierNotifyQuorum = e.source.GetFloat64("notifier_notify_quorum")

	e.PersistenceDbName = e.source.GetString("persistence_db_name")
	e.PersistenceKeyPrefix = e.source.GetString("persistence_key_prefix")
	e.PersistenceMaxSize = e.source.GetInt("persistence_max_size")
	e.PersistencePoolRoutingStrategy = e.source.GetString("persistence_pool_routing_strategy")

	e.PersistenceInsertStrategy = e.source.GetString("persistence_insert_strategy")
	e.PersistenceInsertTactic = e.source.GetString("persistence_insert_tactic")
	e.PersistenceInsertPerDuration = e.source.GetInt("persistence_insert_per_duration")
	e.PersistenceInsertDuration = e.source.GetString("persistence_insert_duration")
	e.PersistenceInsertQuorum = e.source.GetFloat64("persistence_insert_quorum")

	e.PersistenceDeleteStrategy = e.source.GetString("persistence_delete_strategy")
	e.PersistenceDeleteTactic = e.source.GetString("persistence_delete_tactic")
	e.PersistenceDeletePerDuration = e.source.GetInt("persistence_delete_per_duration")
	e.PersistenceDeleteDuration = e.source.GetString("persistence_delete_duration")
	e.PersistenceDeleteQuorum = e.source.GetFloat64("persistence_delete_quorum")

	e.PersistenceRepairStrategy = e.source.GetString("persistence_repair_strategy")
	e.PersistenceRepairTactic = e.source.GetString("persistence_repair_tactic")
	e.PersistenceRepairPerDuration = e.source.GetInt("persistence_repair_per_duration")
	e.PersistenceRepairDuration = e.source.GetString("persistence_repair_duration")
	e.PersistenceRepairQuorum = e.source.GetFloat64("persistence_repair_quorum")

	e.ManagerRepairStrategy = e.source.GetString("manager_repair_strategy")
	e.ManagerRepairTactic = e.source.GetString("manager_repair_tactic")
	e.ManagerRepairPerDuration = e.source.GetInt("manager_repair_per_duration")
	e.ManagerRepairDuration = e.source.GetString("manager_repair_duration")

	e.ConsulInstances = e.source.GetString("consul_instances")

	e.ConsulMaxSize = e.source.GetInt("consul_max_size")
	e.ConsulCheckId = e.source.GetString("consul_check_id")
	e.ConsulOutput = e.source.GetString("consul_output")

	e.ConsulSemaphoreStrategy = e.source.GetString("consul_semaphore_strategy")
	e.ConsulSemaphoreTactic = e.source.GetString("consul_semaphore_tactic")
	e.ConsulSemaphorePerDuration = e.source.GetInt("consul_semaphore_per_duration")
	e.ConsulSemaphoreDuration = e.source.GetString("consul_semaphore_duration")

	e.ConsulHeartbeatStrategy = e.source.GetString("consul_heartbeat_strategy")
	e.ConsulHeartbeatTactic = e.source.GetString("consul_heartbeat_tactic")
	e.ConsulHeartbeatPerDuration = e.source.GetInt("consul_heartbeat_per_duration")
	e.ConsulHeartbeatDuration = e.source.GetString("consul_heartbeat_duration")
	e.ConsulHeartbeatFrequency = e.source.GetDuration("consul_heartbeat_frequency")

	e.ConsulKeyStoreStrategy = e.source.GetString("consul_keystore_strategy")
	e.ConsulKeyStoreTactic = e.source.GetString("consul_keystore_tactic")
	e.ConsulKeyStorePerDuration = e.source.GetInt("consul_keystore_per_duration")
	e.ConsulKeyStoreDuration = e.source.GetString("consul_keystore_duration")
}

// Watch allows the watching of the config file, which inturn allows to be
// notified about when the underlying configuration file changes.
func (e *Env) Watch() <-chan *Env {
	res := make(chan *Env)
	go func() {
		// Make sure we don't block
		e.source.WatchConfig()
		e.source.OnConfigChange(func(in fsnotify.Event) {
			e.read()
			res <- e
		})
	}()
	return res
}

// GetSelectOptions returns all the selection options required to run a
// selection in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetSelectOptions(t Type) StrategyOptions {
	switch t {
	case Store:
		return StrategyOptions{
			e.StoreSelectStrategy,
			e.StoreSelectTactic,
			e.StoreSelectPerDuration,
			e.StoreSelectDuration,
			e.StoreSelectQuorum,
		}
	}
	return StrategyOptions{}
}

// GetInsertOptions returns all the insertion options required to run a
// insertion in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetInsertOptions(t Type) StrategyOptions {
	switch t {
	case Counter:
		return StrategyOptions{
			e.CounterInsertStrategy,
			e.CounterInsertTactic,
			e.CounterInsertPerDuration,
			e.CounterInsertDuration,
			e.CounterInsertQuorum,
		}
	case Store:
		return StrategyOptions{
			e.StoreInsertStrategy,
			e.StoreInsertTactic,
			e.StoreInsertPerDuration,
			e.StoreInsertDuration,
			e.StoreInsertQuorum,
		}
	case Persistence:
		return StrategyOptions{
			e.PersistenceInsertStrategy,
			e.PersistenceInsertTactic,
			e.PersistenceInsertPerDuration,
			e.PersistenceInsertDuration,
			e.PersistenceInsertQuorum,
		}
	case Coordinator:
		return StrategyOptions{
			Tactic: e.InsertStrategyTactic,
		}
	}
	return StrategyOptions{}
}

// GetDeleteOptions returns all the deletion options required to run a
// deletion in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetDeleteOptions(t Type) StrategyOptions {
	switch t {
	case Counter:
		return StrategyOptions{
			e.CounterDeleteStrategy,
			e.CounterDeleteTactic,
			e.CounterDeletePerDuration,
			e.CounterDeleteDuration,
			e.CounterDeleteQuorum,
		}
	case Store:
		return StrategyOptions{
			e.StoreDeleteStrategy,
			e.StoreDeleteTactic,
			e.StoreDeletePerDuration,
			e.StoreDeleteDuration,
			e.StoreDeleteQuorum,
		}
	case Persistence:
		return StrategyOptions{
			e.PersistenceDeleteStrategy,
			e.PersistenceDeleteTactic,
			e.PersistenceDeletePerDuration,
			e.PersistenceDeleteDuration,
			e.PersistenceDeleteQuorum,
		}
	}
	return StrategyOptions{}
}

// GetScanOptions returns all the scanning options required to run a
// scanning in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetScanOptions(t Type) StrategyOptions {
	switch t {
	case Counter:
		return StrategyOptions{
			e.CounterScanStrategy,
			e.CounterScanTactic,
			e.CounterScanPerDuration,
			e.CounterScanDuration,
			e.CounterScanQuorum,
		}
	case Store:
		return StrategyOptions{
			e.StoreScanStrategy,
			e.StoreScanTactic,
			e.StoreScanPerDuration,
			e.StoreScanDuration,
			e.StoreScanQuorum,
		}
	}
	return StrategyOptions{}
}

// GetRepairOptions returns all the repairing options required to run a
// repairing in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetRepairOptions(t Type) StrategyOptions {
	switch t {
	case Counter:
		return StrategyOptions{
			e.CounterRepairStrategy,
			e.CounterRepairTactic,
			e.CounterRepairPerDuration,
			e.CounterRepairDuration,
			e.CounterRepairQuorum,
		}
	case Store:
		return StrategyOptions{
			e.StoreRepairStrategy,
			e.StoreRepairTactic,
			e.StoreRepairPerDuration,
			e.StoreRepairDuration,
			e.StoreRepairQuorum,
		}
	case Persistence:
		return StrategyOptions{
			e.PersistenceRepairStrategy,
			e.PersistenceRepairTactic,
			e.PersistenceRepairPerDuration,
			e.PersistenceRepairDuration,
			e.PersistenceRepairQuorum,
		}
	case Coordinator:
		return StrategyOptions{
			Tactic:              e.RepairStrategyTactic,
			RequestsPerDuration: e.RepairStrategyPerDuration,
			RequestsDuration:    e.RepairStrategyDuration,
		}
	case Manager:
		return StrategyOptions{
			Strategy:            e.ManagerRepairStrategy,
			Tactic:              e.ManagerRepairTactic,
			RequestsPerDuration: e.ManagerRepairPerDuration,
			RequestsDuration:    e.ManagerRepairDuration,
		}
	}
	return StrategyOptions{}
}

// GetNotifyOptions returns all the notifying options required to run a
// notification in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetNotifyOptions(t Type) StrategyOptions {
	switch t {
	case Notifier:
		return StrategyOptions{
			e.NotifierNotifyStrategy,
			e.NotifierNotifyTactic,
			e.NotifierNotifyPerDuration,
			e.NotifierNotifyDuration,
			e.NotifierNotifyQuorum,
		}
	}
	return StrategyOptions{}
}

// GetSemaphoreOptions returns all the selection options required to run a
// selection in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetSemaphoreOptions(t Type) StrategyOptions {
	switch t {
	case Consul:
		return StrategyOptions{
			Strategy:            e.ConsulSemaphoreStrategy,
			Tactic:              e.ConsulSemaphoreTactic,
			RequestsPerDuration: e.ConsulSemaphorePerDuration,
			RequestsDuration:    e.ConsulSemaphoreDuration,
		}
	}
	return StrategyOptions{}
}

// GetHeartbeatOptions returns all the selection options required to run a
// selection in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetHeartbeatOptions(t Type) StrategyOptions {
	switch t {
	case Consul:
		return StrategyOptions{
			Strategy:            e.ConsulHeartbeatStrategy,
			Tactic:              e.ConsulHeartbeatTactic,
			RequestsPerDuration: e.ConsulHeartbeatPerDuration,
			RequestsDuration:    e.ConsulHeartbeatDuration,
		}
	}
	return StrategyOptions{}
}

// GetKeyStoreOptions returns all the selection options required to run a
// selection in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetKeyStoreOptions(t Type) StrategyOptions {
	switch t {
	case Consul:
		return StrategyOptions{
			Strategy:            e.ConsulKeyStoreStrategy,
			Tactic:              e.ConsulKeyStoreTactic,
			RequestsPerDuration: e.ConsulKeyStorePerDuration,
			RequestsDuration:    e.ConsulKeyStoreDuration,
		}
	}
	return StrategyOptions{}
}
