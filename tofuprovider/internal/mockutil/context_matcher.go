package mockutil

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
)

// AnyContext returns a [gomock.Matcher] that matches any [context.Context].
func AnyContext() gomock.Matcher {
	return anyContextMatcher
}

// ContextWithValue returns a [gomock.Matcher] that matches a [context.Context]
// which has the given key/value pair stored in it.
//
// This is intended for tests that verify that context values passed by an
// external caller are correctly propagated to an internal call, such as for
// OpenTelemetry trace span propagation.
func ContextWithValue(key, value any) gomock.Matcher {
	return &contextMatcher{
		check: func(ctx context.Context) bool {
			gotValue := ctx.Value(key)
			return cmp.Equal(gotValue, value, CmpOptions)
		},
		desc: fmt.Sprintf("context.WithValue(..., %#v, %#v)", key, value),
	}
}

var anyContextMatcher = &contextMatcher{
	check: func(ctx context.Context) bool {
		return true
	},
	desc: "any context.Context",
}

type contextMatcher struct {
	check func(context.Context) bool
	desc  string
}

// Matches implements [gomock.Matcher].
func (c *contextMatcher) Matches(x any) bool {
	ctx, ok := x.(context.Context)
	if !ok {
		return false
	}
	return c.check(ctx)
}

// String implements [gomock.Matcher].
func (c *contextMatcher) String() string {
	return c.desc
}
