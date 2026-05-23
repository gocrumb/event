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

type OrderShipped struct {
	OrderID int
}

var (
	OrderPlacedEvent  event.Emitter[OrderPlaced]
	OrderShippedEvent event.Emitter[OrderShipped]
)

func main() {
	OrderPlacedEvent.HandleFunc(func(ctx context.Context, e OrderPlaced) {
		fmt.Println("New Order")
		fmt.Println("Order ID:   ", e.OrderID)
		fmt.Println("Customer ID:", e.CustomerID)
		fmt.Println()
	})

	OrderShippedEvent.HandleFunc(func(ctx context.Context, e OrderShipped) {
		fmt.Println("Order Shipped")
		fmt.Println("Order ID:", e.OrderID)
		fmt.Println()
	})

	ctx := context.Background()

	// From elsewhere in your application:
	OrderPlacedEvent.Trigger(ctx, OrderPlaced{
		OrderID:    5,
		CustomerID: 265,
	})
	OrderShippedEvent.Trigger(ctx, OrderShipped{
		OrderID: 5,
	})
}
