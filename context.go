package gateway

import (
	"context"
	"github.com/getas/alb-gateway/events"
)

// key is the type used for any items added to the request context.
type key int

// requestContextKey is the key for the alb lambda target group proxy `RequestContext`.
const requestContextKey key = iota

// newContext returns a new Context with specific alb lambda target group proxy values.
func newContext(ctx context.Context, e events.LambdaTargetGroupRequest) context.Context {
	return context.WithValue(ctx, requestContextKey, e.RequestContext)
}

// RequestContext returns the ALBLambdaRequestContext value stored in ctx.
func RequestContext(ctx context.Context) (events.LambdaTargetGroupRequest, bool) {
	c, ok := ctx.Value(requestContextKey).(events.LambdaTargetGroupRequest)
	return c, ok
}
