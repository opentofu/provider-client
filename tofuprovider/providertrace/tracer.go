package providertrace

import (
	"context"
	"io"
)

// A Tracer provides a set of hooks that callers of tofuprovider.Start can
// use to be notified when certain interesting events happen during the life
// of the provider client.
//
// Pass a Tracer to tofuprovider.Start by first wrapping it in a
// [context.Context] using [ContextWithTracer] and then passing that context
// (or a child of it with the same values) to tofuprovider.Start.
type Tracer struct {
	// If non-nil, ChildStderr is used as the stderr stream for the plugin
	// child process. If nil then any data the child process writes to stderr
	// is immediately discarded.
	//
	// Plugins are not required to produce data on stderr in any particular
	// format, or even to generate any data at all, but plugins implemented
	// using the HashiCorp Plugin Framework or SDK tend to produce a stream of
	// JSON objects separated by newline characters where each object
	// represents a log message that might be useful for debugging the
	// provider's behavior, and so callers may wish to expose that information
	// somehow.
	ChildStderr io.Writer
}

var defaultTracer = &Tracer{}

// ContextWithTracer returns a new context, child of parent, which carries
// the given [Tracer] object for use in a call to tofuprovider.Start.
func ContextWithTracer(parent context.Context, tracer *Tracer) context.Context {
	if tracer == nil {
		return parent
	}
	return context.WithValue(parent, tracerKey(0), tracer)
}

// TracerFromContext returns the tracer that was previously associated with
// the given context (or one of its parents) using [ContextWithTracer], or
// returns a default no-op tracer if none was explicitly provided.
//
// Callers must not modify the returned tracer because it may be shared with
// other callers running concurrently.
func TracerFromContext(ctx context.Context) *Tracer {
	tracer, ok := ctx.Value(tracerKey(0)).(*Tracer)
	if !ok {
		return defaultTracer
	}
	return tracer
}

type tracerKey int
