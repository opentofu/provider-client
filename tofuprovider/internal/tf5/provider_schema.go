package tf5

import (
	"context"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

func (p *Provider) GetProviderSchema(ctx context.Context, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error) {
	panic("not yet implemented")
}

func (p *Provider) GetFunctions(ctx context.Context, req *providerops.GetFunctionsRequest) (providerops.GetFunctionsResponse, error) {
	panic("unimplemented")
}
