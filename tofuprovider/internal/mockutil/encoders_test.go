package mockutil_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
)

func TestJSON(t *testing.T) {
	// The primary purpose of JSON is to be used when constructing expected
	// DynamicValue objects for comparison using [mockutil.Eq], so we'll
	// test it like that.
	got := &tfplugin6.DynamicValue{
		Json: mockutil.JSON(map[string]any{
			"greeting": "Hello",
		}),
	}
	want := &tfplugin6.DynamicValue{
		Json: []byte(`{ "greeting": "Hello" }`),
	}
	if diff := cmp.Diff(want, got, mockutil.CmpOptions); diff != "" {
		t.Error("wrong result\n" + diff)
	}
}

func TestMsgPack(t *testing.T) {
	// The primary purpose of JSON is to be used when constructing expected
	// DynamicValue objects for comparison using [mockutil.Eq], so we'll
	// test it like that.
	got := &tfplugin6.DynamicValue{
		Msgpack: mockutil.MsgPack(map[string]any{
			"greeting": "Hello",
		}),
	}
	want := &tfplugin6.DynamicValue{
		Msgpack: []byte{0xde, 0x00, 0x01, 0xa8, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0xa5, 0x48, 0x65, 0x6c, 0x6c, 0x6f},
	}
	if diff := cmp.Diff(want, got, mockutil.CmpOptions); diff != "" {
		t.Error("wrong result\n" + diff)
	}
}
