package redis

import "github.com/SimonRichardson/colossus/fusion/tests/stubs"

type Conn struct {
	*stubs.Actions
}

func New(actions *stubs.Actions) *Conn {
	return &Conn{actions}
}

func (c *Conn) Send(commandName string, args ...interface{}) error {
	_, err := c.Run("Send", func(action stubs.Action) (interface{}, error) {
		return action.Run(append([]interface{}{commandName}, args...)...)
	})
	return err
}

func (c *Conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	return c.Run("Do", func(action stubs.Action) (interface{}, error) {
		return action.Run(append([]interface{}{commandName}, args...)...)
	})
}

func (c *Conn) Flush() error {
	_, err := c.Run("Flush", func(action stubs.Action) (interface{}, error) {
		return action.Run()
	})
	return err
}

func (c *Conn) Receive() (interface{}, error) {
	return c.Run("Receive", func(action stubs.Action) (interface{}, error) {
		return action.Run()
	})
}

func (c *Conn) Close() error {
	_, err := c.Run("Close", func(action stubs.Action) (interface{}, error) {
		return action.Run()
	})
	return err
}

func (c *Conn) Err() error {
	if _, err := c.Run("Err", func(action stubs.Action) (interface{}, error) {
		return action.Run()
	}); err != nil {
		return err.(error)
	}
	return nil
}
