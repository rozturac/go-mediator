package mediator

import (
	"context"
	"fmt"
	"testing"
)

func TestRegister_RegisterEvent(t *testing.T) {
	type TestEvent struct {
		name string
	}
	event := &TestEvent{name: "TestEvent"}
	expectedHandler := func(ctx context.Context, event Event) error {
		return nil
	}
	register := &register{
		registeredEvents: map[string][]EventHandler{},
	}

	if err := register.RegisterEvent(event, expectedHandler); err != nil {
		t.Errorf("register.RegisterEvent(%v, eventHandler) = %v, expected = <nil>", event, err)
	}

	if _, ok := register.GetEventHandler(event); !ok {
		t.Errorf("register.GetEventHandler(%v)= _, false expected _, true", event)
	}
}

func TestRegister_RegisterCommand(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedHandler := func(ctx context.Context, command Command) (interface{}, error) {
		return nil, nil
	}

	register := newRegister()
	if err := register.RegisterCommand(command, expectedHandler); err != nil {
		t.Errorf("err := register.RegisterCommand(%v, eventHandler), expected: err = <nil>, received: %v", command, err)
	}

	if _, ok := register.GetCommandHandler(command); !ok {
		t.Errorf("_, ok := register.GetCommandHandler(%v), expected: ok = true, received: ok = false", command)
	}
}

func TestRegister_GetCommandHandler(t *testing.T) {
	type TestCommand struct {
		name string
	}
	command := &TestCommand{name: "TestCommand"}
	expectedResult := fmt.Sprintf("Called Command: %v", command)
	expectedHandler := func(ctx context.Context, command Command) (interface{}, error) {
		return fmt.Sprintf("Called Command: %v", command), nil
	}

	register := newRegister()
	if err := register.RegisterCommand(command, expectedHandler); err != nil {
		t.Errorf("err := register.RegisterCommand(%v, eventHandler), expected: err = <nil>, received: %v", command, err)
	}

	if handler, ok := register.GetCommandHandler(command); ok {
		if result, err := handler(context.Background(), command); err != nil || result != expectedResult {
			t.Errorf("result, err := handler(ctx, %v), expected: (result = %v, err = <nil>), received: (result = %v, err = %v)", command, expectedResult, result, err)
		}
	}
}

func TestRegister_GetEventHandler(t *testing.T) {
	type TestEvent struct {
		name string
	}
	event := &TestEvent{name: "TestEvent"}
	expectedError := fmt.Errorf("test error")
	expectedHandler := func(ctx context.Context, event Event) error {
		return expectedError
	}

	register := newRegister()
	if err := register.RegisterEvent(event, expectedHandler, expectedHandler); err != nil {
		t.Errorf("register.RegisterEvent(%v, eventHandler) = %v, expected = <nil>", event, err)
	}

	if handlers, ok := register.GetEventHandler(event); ok {
		for _, handler := range handlers {
			if err := handler(context.Background(), event); err != expectedError {
				t.Errorf("Expected err: %v, received: %v", expectedError, err)
			}
		}
	}
}
