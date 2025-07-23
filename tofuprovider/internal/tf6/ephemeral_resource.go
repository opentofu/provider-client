package tf6

import (
	"context"

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
	panic("unimplemented")
}
