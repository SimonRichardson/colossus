package client

import (
	"encoding/json"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/hashicorp/consul/api"
)

const (
	Passing     selectors.HealthStatus = "passing"
	Warning     selectors.HealthStatus = "warning"
	Critical    selectors.HealthStatus = "critical"
	Maintenance selectors.HealthStatus = "maintenance"
)

type ClientCreator func(string, string, string) Client

type Client interface {
	Lock(selectors.Namespace) (selectors.SemaphoreUnlock, error)
	Heartbeat(selectors.HealthStatus) error
	List(fs.Prefix) (map[string]int, error)
}

type client struct {
	api     *api.Client
	checkId string
	output  string
}

func NewClient(address, checkId, output string) Client {
	config := api.DefaultConfig()
	config.Address = address

	consul, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &client{
		api:     consul,
		checkId: checkId,
		output:  output,
	}
}

func (c *client) Lock(ns selectors.Namespace) (selectors.SemaphoreUnlock, error) {
	lock, err := c.api.LockKey(ns.Lock())
	if err != nil {
		return noop, err
	}

	ch, err := lock.Lock(nil)
	if err != nil {
		return noop, err
	}
	if ch == nil {
		return noop, typex.Errorf(errors.Source, errors.UnexpectedResults, "Lock not held.")
	}

	return func() error {
		return lock.Unlock()
	}, nil
}

func (c *client) Heartbeat(s selectors.HealthStatus) error {
	return c.api.Agent().UpdateTTL(c.checkId, c.output, status(s))
}

type KeyValueStoreElement struct {
	Key  string `json:"key"`
	Node int    `json:"node"`
}

func (c *client) List(p fs.Prefix) (map[string]int, error) {
	kv := c.api.KV()
	pairs, _, err := kv.List(p.String(), nil)
	if err != nil {
		return nil, err
	}

	res := make(map[string]int)
	for _, v := range pairs {
		var elm KeyValueStoreElement
		if err := json.Unmarshal(v.Value, &elm); err != nil {
			teleprinter.L.Info().Printf("Error decoding KV List %s\n", err.Error())
			continue
		}
		res[elm.Key] = elm.Node
	}
	return res, nil
}

func noop() error {
	return nil
}

func status(s selectors.HealthStatus) string {
	switch s {
	case Passing:
		return api.HealthPassing
	case Warning:
		return api.HealthWarning
	case Critical:
		return api.HealthCritical
	case Maintenance:
		return api.HealthMaint
	default:
		return api.HealthAny
	}
}
