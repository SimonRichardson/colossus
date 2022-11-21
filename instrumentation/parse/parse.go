package parse

import (
	"io"
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	p "github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/SimonRichardson/colossus/instrumentation/multi"
	"github.com/SimonRichardson/colossus/instrumentation/noop"
	"github.com/SimonRichardson/colossus/instrumentation/plaintext"
	"github.com/SimonRichardson/colossus/instrumentation/prometheus"
	r "github.com/SimonRichardson/colossus/instrumentation/redis"
	"github.com/SimonRichardson/colossus/instrumentation/statsd"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/peterbourgon/g2s"
)

type InstrumentationOptions struct {
	StatsdAddress       string
	StatsdSampleRate    float32
	PlaintextWriter     io.Writer
	RedisAddress        string
	RedisBufferDuration time.Duration
	RedisTimeout        string
}

func ParseString(value string,
	options InstrumentationOptions,
) (instrumentation.Instrumentation, error) {
	parts := strings.Split(value, ";")
	switch common.StripWhitespace(strings.ToLower(parts[0])) {
	case "noop":
		return noop.New(), nil
	case "plaintext":
		return plaintext.New(options.PlaintextWriter), nil
	case "statsd":
		statter := g2s.Noop()
		if options.StatsdAddress != "" {
			var err error
			if statter, err = g2s.Dial("udp", options.StatsdAddress); err != nil {
				typex.Fatal(err)
			}
		}
		return statsd.New(statter, options.StatsdSampleRate), nil
	case "prometheus":
		return prometheus.New("colossus", time.Second*10), nil
	case "redis":
		host := options.RedisAddress
		if err := p.ValidRedisHost(host); err != nil {
			return nil, err
		}

		timeout := options.RedisTimeout
		connTimeout, routing, err := p.Parse(timeout, timeout, timeout, "hash", nil)
		if err != nil {
			return nil, err
		}

		return r.New(p.New(
			[]string{host},
			routing,
			connTimeout,
			100,
			nil,
		), options.RedisBufferDuration), nil
	case "multi":
		instruments := []instrumentation.Instrumentation{}
		for _, v := range parts[1:] {
			if instr, err := ParseString(v, options); err != nil {
				return noop.New(), err
			} else {
				instruments = append(instruments, instr)
			}
		}
		return multi.New(instruments...), nil
	}
	return noop.New(), typex.Errorf(errors.Source, errors.NoCaseFound, "Invalid instrumentation %q", value)
}
