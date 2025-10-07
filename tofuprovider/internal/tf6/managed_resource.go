package tf6

import (
	"context"
	"fmt"
	"iter"
	"slices"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// ApplyManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ApplyManagedResourceChange(ctx context.Context, req *providerops.ApplyManagedResourceChangeRequest) (providerops.ApplyManagedResourceChangeResponse, error) {
	priorState, err := makeDynamicValueMsgpack(req.PriorState)
	if err != nil {
		return nil, fmt.Errorf("invalid PriorState value: %w", err)
	}
	plannedNewState, err := makeDynamicValueMsgpack(req.PlannedNewState)
	if err != nil {
		return nil, fmt.Errorf("invalid PlannedNewState value: %w", err)
	}
	config, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}

	var providerMeta *tfplugin6.DynamicValue
	if req.ProviderMeta != providerschema.NoDynamicValue {
		providerMeta, err = makeDynamicValueMsgpack(req.ProviderMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid ProviderMeta value: %w", err)
		}
	}

	protoReq := &tfplugin6.ApplyResourceChange_Request{
		TypeName:       req.ResourceType,
		PriorState:     priorState,
		PlannedState:   plannedNewState,
		Config:         config,
		PlannedPrivate: req.PlannedProviderInternal,
		ProviderMeta:   providerMeta,
	}

	protoResp, err := p.client.ApplyResourceChange(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return applyManagedResourceChangeResponse{proto: protoResp}, nil
}

// ImportManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ImportManagedResourceState(ctx context.Context, req *providerops.ImportManagedResourceStateRequest) (providerops.ImportManagedResourceStateResponse, error) {
	protoReq := &tfplugin6.ImportResourceState_Request{
		TypeName:           req.ResourceType,
		Id:                 req.ID,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.ImportResourceState(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return importManagedResourceStateResponse{proto: protoResp}, nil
}

// MoveManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) MoveManagedResourceState(ctx context.Context, req *providerops.MoveManagedResourceStateRequest) (providerops.MoveManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// PlanManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) PlanManagedResourceChange(ctx context.Context, req *providerops.PlanManagedResourceChangeRequest) (providerops.PlanManagedResourceChangeResponse, error) {
	priorState, err := makeDynamicValueMsgpack(req.PriorState)
	if err != nil {
		return nil, fmt.Errorf("invalid PriorState value: %w", err)
	}
	proposedNewState, err := makeDynamicValueMsgpack(req.ProposedNewState)
	if err != nil {
		return nil, fmt.Errorf("invalid ProposedNewState value: %w", err)
	}
	config, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}

	var providerMeta *tfplugin6.DynamicValue
	if req.ProviderMeta != providerschema.NoDynamicValue {
		providerMeta, err = makeDynamicValueMsgpack(req.ProviderMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid ProviderMeta value: %w", err)
		}
	}

	protoReq := &tfplugin6.PlanResourceChange_Request{
		TypeName:           req.ResourceType,
		PriorState:         priorState,
		ProposedNewState:   proposedNewState,
		Config:             config,
		PriorPrivate:       req.PriorProviderInternal,
		ProviderMeta:       providerMeta,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.PlanResourceChange(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return planManagedResourceChangeResponse{proto: protoResp}, nil
}

// ReadManagedResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error) {
	currentState, err := makeDynamicValueMsgpack(req.CurrentState)
	if err != nil {
		return nil, fmt.Errorf("invalid CurrentState value: %w", err)
	}

	var providerMeta *tfplugin6.DynamicValue
	if req.ProviderMeta != providerschema.NoDynamicValue {
		providerMeta, err = makeDynamicValueMsgpack(req.ProviderMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid ProviderMeta value: %w", err)
		}
	}

	protoReq := &tfplugin6.ReadResource_Request{
		TypeName:           req.ResourceType,
		CurrentState:       currentState,
		Private:            req.ProviderInternal,
		ProviderMeta:       providerMeta,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.ReadResource(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return readManagedResourceResponse{proto: protoResp}, nil
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

type planManagedResourceChangeResponse struct {
	proto *tfplugin6.PlanResourceChange_Response
	common.SealedImpl
}

// Diagnostics implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: p.proto.Diagnostics}
}

// PlannedNewState implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) PlannedNewState() providerschema.DynamicValueOut {
	return dynamicValue{proto: p.proto.PlannedState}
}

// PlannedProviderInternal implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) PlannedProviderInternal() []byte {
	return p.proto.PlannedPrivate
}

// LegacyTypeSystem implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) LegacyTypeSystem() bool {
	return p.proto.LegacyTypeSystem
}

// Deferred implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) Deferred() providerops.Deferred {
	if p.proto.Deferred == nil {
		return nil
	}
	return deferred{proto: p.proto.Deferred}
}

type applyManagedResourceChangeResponse struct {
	proto *tfplugin6.ApplyResourceChange_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ApplyManagedResourceChangeResponse.
func (a applyManagedResourceChangeResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: a.proto.Diagnostics}
}

// PlannedNewState implements providerops.ApplyManagedResourceChangeResponse.
func (a applyManagedResourceChangeResponse) PlannedNewState() providerschema.DynamicValueOut {
	if a.proto.NewState == nil {
		return nil
	}
	return dynamicValue{proto: a.proto.NewState}
}

// ProviderInternal implements providerops.ApplyManagedResourceChangeResponse.
func (a applyManagedResourceChangeResponse) ProviderInternal() []byte {
	return a.proto.Private
}

// LegacyTypeSystem implements providerops.ApplyManagedResourceChangeResponse.
func (a applyManagedResourceChangeResponse) LegacyTypeSystem() bool {
	return a.proto.LegacyTypeSystem
}

type readManagedResourceResponse struct {
	proto *tfplugin6.ReadResource_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ReadManagedResourceResponse.
func (r readManagedResourceResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: r.proto.Diagnostics}
}

// NewState implements providerops.ReadManagedResourceResponse.
func (r readManagedResourceResponse) NewState() providerschema.DynamicValueOut {
	if r.proto.NewState == nil {
		return nil
	}
	return dynamicValue{proto: r.proto.NewState}
}

// ProviderInternal implements providerops.ReadManagedResourceResponse.
func (r readManagedResourceResponse) ProviderInternal() []byte {
	return r.proto.Private
}

// Deferred implements providerops.ReadManagedResourceResponse.
func (r readManagedResourceResponse) Deferred() providerops.Deferred {
	if r.proto.Deferred == nil {
		return nil
	}
	return deferred{proto: r.proto.Deferred}
}

type importManagedResourceStateResponse struct {
	proto *tfplugin6.ImportResourceState_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ImportManagedResourceStateResponse.
func (i importManagedResourceStateResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: i.proto.Diagnostics}
}

// ImportedResources implements providerops.ImportManagedResourceStateResponse.
func (i importManagedResourceStateResponse) ImportedResources() iter.Seq[providerops.ImportedManagedResource] {
	return common.MapSeq(slices.Values(i.proto.ImportedResources), func(protoRes *tfplugin6.ImportResourceState_ImportedResource) providerops.ImportedManagedResource {
		return importedManagedResource{proto: protoRes}
	})
}

type importedManagedResource struct {
	proto *tfplugin6.ImportResourceState_ImportedResource
	common.SealedImpl
}

// ResourceType implements providerops.ImportedManagedResource.
func (i importedManagedResource) ResourceType() string {
	return i.proto.TypeName
}

// State implements providerops.ImportedManagedResource.
func (i importedManagedResource) State() providerschema.DynamicValueOut {
	if i.proto.State == nil {
		return nil
	}
	return dynamicValue{proto: i.proto.State}
}

// ProviderInternal implements providerops.ImportedManagedResource.
func (i importedManagedResource) ProviderInternal() []byte {
	return i.proto.Private
}
