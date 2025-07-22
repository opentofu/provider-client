package tf5

import (
	"context"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

func (p *Provider) ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error) {
	panic("not yet implemented")
}

func (p *Provider) ConfigureProvider(ctx context.Context, req *providerops.ConfigureProviderRequest) (providerops.ConfigureProviderResponse, error) {
	panic("unimplemented")
}
