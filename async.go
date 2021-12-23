package mediator

import "context"

type Async interface {
	Await() (interface{}, error)
}

type async struct {
	ctx   context.Context
	await func(ctx context.Context) (interface{}, error)
}

func (a *async) Await() (interface{}, error) {
	return a.await(a.ctx)
}
