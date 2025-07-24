package tf5

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin5"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

// CloseEphemeralResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) CloseEphemeralResource(ctx context.Context, req *providerops.CloseEphemeralResourceRequest) (providerops.CloseEphemeralResourceResponse, error) {
	panic("unimplemented")
}

// OpenEphemeralResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) OpenEphemeralResource(ctx context.Context, req *providerops.OpenEphemeralResourceRequest) (providerops.OpenEphemeralResourceResponse, error) {
	panic("unimplemented")
}

// RenewEphemeralResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) RenewEphemeralResource(ctx context.Context, req *providerops.RenewEphemeralResourceRequest) (providerops.RenewEphemeralResourceResponse, error) {
	panic("unimplemented")
}

// ValidateEphemeralResourceConfig implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ValidateEphemeralResourceConfig(ctx context.Context, req *providerops.ValidateEphemeralResourceConfigRequest) (providerops.ValidateEphemeralResourceConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin5.ValidateEphemeralResourceConfig_Request{
		TypeName: req.ResourceType,
		Config:   configVal,
	}

	protoResp, err := p.client.ValidateEphemeralResourceConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateEphemeralResourceConfigResponse{proto: protoResp}, nil
}

type validateEphemeralResourceConfigResponse struct {
	proto *tfplugin5.ValidateEphemeralResourceConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateEphemeralResourceConfigResponse.
func (v validateEphemeralResourceConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
