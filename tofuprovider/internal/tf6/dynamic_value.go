package tf6

import (
	"fmt"

	"github.com/apparentlymart/opentofu-providers/internal/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

func makeDynamicValueMsgpack(dv providerschema.DynamicValueIn) (*tfplugin6.DynamicValue, error) {
	if dv == providerschema.NoDynamicValue {
		return nil, fmt.Errorf("missing required value")
	}
	buf, err := common.CtyValueAsMsgpack(dv.Value(), dv.SerializationType())
	if err != nil {
		return nil, fmt.Errorf("cannot serialize to MessagePack: %w", err)
	}
	return &tfplugin6.DynamicValue{
		Msgpack: buf,
	}, nil
}
