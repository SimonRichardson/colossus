package env

import (
	c "github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/spf13/viper"
)

type Env struct {
	source *viper.Viper

	C *c.Env

	// Shim

	ShimHttpAddress string

	// Score

	ScoreMaxSize             int
	ScorePoolRoutingStrategy string

	ScoreIncrementStrategy    string
	ScoreIncrementTactic      string
	ScoreIncrementPerDuration int
	ScoreIncrementDuration    string
	ScoreIncrementQuorum      float64

	// Redis

	ShimRedisInstances      string
	ShimRedisConnectTimeout string
	ShimRedisReadTimeout    string
	ShimRedisWriteTimeout   string
}

// Type describes what stragegy options are available
type Type int

const (
	Score = iota
)

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

	v.SetDefault("shim_http_address", ":9202")

	v.SetDefault("score_max_size", 100)
	v.SetDefault("score_pool_routing_strategy", "Hash")

	v.SetDefault("score_increment_strategy", "IncrementFromTime")
	v.SetDefault("score_increment_tactic", "NonBlocking")
	v.SetDefault("score_increment_per_duration", 0)
	v.SetDefault("score_increment_duration", 0)
	v.SetDefault("score_increment_quorum", 0.51)

	v.SetDefault("shim_redis_instances", "tcp://notifier:6379")
	v.SetDefault("shim_redis_connect_timeout", "1m")
	v.SetDefault("shim_redis_read_timeout", "30s")
	v.SetDefault("shim_redis_write_timeout", "30s")

	e := &Env{
		source: v,
		C:      c.New(paths),
	}

	e.read()

	return e
}

func (e *Env) read() {
	e.ShimHttpAddress = e.source.GetString("shim_http_address")

	e.ScoreMaxSize = e.source.GetInt("score_max_size")
	e.ScorePoolRoutingStrategy = e.source.GetString("score_pool_routing_strategy")

	e.ScoreIncrementStrategy = e.source.GetString("score_increment_strategy")
	e.ScoreIncrementTactic = e.source.GetString("score_increment_tactic")
	e.ScoreIncrementPerDuration = e.source.GetInt("score_increment_per_duration")
	e.ScoreIncrementDuration = e.source.GetString("score_increment_duration")
	e.ScoreIncrementQuorum = e.source.GetFloat64("score_increment_quorum")

	e.ShimRedisInstances = e.source.GetString("shim_redis_instances")
	e.ShimRedisConnectTimeout = e.source.GetString("shim_redis_connect_timeout")
	e.ShimRedisReadTimeout = e.source.GetString("shim_redis_read_timeout")
	e.ShimRedisWriteTimeout = e.source.GetString("shim_redis_write_timeout")
}

// GetIncrementOptions returns all the increments options required to run a
// increments in the application. It takes a Type argument to switch over the
// storage strategy.
func (e *Env) GetIncrementOptions(t Type) c.StrategyOptions {
	switch t {
	case Score:
		return c.StrategyOptions{
			e.ScoreIncrementStrategy,
			e.ScoreIncrementTactic,
			e.ScoreIncrementPerDuration,
			e.ScoreIncrementDuration,
			e.ScoreIncrementQuorum,
		}
	}
	return c.StrategyOptions{}
}
