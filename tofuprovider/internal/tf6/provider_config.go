package tf6

import (
	"context"
	"fmt"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
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

func (p *Provider) ConfigureProvider(ctx context.Context, req *providerops.ConfigureProviderRequest) (providerops.ConfigureProviderResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin6.ConfigureProvider_Request{
		Config:             configVal,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.ConfigureProvider(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return configureProviderResponse{proto: protoResp}, nil
}

type configureProviderResponse struct {
	proto *tfplugin6.ConfigureProvider_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateProviderConfigResponse.
func (v configureProviderResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
