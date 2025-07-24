package tf6

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

// ApplyManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ApplyManagedResourceChange(ctx context.Context, req *providerops.ApplyManagedResourceChangeRequest) (providerops.ApplyManagedResourceChangeResponse, error) {
	panic("unimplemented")
}

// ImportManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ImportManagedResourceState(ctx context.Context, req *providerops.ImportManagedResourceStateRequest) (providerops.ImportManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// MoveManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) MoveManagedResourceState(ctx context.Context, req *providerops.MoveManagedResourceStateRequest) (providerops.MoveManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// PlanManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) PlanManagedResourceChange(ctx context.Context, req *providerops.PlanManagedResourceChangeRequest) (providerops.PlanManagedResourceChangeResponse, error) {
	panic("unimplemented")
}

// ReadManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error) {
	panic("unimplemented")
}

// UpgradeManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) UpgradeManagedResourceState(ctx context.Context, req *providerops.UpgradeManagedResourceStateRequest) (providerops.UpgradeManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// ValidateManagedResourceConfig implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ValidateManagedResourceConfig(ctx context.Context, req *providerops.ValidateManagedResourceConfigRequest) (providerops.ValidateManagedResourceConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin6.ValidateResourceConfig_Request{
		TypeName: req.ResourceType,
		Config:   configVal,
	}

	protoResp, err := p.client.ValidateResourceConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateManagedResourceConfigResponse{proto: protoResp}, nil
}

type validateManagedResourceConfigResponse struct {
	proto *tfplugin6.ValidateResourceConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateEphemeralResourceConfigResponse.
func (v validateManagedResourceConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
