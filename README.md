# 🚀 Go-Mediator

[![GitHub license](https://img.shields.io/github/license/rozturac/go-mediator.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://github.com/rozturac/go-mediator/blob/main/LICENSE)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/rozturac/go-mediator)
[![Release](https://img.shields.io/github/tag/rozturac/go-mediator.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/rozturac/go-mediator/releases/latest)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge)](https://goreportcard.com/report/github.com/rozturac/go-mediator)

## Installation

Via go packages:
```go get github.com/rozturac/go-mediator```

## Usage

### Register Command

```go
package main

import (
	"context"
	"github.com/rozturac/go-mediator"
	"log"
)

type TestCommand struct {
	name string
}

func TestCommandHandler(ctx context.Context, command mediator.Command) (interface{}, error) {
	return "Success", nil
}

func main() {
	mediator := mediator.Create()
	if err := mediator.RegisterCommand(&TestCommand{}, TestCommandHandler); err != nil {
		log.Fatal(err)
	}
}

```

### Register Behavior
```go
func LoggingBehavior(ctx context.Context, command mediator.Command, next mediator.CommandHandler) (interface{}, error) {
    log.Println(fmt.Sprintf("Request INFO - Command: %v, CorrelationId: %v", command, ctx.Value("CorrelationId")))
    result, err := next(ctx, command)
    log.Println(fmt.Sprintf("Response INFO - Result: %v, Error: %v", result, err))
    return result, err
}

mediator.WithBehavior(LoggingBehavior)
```

### Send Command
```go
command := &TestCommand{name: "Test Command"}
if result, err := mediator.Send(context.Background(), command); err != nil {
	log.Fatal(err)
} else {
	fmt.Println(fmt.Sprintf("Response: %v", result))
}
```

### Send Command As Async
```go
command := &TestCommand{name: "Test Command"}
async := mediator.SendAsync(context.Background(), command)
//We can do something here in the same time..
if result, err := async.Await(); err != nil {
	log.Fatal(err)
} else {
	fmt.Println(fmt.Sprintf("Response: %v", result))
}
```

## License
MIT License

Copyright (c) 2021 Rıdvan ÖZTURAÇ

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
