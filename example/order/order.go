package main

import (
	"context"
	"fmt"
)

//go:generate event -type=OrderPlaced,OrderShipped

type OrderPlaced struct {
	OrderID    int
	CustomerID int
}

type OrderShipped struct {
	OrderID int
}

func main() {
	OnOrderPlaced(func(ctx context.Context, e OrderPlaced) {
		fmt.Println("New Order")
		fmt.Println("Order ID:   ", e.OrderID)
		fmt.Println("Customer ID:", e.CustomerID)
		fmt.Println()
	})

	OnOrderShipped(func(ctx context.Context, e OrderShipped) {
		fmt.Println("Order Shipped")
		fmt.Println("Order ID:", e.OrderID)
		fmt.Println()
	})

	ctx := context.Background()

	// From elsewhere in your application:
	EmitOrderPlaced(ctx, OrderPlaced{
		OrderID:    5,
		CustomerID: 265,
	})
	EmitOrderShipped(ctx, OrderShipped{
		OrderID: 5,
	})
}
