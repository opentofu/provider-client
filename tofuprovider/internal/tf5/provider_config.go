package tf5

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin5"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

func (p *Provider) ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin5.PrepareProviderConfig_Request{
		Config: configVal,
	}

	// NOTE: Protocol version 5 uses a slightly different model where in
	// principle it is possible for a provider to return a modified version
	// of the configuration to pass to ConfigureProvider, mimicking how
	// PlanResourceChange can return adjusted values (e.g. inserting default
	// values) to send to ApplyResourceChange.
	//
	// However, in practice providers did not make use of this because the
	// provider configuration values never get saved anywhere anyway and
	// so the provider can just do the same default value insertion etc
	// inside the ConfigureProvider implementation.
	//
	// Therefore modern OpenTofu just ignores the "prepared config" and
	// sends the original config to ConfigureProvider as long as validation
	// succeeds. This library follows that shape, effectively forcing all
	// other callers to treat this like modern OoenTofu does, and also how
	// things work in protocol version 6. A hypothetical provider relying
	// on the ability to "prepare" its configuration would not be usable
	// through this library, but such a provider has not been usable with
	// either OpenTofu or Terraform for quite some time either.
	protoResp, err := p.client.PrepareProviderConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateProviderConfigResponse{proto: protoResp}, nil
}

type validateProviderConfigResponse struct {
	proto *tfplugin5.PrepareProviderConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateProviderConfigResponse.
func (v validateProviderConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}

func (p *Provider) ConfigureProvider(ctx context.Context, req *providerops.ConfigureProviderRequest) (providerops.ConfigureProviderResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin5.Configure_Request{
		Config:             configVal,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.Configure(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return configureProviderResponse{proto: protoResp}, nil
}

type configureProviderResponse struct {
	proto *tfplugin5.Configure_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateProviderConfigResponse.
func (v configureProviderResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}
