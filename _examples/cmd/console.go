package main

import (
	"context"
	"fmt"
	"github.com/rozturac/go-mediator"
	"log"
)

type TestCommand struct {
	name string
}

func TestCommandHandler(ctx context.Context, command mediator.Command) (interface{}, error) {
	fmt.Println("Called")
	return "Success", nil
}

func LoggingBehavior(ctx context.Context, command mediator.Command, next mediator.CommandHandler) (interface{}, error) {
	log.Println(fmt.Sprintf("Request INFO - Command: %v, CorrelationId: %v", command, ctx.Value("CorrelationId")))
	result, err := next(ctx, command)
	log.Println(fmt.Sprintf("Response INFO - Result: %v, Error: %v", result, err))
	return result, err
}

func main() {
	mediator := mediator.Create()
	mediator.WithBehavior(LoggingBehavior)
	if err := mediator.RegisterCommand(&TestCommand{}, TestCommandHandler); err != nil {
		log.Fatal(err)
	}

	command := &TestCommand{name: "Test Command"}
	async := mediator.SendAsync(context.Background(), command)
	//We can do something in the same time..
	if result, err := async.Await(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(fmt.Sprintf("Response: %v", result))
	}
}
