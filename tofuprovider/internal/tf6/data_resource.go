package tf6

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

// ReadDataResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadDataResource(ctx context.Context, req *providerops.ReadDataResourceRequest) (providerops.ReadDataResourceResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	providerMetaVal, err := makeDynamicValueMsgpack(req.ProviderMeta)
	if err != nil {
		return nil, fmt.Errorf("invalid ProviderMeta value: %w", err)
	}
	protoReq := &tfplugin6.ReadDataSource_Request{
		TypeName:           req.ResourceType,
		Config:             configVal,
		ProviderMeta:       providerMetaVal,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.ReadDataSource(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return readDataResourceResponse{proto: protoResp}, nil

}

type readDataResourceResponse struct {
	proto *tfplugin6.ReadDataSource_Response
	common.SealedImpl
}

// Deferred implements providerops.ReadDataResourceResponse.
func (r readDataResourceResponse) Deferred() providerops.Deferred {
	if r.proto.Deferred == nil {
		return nil
	}
	return deferred{proto: r.proto.Deferred}
}

// State implements providerops.ReadDataResourceResponse.
func (r readDataResourceResponse) State() providerschema.DynamicValueOut {
	if r.proto.State == nil {
		return nil
	}
	return dynamicValue{proto: r.proto.State}
}

// Diagnostics implements providerops.ReadDataResourceResponse.
func (r readDataResourceResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: r.proto.Diagnostics}
}

// ValidateDataResourceConfig implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ValidateDataResourceConfig(ctx context.Context, req *providerops.ValidateDataResourceConfigRequest) (providerops.ValidateDataResourceConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin6.ValidateDataResourceConfig_Request{
		TypeName: req.ResourceType,
		Config:   configVal,
	}

	protoResp, err := p.client.ValidateDataResourceConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateDataResourceConfigResponse{proto: protoResp}, nil
}

type validateDataResourceConfigResponse struct {
	proto *tfplugin6.ValidateDataResourceConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateDataResourceConfigResponse.
func (v validateDataResourceConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
