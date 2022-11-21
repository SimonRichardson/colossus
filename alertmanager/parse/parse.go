package parse

import (
	"os"
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
	"github.com/SimonRichardson/colossus/alertmanager/multi"
	"github.com/SimonRichardson/colossus/alertmanager/noop"
	"github.com/SimonRichardson/colossus/alertmanager/plaintext"
	"github.com/SimonRichardson/colossus/alertmanager/prometheus"
	"github.com/SimonRichardson/colossus/alertmanager/statsd"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/peterbourgon/g2s"
)

type AlertManagerOptions struct {
	StatsdAddress    string
	StatsdSampleRate float32
}

func ParseString(value string,
	options AlertManagerOptions,
) (alertmanager.AlertManager, error) {
	parts := strings.Split(value, ";")
	switch common.StripWhitespace(strings.ToLower(parts[0])) {
	case "noop":
		return noop.New(), nil
	case "plaintext":
		return plaintext.New(os.Stderr), nil
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
		return prometheus.New("bombe", time.Second*10), nil
	case "multi":

		managers := []alertmanager.AlertManager{}
		for _, v := range parts[1:] {
			if instr, err := ParseString(v, options); err != nil {
				return noop.New(), err
			} else {
				managers = append(managers, instr)
			}
		}
		return multi.New(managers...), nil
	}
	return noop.New(), typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid alertmanager %q", value)
}
