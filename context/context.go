package context

import (
	"context"

	"github.com/google/uuid"
)

// fleetStateString is a type definition over string to use as a key in context
type fleetStateString string

// package constant
const (
	processID fleetStateString = "processID"
)

// SetProcessID generates new uuid and sets it as a processID into the context
// if en error occures during creation of uuid, parent context and an error will be returned
func SetProcessID(ctx context.Context) (context.Context, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, processID, id.String()), nil
}

// ProcessID gets a as a processID from the context
func ProcessID(ctx context.Context) string {
	value := ctx.Value(processID)
	if value == nil {
		return ""
	}
	return value.(string)
}
