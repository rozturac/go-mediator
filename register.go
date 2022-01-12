package mediator

import (
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
	if t.Kind() == reflect.Struct {
		return fmt.Errorf("event must be pointer")
	} else if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("event type must be struct")
	}

	for _, handler := range handlers {
		r.registeredEvents[t.String()] = append(r.registeredEvents[t.String()], handler)
	}
	return nil
}

func (r register) GetEventHandler(event Event) ([]EventHandler, bool) {
	t := reflect.TypeOf(event)
	if fns, ok := r.registeredEvents[t.String()]; ok {
		return fns, true
	}
	return nil, false
}

func (r register) RegisterCommand(command Command, handler CommandHandler) error {
	t := reflect.TypeOf(command)
	if t.Kind() == reflect.Struct {
		return fmt.Errorf("command must be pointer")
	} else if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("command type must be struct")
	}

	if _, ok := r.registeredCommands[t.String()]; ok {
		return fmt.Errorf("%s has been already added", t.Name())
	} else {
		r.registeredCommands[t.String()] = handler
	}
	return nil
}

func (r register) GetCommandHandler(command Command) (CommandHandler, bool) {
	t := reflect.TypeOf(command)
	if fns, ok := r.registeredCommands[t.String()]; ok {
		return fns, true
	}
	return nil, false
}
