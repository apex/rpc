package rpc

import (
	"context"
	"net/http"
)

// ctxKey is a private context key.
type ctxKey struct{}

// NewRequestContext returns a new context with ctx.
func NewRequestContext(ctx context.Context, v *http.Request) context.Context {
	return context.WithValue(ctx, ctxKey{}, v)
}

// RequestFromContext returns ctx from context.
func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	v, ok := ctx.Value(ctxKey{}).(*http.Request)
	return v, ok
}
