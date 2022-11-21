package lorenz

import (
	"encoding/json"
	"sync"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/blist/services/lorenz/client"
	"github.com/SimonRichardson/colossus/typex"
)

// Cluster defines what
type Cluster interface {
	sv.Charger
	sv.EventSelector
	sv.CodeSelector
	sv.Inspector
}

type cluster struct {
	client client.Client
}

func newCluster(client client.Client) *cluster {
	return &cluster{client}
}

func (c *cluster) Charge(event selectors.Event,
	user selectors.User,
	payment selectors.Payment,
) <-chan sv.Element {
	return c.common(func(dst chan sv.Element) {
		if result, err := charge(c.client, event, user, payment); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewKeyTxnElement(result.Key, result.Txn)
		}
	})
}

func (c *cluster) SelectEventByKey(key selectors.Key) <-chan sv.Element {
	return c.common(func(dst chan sv.Element) {
		if result, err := readEvent(c.client, key); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewEventElement(result)
		}
	})
}

func (c *cluster) SelectEventsByOffset(offset, limit int) <-chan sv.Element {
	return c.common(func(dst chan sv.Element) {
		if result, err := readEvents(c.client, offset, limit); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewEventsElement(result)
		}
	})
}

func (c *cluster) SelectCodeForEvent(event selectors.Event, user selectors.User) <-chan sv.Element {
	return c.common(func(dst chan sv.Element) {
		if result, err := readCode(c.client, event, user); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewCodeSetElement(result)
		}
	})
}

func (c *cluster) Version() <-chan sv.Element {
	return c.common(func(dst chan sv.Element) {
		if result, err := version(c.client); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewVersionElement(result)
		}
	})
}

func (c *cluster) common(fn func(dst chan sv.Element)) <-chan sv.Element {
	out := make(chan sv.Element)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() { wg.Wait(); close(out) }()
		go func() {
			defer wg.Done()

			fn(out)
		}()
	}()
	return out
}

type Error struct {
	Error       string `json:"error"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func readError(response *client.Response, status int) error {
	if response.StatusCode != status {
		var record Error
		if err := json.Unmarshal(response.Bytes, &record); err != nil {
			return typex.Errorf(errors.Source, errors.UnexpectedArgument,
				"Error occured (%s).", err.Error())
		}

		return typex.Errorf(errors.Source, errors.UnexpectedResults, record.Error)
	}
	return nil
}

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}
