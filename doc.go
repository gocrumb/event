// Event is a tool to generate observer pattern code for your events. Given
// the name(s) of event types, this tool will generate the following
// symbols in a self-contained Go source file:
//
//	 func On{T}(f func(context.Context, {T})) { ... }
//	 func Emit{T}(ctx context.Context, e {T}) { ... }
//
// For example, given this snippet,
//
// 	package ev
//
// 	type OrderPlaced struct {
// 		OrderID    int
// 		CustomerID int
// 	}
//
// running this command
//
// 	event -type=OrderPlaced
//
// in the same directory will create the file event.go, in package ev, containing the definitions of all
// the relevant symbols.
//
// You can then listen for events by doing this:
//
// 	OnOrderPlaced(func(ctx context.Context, e OrderPlaced) {
// 		...
// 	})
//
// And, emit events by doing this:
//
// 	EmitOrderPlaced(ctx, OrderPlaced{ ... })
//
// Typically this process would be run using go generate, like
// this:
//
// 	//go:generate event -type=OrderPlaced
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate symbols for multiple types. The default output file is
// event.go. It can be overridden with the -output flag.
//
// A more complete list of symbols generated for each type is as follows:
//
//	 const Event{T} EventType = iota
//
//	 func (e {T}) Type() EventType { ... }
//
//	 func (e {T}) Trigger(ctx context.Context) { ... }
//
//	 type {T}Handler interface{
//	 	Handle(context.Context, {T})
//	 }
//
//	 type {T}HandlerFunc func(context.Context, {T})
//
//	 func (f {T}HandlerFunc) Handle(ctx context.Context, e {T}) { ... }
//
//	 type {T}Emitter struct { ... }
//
//	 func (m *{T}Emitter) Trigger(ctx context.Context, e {T}) { ... }
//
//	 func (m *{T}Emitter) Handle(h {T}Handler) { ... }
//
//	 func (m *{T}Emitter) HandleFunc(f func(context.Context, {T})) { ... }
//
//	 func On{T}(f func(context.Context, {T})) { ... }
//
//	 func Emit{T}(ctx context.Context, e {T}) { ... }
//
//	 var Emitter{T}  = {T}Emitter{}
//
package main
