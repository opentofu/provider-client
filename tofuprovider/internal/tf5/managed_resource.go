package tf5

import (
	"context"

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

// ReadManagedResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error) {
	panic("unimplemented")
}

// UpgradeManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) UpgradeManagedResourceState(ctx context.Context, req *providerops.UpgradeManagedResourceStateRequest) (providerops.UpgradeManagedResourceStateResponse, error) {
	panic("unimplemented")
}
