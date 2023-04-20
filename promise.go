package promise

import (
	"context"
	"errors"
	"github.com/semichkin-gopkg/conf"
)

var ErrCanceled = errors.New("promise: canceled by ctx.Done")

type Response[T any] struct {
	Payload T
	Error   error
}

type Promise[T any] struct {
	config Conf[T]

	doneCtx    context.Context
	notifyDone context.CancelFunc
	response   *Response[T]
}

func NewPromise[T any](
	updaters ...conf.Updater[Conf[T]],
) *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())

	return &Promise[T]{
		config: conf.New[Conf[T]]().
			Append(updaters...).
			Build(),
		doneCtx:    ctx,
		notifyDone: cancel,
	}
}

func (p *Promise[T]) Apply(
	updaters ...conf.Updater[Conf[T]],
) {
	p.config = conf.New[Conf[T]]().
		Append(func(c *Conf[T]) {
			*c = p.config
		}).
		Append(updaters...).
		Build()
}

func (p *Promise[T]) Resolve(payload T) {
	if p.response != nil {
		return
	}

	for _, h := range p.config.ResolveHandlers {
		h(payload)
	}

	p.response = &Response[T]{Payload: payload}
	p.notifyDone()
}

func (p *Promise[T]) Reject(err error) {
	if p.response != nil {
		return
	}

	for _, h := range p.config.RejectHandlers {
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
