package promise

import (
	"github.com/semichkin-gopkg/conf"
)

type Conf[T any] struct {
	ResolveHandlers []func(payload T)
	RejectHandlers  []func(err error)
}

func WithResolveHandler[T any](handler func(payload T)) conf.Updater[Conf[T]] {
	return func(c *Conf[T]) {
		if handler != nil {
			c.ResolveHandlers = append(c.ResolveHandlers, handler)
		}
	}
}

func WithRejectHandler[T any](handler func(err error)) conf.Updater[Conf[T]] {
	return func(c *Conf[T]) {
		if handler != nil {
			c.RejectHandlers = append(c.RejectHandlers, handler)
		}
	}
}
