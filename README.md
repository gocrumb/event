# Event

Manage and dispatch events within your Go application. This is somewhat inspired from Laravel's simple approach to the observer pattern.

## Installation

Install grumbs/event using the go get command:

```
$ go get github.com/grumbs/event
```

The package requires no additional dependencies other than Go itself.

## Example

Run `go generate` with a file like the following:

``` golang
package main

//go:generate event -type=OrderPlaced,OrderShipped

type OrderPlaced struct {
	OrderID    int
	CustomerID int
}

type OrderShipped struct {
	OrderID int
}
```

Use events from everywhere within your applications:

``` golang
OnOrderPlaced(func(e OrderPlaced) {
	fmt.Println("New Order")
	fmt.Println("Order ID:   ", e.OrderID)
	fmt.Println("Customer ID:", e.CustomerID)
	fmt.Println()
})

OnOrderShipped(func(e OrderShipped) {
	fmt.Println("Order Shipped")
	fmt.Println("Order ID:", e.OrderID)
	fmt.Println()
})

// From elsewhere in your application:
EmitOrderPlaced(OrderPlaced{
	OrderID:    5,
	CustomerID: 265,
})
EmitOrderShipped(OrderShipped{
	OrderID: 5,
})
```