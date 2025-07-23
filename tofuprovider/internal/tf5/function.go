package tf5

import (
	"context"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

// CallFunction implements tofuprovider.GRPCPluginProvider.
func (p *Provider) CallFunction(ctx context.Context, req *providerops.CallFunctionRequest) (providerops.CallFunctionResponse, error) {
	panic("unimplemented")
}
