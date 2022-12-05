package promise

import (
	"context"
	"errors"
)

var ErrCanceled = errors.New("promise: context canceled")

type Response[T any] struct {
	Payload T
	Error   error
}

type Promise[T any] struct {
	doneCtx    context.Context
	doneNotify context.CancelFunc
	response   *Response[T]
}

func NewPromise[T any]() *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())

	return &Promise[T]{
		doneCtx:    ctx,
		doneNotify: cancel,
	}
}

func (f *Promise[T]) Resolve(payload T) {
	f.response = &Response[T]{Payload: payload}
	f.doneNotify()
}

func (f *Promise[T]) Reject(err error) {
	f.response = &Response[T]{Error: err}
	f.doneNotify()
}

func (f *Promise[T]) Wait(ctx context.Context) Response[T] {
	select {
	case <-ctx.Done():
		if f.response == nil {
			return Response[T]{
				Error: ErrCanceled,
			}
		}
	case <-f.doneCtx.Done():
	}

	return *f.response
}
