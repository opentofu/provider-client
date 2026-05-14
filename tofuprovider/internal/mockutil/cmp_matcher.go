package mockutil

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
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
