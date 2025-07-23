package tf6

import (
	"context"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

// ReadDataResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadDataResource(ctx context.Context, req *providerops.ReadDataResourceRequest) (providerops.ReadDataResourceResponse, error) {
	panic("unimplemented")
}

// ValidateDataResourceConfig implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ValidateDataResourceConfig(ctx context.Context, req *providerops.ValidateDataResourceConfigRequest) (providerops.ValidateDataResourceConfigResponse, error) {
	panic("unimplemented")
}
