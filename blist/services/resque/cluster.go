package resque

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/typex"
	r "github.com/garyburd/redigo/redis"
)

const (
	defaultKey             string        = "uE7Lz8Swq0iDI51XtO6cfnOM1DIOMO8IKZPBjJmO"
	defaultBlockingTimeout time.Duration = time.Second * 10
)

// Cluster defines a interface for parallelised requesting.
type Cluster interface {
	sv.Enqueuer
	sv.Register
}

type cluster struct {
	pool *redis.Pool
}

func newCluster(pool *redis.Pool) *cluster {
	return &cluster{pool}
}

func (c *cluster) EnqueueBytes(queue selectors.Queue, class selectors.Class, value []byte) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		err := enqueueBytes(conn, queue, class, value)
		dst <- sv.NewErrorElement(err)
	})
}

func (c *cluster) DequeueBytes(queue selectors.Queue, class selectors.Class) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		if res, err := dequeueBytes(conn, queue, class); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewBytesElement(res)
		}
	})
}

func (c *cluster) RegisterFailure(queue selectors.Queue,
	class selectors.Class,
	failure selectors.Failure,
) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		err := registerFailure(conn, queue, class, failure)
		dst <- sv.NewErrorElement(err)
	})
}

func (c *cluster) common(fn func(r.Conn, chan sv.Element)) <-chan sv.Element {
	out := make(chan sv.Element)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() { wg.Wait(); close(out) }()
		go func() {
			defer wg.Done()

			if err := c.pool.With(defaultKey, func(conn r.Conn) error {
				fn(conn, out)
				return nil
			}); err != nil {
				out <- sv.NewErrorElement(err)
			}
		}()
	}()
	return out
}

var (
	headerPrefix = []byte{123, 34, 99, 108, 97, 115, 115, 34, 58, 34}
	headerSuffix = []byte{34, 44, 34, 97, 114, 103, 115, 34, 58}
	headerLen    = len(headerPrefix) + len(headerSuffix)
	footerSuffix = []byte{125}
	footerLen    = len(footerSuffix)
)

func header(class string) []byte {
	return append(headerPrefix, append([]byte(class), headerSuffix...)...)
}

func footer() []byte {
	return footerSuffix
}

func encode(class string, bytes []byte) []byte {
	return append(header(class), append(bytes, footer()...)...)
}

func decode(bytes []byte) (string, []byte) {
	// class
	var class []byte
	for i := len(headerPrefix); i < len(bytes); i++ {
		if match(headerSuffix, bytes[i:]) {
			break
		}
		class = append(class, bytes[i])
	}
	var (
		start = headerLen + len(class)
		end   = len(bytes) - footerLen
	)
	return string(class), bytes[start:end]
}

func match(a, b []byte) bool {
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}

func enqueueBytes(conn r.Conn, queue selectors.Queue, class selectors.Class, value []byte) error {
	bytes := encode(class.String(), value)
	_, err := conn.Do("RPUSH", fmt.Sprintf("resque:queue:%s", queue), bytes)
	return err
}

func dequeueBytes(conn r.Conn, queue selectors.Queue, class selectors.Class) ([]byte, error) {
	var (
		channel  = fmt.Sprintf("resque:queue:%s", queue)
		res, err = r.Values(conn.Do("BLPOP", channel, defaultBlockingTimeout.Seconds()))
	)
	if err != nil {
		return nil, err
	}
	if len(res) < 2 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Invalid byte values.")
	}

	if bytes, ok := res[1].([]byte); ok {
		c, value := decode(bytes)
		if c != class.String() {
			return value, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Invalid class found (expected: %q, actual: %q).", c, class)
		}
		return value, nil
	}

	return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
		"Invalid types found (%q, %q)", queue, class)
}

func registerFailure(conn r.Conn,
	queue selectors.Queue,
	class selectors.Class,
	failure selectors.Failure,
) error {
	var (
		channel   = fmt.Sprintf("resque:failed")
		data, err = json.Marshal(failurePayload{
			FailedAt:  time.Now(),
			Exception: "Error",
			Error:     failure.Error.Error(),
			Queue:     queue.String(),
			Class:     class.String(),
		})
	)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", channel, data)
	return err
}

type failurePayload struct {
	FailedAt  time.Time `json:"failed_at"`
	Exception string    `json:"exception"`
	Error     string    `json:"error"`
	Queue     string    `json:"queue"`
	Class     string    `json:"class"`
}
