package promise

import (
	"github.com/semichkin-gopkg/configurator"
)

type Configuration[T any] struct {
	ResolveHandlers []func(payload T)
	RejectHandlers  []func(err error)
}

func WithResolveHandler[T any](handler func(payload T)) configurator.Updater[Configuration[T]] {
	return func(c *Configuration[T]) {
		if handler != nil {
			c.ResolveHandlers = append(c.ResolveHandlers, handler)
		}
	}
}

func WithRejectHandler[T any](handler func(err error)) configurator.Updater[Configuration[T]] {
	return func(c *Configuration[T]) {
		if handler != nil {
			c.RejectHandlers = append(c.RejectHandlers, handler)
		}
	}
}
