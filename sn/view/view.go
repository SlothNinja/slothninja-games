package view

import "golang.org/x/net/context"

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// dataKey is the key for Data values in Contexts. It is
// unexported; clients use view.NewContext and view.FromContext
// instead of using this key directly.
var dataKey key = 0

// Data provides storage for a view data
type Data map[string]interface{}

// NewContext returns a new Context that carries Data
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, dataKey, make(Data, 0))
}

// FromContext returns the Data value stored in ctx, if any.
func FromContext(ctx context.Context) Data {
	if d, ok := ctx.Value(dataKey).(Data); !ok {
		return nil
	} else {
		return d
	}
}
