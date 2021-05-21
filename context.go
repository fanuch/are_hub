package are_hub

import "context"

// Unique, unexported key type to prevent collisions
type ctxKey uint

// See: https://golang.org/ref/spec#Constant_declarations
// and: https://golang.org/ref/spec#Iota
// and: https://golang.org/doc/effective_go.html#constants
const (
	keyChannel ctxKey = iota
)

// Types implementing this interface can be stored in a context.
type ContextEmbeddable interface {
	ToCtx(context.Context) context.Context
}
