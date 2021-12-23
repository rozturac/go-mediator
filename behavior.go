package mediator

import "context"

type BehaviorHandler func(ctx context.Context, command Command, next CommandHandler) (interface{}, error)

type Behavior interface {
	WithBehavior(handler BehaviorHandler)
	use(ctx context.Context, command Command, handler CommandHandler, index int) (interface{}, error)
}

type behavior struct {
	behaviorHandlers []BehaviorHandler
}

func newBehavior() *behavior {
	return &behavior{
		behaviorHandlers: []BehaviorHandler{},
	}
}

func (b *behavior) WithBehavior(handler BehaviorHandler) {
	b.behaviorHandlers = append(b.behaviorHandlers, handler)
}

func (b *behavior) use(ctx context.Context, command Command, handler CommandHandler, index int) (interface{}, error) {
	if index >= len(b.behaviorHandlers) {
		return handler(ctx, command)
	}

	return b.behaviorHandlers[index](ctx, command, func(ctx context.Context, command Command) (interface{}, error) {
		return b.use(ctx, command, handler, index+1)
	})
}
