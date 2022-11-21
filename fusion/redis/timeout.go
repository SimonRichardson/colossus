package redis

import "time"

type ConnectionTimeout struct {
	connect time.Duration
	read    time.Duration
	write   time.Duration
}

func newConnectionTimeout() *ConnectionTimeout {
	return &ConnectionTimeout{
		connect: time.Minute,
		read:    time.Second * 10,
		write:   time.Second * 30,
	}
}

func (c *ConnectionTimeout) All(timeout time.Duration) *ConnectionTimeout {
	c.connect = timeout
	c.read = timeout
	c.write = timeout
	return c
}
