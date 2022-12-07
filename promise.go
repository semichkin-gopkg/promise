package promise

import (
	"context"
	"errors"
	"github.com/semichkin-gopkg/configurator"
)

var ErrCanceled = errors.New("promise: canceled by ctx.Done")

type Response[T any] struct {
	Payload T
	Error   error
}

type Promise[T any] struct {
	configuration Configuration[T]

	doneCtx    context.Context
	notifyDone context.CancelFunc
	response   *Response[T]
}

func NewPromise[T any](
	updaters ...configurator.Updater[Configuration[T]],
) *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())

	return &Promise[T]{
		configuration: configurator.New[Configuration[T]]().
			Append(updaters...).
			Apply(),
		doneCtx:    ctx,
		notifyDone: cancel,
	}
}

func (p *Promise[T]) Apply(
	updaters ...configurator.Updater[Configuration[T]],
) {
	p.configuration = configurator.New[Configuration[T]]().
		Append(func(c *Configuration[T]) {
			*c = p.configuration
		}).
		Append(updaters...).
		Apply()
}

func (p *Promise[T]) Resolve(payload T) {
	if p.response != nil {
		return
	}

	for _, h := range p.configuration.ResolveHandlers {
		h(payload)
	}

	p.response = &Response[T]{Payload: payload}
	p.notifyDone()
}

func (p *Promise[T]) Reject(err error) {
	if p.response != nil {
		return
	}

	for _, h := range p.configuration.RejectHandlers {
		h(err)
	}

	p.response = &Response[T]{Error: err}
	p.notifyDone()
}

func (p *Promise[T]) Wait(ctx context.Context) Response[T] {
	select {
	case <-ctx.Done():
		p.Reject(ErrCanceled)
	case <-p.doneCtx.Done():
	}

	return *p.response
}
