package redis

import (
	"fmt"

	"github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/fusion/tests/stubs"
	r "github.com/garyburd/redigo/redis"
)

func SendAndReceive(
	fn1 func(commandName string, args ...interface{}) error,
	fn2 func() (interface{}, error),
) redis.RedisCreator {
	return func(address, password string, timeout *redis.ConnectionTimeout) (r.Conn, error) {
		actions := stubs.NewActions()
		actions.
			On("Send", func(args ...interface{}) (interface{}, error) {
				if m, ok := args[0].(string); ok {
					err := fn1(m, args[1:]...)
					return nil, err
				}
				return nil, fmt.Errorf("Unexpected arguments")
			}).
			On("Receive", func(args ...interface{}) (interface{}, error) {
				return fn2()
			}).
			On("Close", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}).
			On("Flush", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}).
			On("Err", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			})
		return New(actions), nil
	}
}

func DoSendAndReceive(
	fn1 func(commandName string, args ...interface{}) (interface{}, error),
	fn2 func(commandName string, args ...interface{}) error,
	fn3 func() (interface{}, error),
) redis.RedisCreator {
	return func(address, password string, timeout *redis.ConnectionTimeout) (r.Conn, error) {
		actions := stubs.NewActions()
		actions.
			On("Do", func(args ...interface{}) (interface{}, error) {
				if m, ok := args[0].(string); ok {
					return fn1(m, args[1:]...)
				}
				return nil, fmt.Errorf("Unexpected arguments")
			}).
			On("Send", func(args ...interface{}) (interface{}, error) {
				if m, ok := args[0].(string); ok {
					err := fn2(m, args[1:]...)
					return nil, err
				}
				return nil, fmt.Errorf("Unexpected arguments")
			}).
			On("Receive", func(args ...interface{}) (interface{}, error) {
				return fn3()
			}).
			On("Close", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}).
			On("Flush", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}).
			On("Err", func(args ...interface{}) (interface{}, error) {
				return nil, nil
			})
		return New(actions), nil
	}
}
