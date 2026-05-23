# Event

[![Go Reference](https://pkg.go.dev/badge/github.com/gocrumb/event.svg)](https://pkg.go.dev/github.com/gocrumb/event)

A small, generic, in-process event emitter for Go. Inspired by Laravel's simple take on the observer pattern.

## Installation

```
$ go get github.com/gocrumb/event
```

Requires Go 1.21+.

## Usage

Declare an event type and a typed emitter for it:

``` golang
package main

import (
	"context"
	"fmt"

	"github.com/gocrumb/event"
)

type OrderPlaced struct {
	OrderID    int
	CustomerID int
}

var OrderPlacedEvent event.Emitter[OrderPlaced]
```

Register handlers:

``` golang
OrderPlacedEvent.HandleFunc(func(ctx context.Context, e OrderPlaced) {
	fmt.Println("New Order")
	fmt.Println("Order ID:   ", e.OrderID)
	fmt.Println("Customer ID:", e.CustomerID)
})
```

Trigger the event from elsewhere in your application:

``` golang
OrderPlacedEvent.Trigger(ctx, OrderPlaced{
	OrderID:    5,
	CustomerID: 265,
})
```

The zero value of `Emitter` is ready to use. Handlers run synchronously, in registration order, on the goroutine that calls `Trigger`.

## Documentation

- [Reference](https://pkg.go.dev/github.com/gocrumb/event)

## Contributing

Contributions are welcome.

## License

This package is available under the [BSD (3-Clause) License](https://opensource.org/licenses/BSD-3-Clause).
