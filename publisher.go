package mediator

import (
	"context"
	"fmt"
)

type Event interface{}
type EventHandler func(ctx context.Context, event Event) error

type Publisher interface {
	Publish(ctx context.Context, event Event) error
	PublishAsync(ctx context.Context, event Event) Async
}

type publisher struct {
	register Register
}

func newPublisher(register Register) Publisher {

	if register == nil {
		panic(fmt.Errorf("register value cannot be nil"))
	}

	return &publisher{
		register: register,
	}
}

func (p *publisher) Publish(ctx context.Context, event Event) error {
	if fns, ok := p.register.GetEventHandler(event); ok {
		for _, fn := range fns {
			if err := fn(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("event handler not found")
}

func (p *publisher) PublishAsync(ctx context.Context, event Event) Async {
	var (
		err      error
		selector = make(chan struct{})
	)

	go func() {
		defer close(selector)
		if fns, ok := p.register.GetEventHandler(event); ok {
			for _, fn := range fns {
				if err = fn(ctx, event); err != nil {
					break
				}
			}
		} else {
			err = fmt.Errorf("event handler not found")
		}
	}()

	return &async{
		ctx: ctx,
		await: func(ctx context.Context) (interface{}, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-selector:
				return nil, err
			}
		},
	}
}
