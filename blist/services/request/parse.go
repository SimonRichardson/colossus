package request

import (
	"net/url"
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/request/client"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	empty = []Cluster{}
)

// DefaultConfig creates a default configuration for the service
func DefaultConfig(address string) ([]Cluster, error) {
	return ParseString(address, 100, nil)
}

// StrategyOptions defines what options are available when creating a and using
// a stragegy.
type StrategyOptions struct {
	Strategy            string
	Tactic              string
	RequestsPerDuration int
	RequestsDuration    string
	Quorum              float64
}

// ParseString parses various inputs and returns a slice of clusters that we can
// then use with in the service.
func ParseString(address string,
	maxSize int,
	creator client.ClientCreator,
) ([]Cluster, error) {
	var (
		clusters = []Cluster{}
		fn       = clientCreator(creator)
	)

	for k, v := range strings.Split(common.StripWhitespace(address), ";") {
		uri, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, newCluster(
			client.New(*uri, k, maxSize, fn),
			k,
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q for request", address)
	}

	return clusters, nil
}

func clientCreator(creator client.ClientCreator) (res client.ClientCreator) {
	if creator == nil {
		res = func(address url.URL, index int) client.Client {
			return client.NewClient(address, index)
		}
	} else {
		res = creator
	}
	return
}

type requestStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.Request
	Tactic   Tactic
}

func (o requestStategyOpts) Apply(s *Service) selectors.Request { return o.Strategy(s, o.Tactic) }

func ParseRequestStrategy(opts StrategyOptions) (requestStategyOpts, error) {
	var (
		strategy requestStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseRequestStrategy(opts.Strategy); err != nil {
		return requestStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return requestStategyOpts{}, err
	}

	return requestStategyOpts{strategy, tactic}, nil
}

func parseRequestStrategy(strategy string) (requestStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopRequest, nil
	case "request":
		return Request, nil
	}
	return NoopRequest, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid request strategy %q", strategy)
}

func readTactic(tactic string,
	requestsPerDuration int,
	requestsDuration string,
) (Tactic, error) {
	dur, err := time.ParseDuration(requestsDuration)
	if err != nil {
		return noopTactic, err
	}

	switch common.Normalise(tactic) {
	case "noop":
		return noopTactic, nil
	case "nonblocking":
		return nonBlocking, nil
	case "ratelimited":
		return rateLimited(requestsPerDuration, dur), nil
	}
	return noopTactic, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid request tactic %q", tactic)
}

func ParseFittingProvider(provider string) (Provider, error) {
	switch common.Normalise(provider) {
	case "first":
		return First(), nil
	case "race":
		return Race(), nil
	case "ignore":
		return Ignore(), nil
	case "random":
		return Random(0, 0.5), nil
	case "ratelimit":
		return RateLimit(0, time.Millisecond), nil
	}
	return Ignore(), typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid request fitting %q", provider)
}
