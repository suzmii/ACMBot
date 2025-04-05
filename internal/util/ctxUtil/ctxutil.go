package ctxUtil

import (
	"context"
)

type key[T any] struct{}

func Set[T any](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, key[T]{}, val)
}

func Get[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}

func MustGet[T any](ctx context.Context) T {
	return ctx.Value(key[T]{}).(T)
}
