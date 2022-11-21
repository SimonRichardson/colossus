package agents

import (
	"fmt"
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-walker/common"
	"github.com/SimonRichardson/colossus/teleprinter"
)

const (
	defaultWalkerNamespace = selectors.Namespace("colossus_walker")
	defaultMaxSize         = 99999
	defaultExpiry          = time.Hour * 24 * 7 * 30
)

type Walk struct{}

func (a Walk) Init(opts AgentOptions) error {
	var (
		co    = opts.Coordinator
		timer = time.NewTicker(time.Minute)
	)

	go func() {

	loop:
		// Walk over all the keys!
		for {

			select {
			case <-timer.C:
				keys, err := co.Keys()
				if err != nil {
					continue loop
				}

				// We should do something clever like only walk some of them at
				// a time.

				for _, key := range keys {
					func() {
						var (
							ns          = key.Namespace()
							unlock, err = co.Lock(ns.Prefix(defaultWalkerNamespace))
						)
						if err != nil {
							teleprinter.L.Info().Printf("Unable to process repair walker, as event is locked : %s\n", err)
							return
						}
						defer unlock()

						if _, err := common.Get(fmt.Sprintf("%s/http/v1/%s/select?limit=99999&size=%d&expiry=%d",
							opts.HttpAddress,
							key.String(),
							defaultMaxSize,
							defaultExpiry.Nanoseconds(),
						)); err != nil {
							teleprinter.L.Error().Printf("Error processing all events with : %s\n", err)
							return
						}
					}()
				}
			}
		}
	}()

	return nil
}
