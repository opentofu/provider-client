package tf5

import (
	"fmt"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
	"github.com/zclconf/go-cty/cty"
)

func makeDynamicValueMsgpack(dv providerschema.DynamicValueIn) (*tfplugin5.DynamicValue, error) {
	if dv == providerschema.NoDynamicValue {
		return nil, fmt.Errorf("missing required value")
	}
	buf, err := common.CtyValueAsMsgpack(dv.Value(), dv.SerializationType())
	if err != nil {
		return nil, fmt.Errorf("cannot serialize to MessagePack: %w", err)
	}
	return &tfplugin5.DynamicValue{
		Msgpack: buf,
	}, nil
}

type dynamicValue struct {
	proto *tfplugin5.DynamicValue
	common.SealedImpl
}

// AsCtyValue implements providerschema.DynamicValueOut.
func (d dynamicValue) AsCtyValue(withType cty.Type) (cty.Value, error) {
	switch {
	case len(d.proto.Msgpack) != 0:
		raw := common.CtyValueMsgpack(d.proto.Msgpack)
		return raw.AsCtyValue(withType)
	case len(d.proto.Json) != 0:
		raw := common.CtyValueJSON(d.proto.Json)
		return raw.AsCtyValue(withType)
	default:
		return cty.NilVal, fmt.Errorf("unsupported value serialization format")
	}
}
