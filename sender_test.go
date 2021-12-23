package mediator

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestSender_newSender(t *testing.T) {
	sender := newSender(newRegister(), newBehavior())
	if sender == nil {
		t.Errorf("sender := newSender(newRegister(), newBehavior()), expected: (sender = &sender{}), received: (sender = <nil>)")
	}
}

func TestSender_newSenderWithoutRegister(t *testing.T) {
	expectedError := errors.New("register value cannot be nil")
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if err.Error() == expectedError.Error() {
					return
				}
			}
			t.Errorf("Expected error: %v, received: %v", expectedError, r)
		}
		t.Errorf("Expected error: %v, received: %v", expectedError, nil)
	}()

	newSender(nil, newBehavior())
}

func TestSender_newSenderWithoutBehavior(t *testing.T) {
	expectedError := errors.New("behavior value cannot be nil")
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if err.Error() == expectedError.Error() {
					return
				}
			}
			t.Errorf("Expected error: %v, received: %v", expectedError, r)
		}
		t.Errorf("Expected error: %v, received: %v", expectedError, nil)
	}()

	newSender(newRegister(), nil)
}

func TestSender_Send(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedResult := fmt.Sprintf("Called Command: %v", command)
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v", command), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	sender := newSender(register, newBehavior())
	if result, err := sender.Send(context.Background(), command); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}

func TestSender_SendWithCorrelationId(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	ctx := context.WithValue(context.Background(), "correlationId", uuid.NewString())
	expectedResult := fmt.Sprintf("Called Command: %v, with correlationId: %v", command, ctx.Value("correlationId"))
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v, with correlationId: %v", command, ctx.Value("correlationId")), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	sender := newSender(register, newBehavior())
	if result, err := sender.Send(ctx, command); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}

func TestSender_SendWithBehavior(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedResult := fmt.Sprintf("WrappedResult: Called Command: %v", command)
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v", command), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	behavior := newBehavior()
	behavior.WithBehavior(func(ctx context.Context, command Command, next CommandHandler) (interface{}, error) {
		result, err := next(ctx, command)
		return fmt.Sprintf("WrappedResult: %v", result), err
	})
	sender := newSender(register, behavior)
	if result, err := sender.Send(context.Background(), command); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}

func TestSender_SendAsync(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedResult := fmt.Sprintf("Called Command: %v", command)
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v", command), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	sender := newSender(register, newBehavior())
	async := sender.SendAsync(context.Background(), command)

	if result, err := async.Await(); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}

func TestSender_SendAsyncWithCorrelationId(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	ctx := context.WithValue(context.Background(), "correlationId", uuid.NewString())
	expectedResult := fmt.Sprintf("Called Command: %v, with correlationId: %v", command, ctx.Value("correlationId"))
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v, with correlationId: %v", command, ctx.Value("correlationId")), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	sender := newSender(register, newBehavior())
	async := sender.SendAsync(ctx, command)

	if result, err := async.Await(); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}

func TestSender_SendAsyncWithBehavior(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedResult := fmt.Sprintf("WrappedResult: Called Command: %v", command)
	handler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v", command), nil
	}

	register := newRegister()
	register.RegisterCommand(&TestCommand{}, handler)
	behavior := newBehavior()
	behavior.WithBehavior(func(ctx context.Context, command Command, next CommandHandler) (interface{}, error) {
		result, err := next(ctx, command)
		return fmt.Sprintf("WrappedResult: %v", result), err
	})
	sender := newSender(register, behavior)
	async := sender.SendAsync(context.Background(), command)

	if result, err := async.Await(); err != nil || result != expectedResult {
		t.Errorf("result, err :=sender.Send(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
	}
}
