package mockutil

import (
	"encoding/json"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// Comparer is a helper for using [cmp] in conjunction with [gomock], so that
// call arguments can be matched using custom comparers and transformers as
// modelled by [cmp].
type Comparer struct {
	opts cmp.Options
}

// NewComparer builds a new [Comparer] that compares values using [CmpOptions]
// along with any other options provided as arguments.
//
// Use [Comparer.Eq] on the result to get a [gomock.Matcher] for use
// in describing a mock call.
//
// If you're not passing any extra options then consider just calling the
// package-level [Eq] instead, which uses a preconstructed default comparer.
func NewComparer(opts ...cmp.Option) Comparer {
	fullOpts := buildCmpOptions(opts...)
	return Comparer{
		opts: fullOpts,
	}
}

// Eq returns a [gomock.Matcher] that reports a match if [cmp] considers
// the call's value and the given value to match when tested using the
// comparer's [cmp.Options].
func (c Comparer) Eq(v any) gomock.Matcher {
	return &cmpMatcher{v, c.opts}
}

var defaultComparer = NewComparer()

// Eq calls [Comparer.Eq] on a default [Comparer] that only has the built-in
// set of [cmp.Options], for easier use when no special comparers are needed.
func Eq(v any) gomock.Matcher {
	return defaultComparer.Eq(v)
}

// Diff is a thin wrapper around [cmp.Diff] that automatically includes
// this package's default options, from [CmpOptions].
func Diff(a, b any, opts ...cmp.Option) string {
	fullOpts := buildCmpOptions(opts...)
	return cmp.Diff(a, b, fullOpts)
}

type cmpMatcher struct {
	v    any
	opts cmp.Options
}

// Matches implements [gomock.Matcher].
func (m cmpMatcher) Matches(x any) bool {
	return cmp.Equal(x, m.v, m.opts)
}

// String implements [gomock.Matcher].
func (m cmpMatcher) String() string {
	// Trickery: the helpers that gomock.Eq uses to generate a human-readable
	// description of what it does are not exported, so we'll just make a
	// synthetic gomock.Eq comparer and ask it for its string representation.
	eq := gomock.Eq(m.v)
	return eq.String()
}

// CmpOptions is the base set of [cmp] options that is included by default when
// you use [Eq], [Diff], [NewComparer], etc.
//
// This is exported just in case a caller needs to do something unusual that
// requires a direct call to one of the [cmp] functions. If possible, prefer
// to use one of the wrappers elsewhere in this package to make your test
// more concise.
var CmpOptions = cmp.Options{
	ctydebug.CmpOptions,
	cmp.AllowUnexported(providerschema.DynamicValueIn{}),
	cmp.Transformer("decodeDynamicValue5", func(v *tfplugin5.DynamicValue) any {
		// tfplugin5.DynamicValue is compared by the data structure
		// that's encoded rather than by the specific bytes that encode it,
		// since there are multiple valid ways to serialize the same
		// information.
		var out any
		if len(v.Json) != 0 {
			err := json.Unmarshal(v.Json, &out)
			if err != nil {
				return invalidJSON(v.Json)
			}
		} else {
			err := msgpack.Unmarshal(v.Msgpack, &out)
			if err != nil {
				return invalidMessagePack(v.Msgpack)
			}
		}
		return out
	}),
	cmp.Transformer("decodeDynamicValue6", func(v *tfplugin6.DynamicValue) any {
		// tfplugin6.DynamicValue is compared by the data structure
		// that's encoded rather than by the specific bytes that encode it,
		// since there are multiple valid ways to serialize the same
		// information.
		var out any
		if len(v.Json) != 0 {
			err := json.Unmarshal(v.Json, &out)
			if err != nil {
				return invalidJSON(v.Json)
			}
		} else {
			err := msgpack.Unmarshal(v.Msgpack, &out)
			if err != nil {
				return invalidMessagePack(v.Msgpack)
			}
		}
		return out
	}),
	// We need to pre-filter protocmp.Transform because otherwise it'd be
	// ambiguous with our DynamicValue-specific transformers above.
	cmp.FilterValues(
		func(a, b proto.Message) bool {
			return !(isDynamicValueMessage(a) || isDynamicValueMessage(b))
		},
		protocmp.Transform(),
	),
	cmp.Transformer("comparableDiagnostics", cmpTransformDiagnostics),
	interfaceTypeTransformer[providerops.Diagnostic](
		"comparableDiagnostic", func(v any) *ComparableDiagnostic {
			return cmpTransformDiagnostic(v.(providerops.Diagnostic))
		},
	),
	interfaceTypeTransformer[providerops.FunctionError](
		"comparableFunctionError", func(v any) *ComparableFunctionError {
			return cmpTransformFunctionError(v.(providerops.FunctionError))
		},
	),
}

// interfaceTypeTransformer is a helper for using [cmp.Transformer] with
// interface types.
//
// Unfortunately this requires some additional machinery because of how
// [reflect.ValueOf] calls implicitly convert interface-typed values to the
// generic [any], losing any more specific interface type that the value
// had. To work around that we use [cmp.FilterPath] to test whether a value
// implements the interface, and then use a transform function whose argument
// is [any] to make sure it can match any concrete type.
//
// The given transform function should immediately unconditionally type-assert
// the given argument to type Interface, which is guaranteed to succeed due
// to this helper's filtering logic.
func interfaceTypeTransformer[Interface, Concrete any](name string, tr func(any) Concrete) cmp.Option {
	ifaceType := reflect.TypeFor[Interface]()
	concreteType := reflect.TypeFor[Concrete]()
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			step := p.Last()
			gotType := step.Type()
			if gotType == concreteType {
				// Avoid recursively transforming the transformer's own result.
				return false
			}
			if gotType.Implements(ifaceType) {
				// Easy case: it's a concrete type that implements the interface.
				return true
			}
			if gotType.Kind() == reflect.Interface {
				// Unfortunately if the given value was of an interface type
				// already then gotType will just be reflect.TypeFor[any](), in
				// which case we need to check the values themselves.
				v1, v2 := step.Values()
				if !(v1.IsValid() && v2.IsValid()) {
					return false
				}
				if v1.IsNil() || v2.IsNil() {
					return false
				}
				// Since v1 and v2 are both representing interface values,
				// "Elem" here means to take the concrete-typed value inside.
				return v1.Elem().Type().Implements(ifaceType) && v2.Elem().Type().Implements(ifaceType)
			}
			return false
		},
		cmp.Transformer(name, tr),
	)
}

func buildCmpOptions(opts ...cmp.Option) cmp.Options {
	if len(opts) == 0 {
		return CmpOptions // no need to allocate in this common case
	}
	ret := make(cmp.Options, 0, 1+len(opts))
	ret = append(ret, CmpOptions)
	ret = append(ret, opts...)
	return ret
}

func isDynamicValueMessage(msg proto.Message) bool {
	switch msg.(type) {
	case *tfplugin5.DynamicValue, *tfplugin6.DynamicValue:
		return true
	default:
		return false
	}
}

type invalidMessagePack []byte
type invalidJSON []byte
