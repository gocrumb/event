package main

import "fmt"

//go:generate event -type=OrderPlaced,OrderShipped

type OrderPlaced struct {
	OrderID    int
	CustomerID int
}

type OrderShipped struct {
	OrderID int
}

func main() {
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
}
