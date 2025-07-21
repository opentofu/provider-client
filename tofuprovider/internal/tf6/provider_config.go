package tf6

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/internal/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

func (p *Provider) ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin6.ValidateProviderConfig_Request{
		Config: configVal,
	}

	protoResp, err := p.client.ValidateProviderConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateProviderConfigResponse{proto: protoResp}, nil
}

type validateProviderConfigResponse struct {
	proto *tfplugin6.ValidateProviderConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateProviderConfigResponse.
func (v validateProviderConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
