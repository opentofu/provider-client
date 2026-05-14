package mockutil_test

import (
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func TestComparer(t *testing.T) {
	tests := map[string]struct {
		Value, Match, Mismatch any
	}{
		"cty.Value": {
			cty.StringVal("hello"),
			cty.StringVal("hello"),
			cty.StringVal("goodbye"),
		},
		"providerschema.DynamicValueIn": {
			providerschema.NewDynamicValue(cty.StringVal("hello"), cty.String),
			providerschema.NewDynamicValue(cty.StringVal("hello"), cty.String),
			providerschema.NewDynamicValue(cty.StringVal("goodbye"), cty.String),
		},
		"proto.Message": {
			&tfplugin6.CallFunction_Request{Name: "hello"},
			&tfplugin6.CallFunction_Request{Name: "hello"},
			&tfplugin6.CallFunction_Request{Name: "goodbye"},
		},
		"tfplugin5.DynamicValue.Json": {
			&tfplugin5.DynamicValue{
				Json: []byte(`{"name":"Eve","age": 45}`),
			},
			&tfplugin5.DynamicValue{
				// This intentionally differs only by irrelevant serialization details.
				Json: []byte(`{"age": 45,"name":"Eve"}`),
			},
			&tfplugin5.DynamicValue{
				Json: []byte(`{"name":"Eve","age": 46}`),
			},
		},
		"tfplugin5.DynamicValue.Msgpack": {
			&tfplugin5.DynamicValue{
				// {"name":"Eve","age": 45}
				Msgpack: []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2d},
			},
			&tfplugin5.DynamicValue{
				// This is semantically equivalent to the above but using a longer encoding for the age value.
				Msgpack: []byte{0xdf, 0x00, 0x00, 0x00, 0x02, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2d},
			},
			&tfplugin5.DynamicValue{
				// Same as the base value except "age" is now 46
				Msgpack: []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2e},
			},
		},
		"tfplugin6.DynamicValue.Json": {
			&tfplugin6.DynamicValue{
				Json: []byte(`{"name":"Eve","age": 45}`),
			},
			&tfplugin6.DynamicValue{
				// This intentionally differs only by irrelevant serialization details.
				Json: []byte(`{"age": 45,"name":"Eve"}`),
			},
			&tfplugin6.DynamicValue{
				Json: []byte(`{"name":"Eve","age": 46}`),
			},
		},
		"tfplugin6.DynamicValue.Msgpack": {
			&tfplugin6.DynamicValue{
				// {"name":"Eve","age": 45}
				Msgpack: []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2d},
			},
			&tfplugin6.DynamicValue{
				// This is semantically equivalent to the above but using a longer encoding for the age value.
				Msgpack: []byte{0xdf, 0x00, 0x00, 0x00, 0x02, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2d},
			},
			&tfplugin6.DynamicValue{
				// Same as the base value except "age" is now 46
				Msgpack: []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa3, 0x45, 0x76, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x2e},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			matcher := mockutil.Eq(test.Value)
			if !matcher.Matches(test.Match) {
				t.Error("did not match when it should have matched")
			}
			if matcher.Matches(test.Mismatch) {
				t.Error("matched when it should not have matched")
			}

			// While we're here we'll also test Diff, since it has similar inputs.
			if diff := mockutil.Diff(test.Value, test.Match); diff != "" {
				t.Error("unexpected diff\n" + diff)
			}
			if diff := mockutil.Diff(test.Value, test.Mismatch); diff == "" {
				t.Error("missing expected diff")
			}
		})
	}
}
