package promise

import (
	"context"
	"errors"
	"testing"
)

func TestPromise_Wait_RejectAfterResolve(t *testing.T) {
	p := NewPromise[string]()

	p.Resolve("5")
	p.Reject(errors.New("error"))

	resp := p.Wait(context.Background())

	if resp.Payload != "5" {
		t.Error("Wait should return response with 5 payload, got: " + resp.Payload)
	}
}

func TestPromise_Wait_ResolveAfterReject(t *testing.T) {
	p := NewPromise[string]()

	p.Reject(errors.New("error"))
	p.Resolve("5")

	resp := p.Wait(context.Background())

	if resp.Error.Error() != "error" {
		t.Error("Wait should return response with error")
	}
}

func TestPromise_Wait_CanceledByCtx(t *testing.T) {
	p := NewPromise[string]()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := p.Wait(ctx)

	if !errors.Is(resp.Error, ErrCanceled) {
		t.Error("Wait should return response with ErrCanceled error, got: " + resp.Error.Error())
	}
}

func TestPromise_ApplyResolveHandler(t *testing.T) {
	appliedByConstructor := false
	appliedByApplyMethod := false

	p := NewPromise[string](WithResolveHandler(func(payload string) {
		appliedByConstructor = true
	}))
	p.Apply(WithResolveHandler(func(payload string) {
		appliedByApplyMethod = true
	}))

	p.Resolve("")

	if !appliedByConstructor {
		t.Error("Resolve handler was not applied by constructor")
	}

	if !appliedByApplyMethod {
		t.Error("Resolve handler was not applied by Apply method")
	}
}

func TestPromise_ApplyRejectHandler(t *testing.T) {
	appliedByConstructor := false
	appliedByApplyMethod := false

	p := NewPromise[string](WithRejectHandler[string](func(err error) {
		appliedByConstructor = true
	}))
	p.Apply(WithRejectHandler[string](func(err error) {
		appliedByApplyMethod = true
	}))

	p.Reject(errors.New("error"))

	if !appliedByConstructor {
		t.Error("Reject handler was not applied by constructor")
	}

	if !appliedByApplyMethod {
		t.Error("Reject handler was not applied by Apply method")
	}
}
