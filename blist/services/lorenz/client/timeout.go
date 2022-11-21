package client

import "time"

type ConnectionTimeout struct {
	Global time.Duration
}

func NewConnectionTimeout() *ConnectionTimeout {
	return &ConnectionTimeout{
		Global: time.Minute,
	}
}

func (c *ConnectionTimeout) All(timeout time.Duration) *ConnectionTimeout {
	c.Global = timeout
	return c
}
