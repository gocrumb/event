// Package event provides a small, generic, in-process emitter for the
// observer pattern.
//
// Declare an event type and a typed emitter for it:
//
//	type OrderPlaced struct {
//		OrderID    int
//		CustomerID int
//	}
//
//	var OrderPlacedEvent event.Emitter[OrderPlaced]
//
// Register handlers:
//
//	OrderPlacedEvent.HandleFunc(func(ctx context.Context, e OrderPlaced) {
//		// ...
//	})
//
// Trigger the event:
//
//	OrderPlacedEvent.Trigger(ctx, OrderPlaced{OrderID: 5, CustomerID: 265})
//
// To log each triggered event at debug level, set a default logger for all
// emitters:
//
//	event.SetDefaultLogger(logger)
//
// or create an emitter with its own logger:
//
//	var OrderPlacedEvent = event.New[OrderPlaced](event.Logger(logger))
//
// The zero value of Emitter is ready to use. Handlers are invoked
// synchronously, in the order they were registered, on the goroutine that
// called Trigger.
package event
