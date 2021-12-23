package mediator

import (
	"context"
	"errors"
	"testing"
)

func TestPublisher_newSender(t *testing.T) {
	sender := newSender(newRegister(), newBehavior())
	if sender == nil {
		t.Errorf("sender := newSender(newRegister(), newBehavior()), expected: (sender = &sender{}), received: (sender = <nil>)")
	}
}

func TestPublisher_newSenderWithoutRegister(t *testing.T) {
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

func TestPublisher_Publish(t *testing.T) {
	type TestEvent struct {
		name string
	}
	event := &TestEvent{name: "TestEvent"}
	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	register := newRegister()
	register.RegisterEvent(&TestEvent{}, handler)
	publisher := newPublisher(register)

	if err := publisher.Publish(context.Background(), event); err != nil {
		t.Errorf("Expected error: <nil>, received: %v", err)
	}
}

func TestPublisher_PublishAsync(t *testing.T) {
	type TestEvent struct {
		name string
	}
	event := &TestEvent{name: "TestEvent"}
	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	register := newRegister()
	register.RegisterEvent(&TestEvent{}, handler)
	publisher := newPublisher(register)
	async := publisher.PublishAsync(context.Background(), event)
	if err, _ := async.Await(); err != nil {
		t.Errorf("Expected error: <nil>, received: %v", err)
	}
}
