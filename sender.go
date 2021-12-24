package mediator

import (
	"context"
	"fmt"
)

type Command interface{}
type CommandHandler func(ctx context.Context, command Command) (interface{}, error)

type Sender interface {
	Send(ctx context.Context, command Command) (interface{}, error)
	SendAsync(ctx context.Context, command Command) Async
}

type sender struct {
	register Register
	behavior Behavior
}

func newSender(register Register, behavior Behavior) Sender {

	if register == nil {
		panic(fmt.Errorf("register value cannot be nil"))
	}

	if behavior == nil {
		panic(fmt.Errorf("behavior value cannot be nil"))
	}

	return &sender{
		register: register,
		behavior: behavior,
	}
}

func (s *sender) Send(ctx context.Context, command Command) (result interface{}, err error) {
	if fn, ok := s.register.GetCommandHandler(command); ok {
		result, err = s.behavior.use(ctx, command, fn, 0)
	} else {
		err = fmt.Errorf("not found handler registered by %v", command)
	}
	return
}

func (s *sender) SendAsync(ctx context.Context, command Command) Async {
	var (
		selector = make(chan struct{})
		err      error
		result   interface{}
	)
	go func() {
		defer close(selector)
		if fn, ok := s.register.GetCommandHandler(command); ok {
			result, err = s.behavior.use(ctx, command, fn, 0)
		} else {
			err = fmt.Errorf("not found handler registered by %v", command)
		}
	}()
	return &async{
		ctx: ctx,
		await: func(ctx context.Context) (interface{}, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-selector:
				return result, err
			}
		},
	}
}
