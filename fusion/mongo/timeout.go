package mongo

import "time"

type ConnectionTimeout struct {
	global time.Duration
}

func newConnectionTimeout() *ConnectionTimeout {
	return &ConnectionTimeout{
		global: time.Minute,
	}
}

func (c *ConnectionTimeout) All(timeout time.Duration) *ConnectionTimeout {
	c.global = timeout
	return c
}
