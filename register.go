package mediator

import (
	"errors"
	"fmt"
	"reflect"
)

type Register interface {
	RegisterEvent(event Event, handlers ...EventHandler) error
	RegisterCommand(command Command, handler CommandHandler) error
	GetEventHandler(event Event) ([]EventHandler, bool)
	GetCommandHandler(command Command) (CommandHandler, bool)
}

type register struct {
	registeredEvents   map[string][]EventHandler
	registeredCommands map[string]CommandHandler
}

func newRegister() Register {
	return &register{
		registeredCommands: map[string]CommandHandler{},
		registeredEvents:   map[string][]EventHandler{},
	}
}

func (r register) RegisterEvent(event Event, handlers ...EventHandler) error {
	t := reflect.TypeOf(event)
	if t.Elem().Kind() != reflect.Struct {
		return errors.New("event type must be struct")
	}

	for _, handler := range handlers {
		r.registeredEvents[t.Name()] = append(r.registeredEvents[t.Name()], handler)
	}
	return nil
}

func (r register) GetEventHandler(event Event) ([]EventHandler, bool) {
	t := reflect.TypeOf(event)
	if fns, ok := r.registeredEvents[t.Name()]; ok {
		return fns, true
	}
	return nil, false
}

func (r register) RegisterCommand(command Command, handler CommandHandler) error {
	t := reflect.TypeOf(command)
	if t.Elem().Kind() != reflect.Struct {
		return errors.New("event type must be struct")
	}

	if _, ok := r.registeredCommands[t.Name()]; ok {
		return fmt.Errorf("%s has been already added", t.Name())
	} else {
		r.registeredCommands[t.Name()] = handler
	}
	return nil
}

func (r register) GetCommandHandler(command Command) (CommandHandler, bool) {
	t := reflect.TypeOf(command)
	if fns, ok := r.registeredCommands[t.Name()]; ok {
		return fns, true
	}
	return nil, false
}
