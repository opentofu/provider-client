package providerdefs

import (
	"context"

	"github.com/opentofu/provider-client/tofuprovider"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// Provider is a base implementation of [tofuprovider.Provider] where most
// methods immediately fail with an error that passes
// [providerops.IsUnimplementedErr].
//
// Embed this into any external implementation of the interface so that your
// type will automatically have default implementations of any new methods
// added in later versions of this library.
type Provider struct {
	sealedImpl
}

var _ tofuprovider.Provider = Provider{}

// ApplyManagedResourceChange implements [tofuprovider.Provider].
func (Provider) ApplyManagedResourceChange(ctx context.Context, req *providerops.ApplyManagedResourceChangeRequest) (providerops.ApplyManagedResourceChangeResponse, error) {
	return nil, common.ErrUnimplemented
}

// CallFunction implements [tofuprovider.Provider].
func (Provider) CallFunction(ctx context.Context, req *providerops.CallFunctionRequest) (providerops.CallFunctionResponse, error) {
	return nil, common.ErrUnimplemented
}

// CloseEphemeralResource implements [tofuprovider.Provider].
func (Provider) CloseEphemeralResource(ctx context.Context, req *providerops.CloseEphemeralResourceRequest) (providerops.CloseEphemeralResourceResponse, error) {
	return nil, common.ErrUnimplemented
}

// ConfigureProvider implements [tofuprovider.Provider].
func (Provider) ConfigureProvider(ctx context.Context, req *providerops.ConfigureProviderRequest) (providerops.ConfigureProviderResponse, error) {
	return nil, common.ErrUnimplemented
}

// GetFunctions implements [tofuprovider.Provider].
func (Provider) GetFunctions(ctx context.Context, req *providerops.GetFunctionsRequest) (providerops.GetFunctionsResponse, error) {
	return nil, common.ErrUnimplemented
}

// GetProviderSchema implements [tofuprovider.Provider].
func (Provider) GetProviderSchema(ctx context.Context, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error) {
	return nil, common.ErrUnimplemented
}

// GracefulStop implements [tofuprovider.Provider] by doing absolutely nothing.
func (Provider) GracefulStop(ctx context.Context) error {
	// Doing nothing at all is a valid implementation of GracefulStop, so
	// we'll do that as the default behavior.
	return nil
}

// ImportManagedResourceState implements [tofuprovider.Provider].
func (Provider) ImportManagedResourceState(ctx context.Context, req *providerops.ImportManagedResourceStateRequest) (providerops.ImportManagedResourceStateResponse, error) {
	return nil, common.ErrUnimplemented
}

// MoveManagedResourceState implements [tofuprovider.Provider].
func (Provider) MoveManagedResourceState(ctx context.Context, req *providerops.MoveManagedResourceStateRequest) (providerops.MoveManagedResourceStateResponse, error) {
	return nil, common.ErrUnimplemented
}

// OpenEphemeralResource implements [tofuprovider.Provider].
func (Provider) OpenEphemeralResource(ctx context.Context, req *providerops.OpenEphemeralResourceRequest) (providerops.OpenEphemeralResourceResponse, error) {
	return nil, common.ErrUnimplemented
}

// PlanManagedResourceChange implements [tofuprovider.Provider].
func (Provider) PlanManagedResourceChange(ctx context.Context, req *providerops.PlanManagedResourceChangeRequest) (providerops.PlanManagedResourceChangeResponse, error) {
	return nil, common.ErrUnimplemented
}

// ReadDataResource implements [tofuprovider.Provider].
func (Provider) ReadDataResource(ctx context.Context, req *providerops.ReadDataResourceRequest) (providerops.ReadDataResourceResponse, error) {
	return nil, common.ErrUnimplemented
}

// ReadManagedResource implements [tofuprovider.Provider].
func (Provider) ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error) {
	return nil, common.ErrUnimplemented
}

// RenewEphemeralResource implements [tofuprovider.Provider].
func (Provider) RenewEphemeralResource(ctx context.Context, req *providerops.RenewEphemeralResourceRequest) (providerops.RenewEphemeralResourceResponse, error) {
	return nil, common.ErrUnimplemented
}

// UpgradeManagedResourceState implements [tofuprovider.Provider].
func (Provider) UpgradeManagedResourceState(ctx context.Context, req *providerops.UpgradeManagedResourceStateRequest) (providerops.UpgradeManagedResourceStateResponse, error) {
	return nil, common.ErrUnimplemented
}

// ValidateDataResourceConfig implements [tofuprovider.Provider].
func (Provider) ValidateDataResourceConfig(ctx context.Context, req *providerops.ValidateDataResourceConfigRequest) (providerops.ValidateDataResourceConfigResponse, error) {
	return nil, common.ErrUnimplemented
}

// ValidateEphemeralResourceConfig implements [tofuprovider.Provider].
func (Provider) ValidateEphemeralResourceConfig(ctx context.Context, req *providerops.ValidateEphemeralResourceConfigRequest) (providerops.ValidateEphemeralResourceConfigResponse, error) {
	return nil, common.ErrUnimplemented
}

// ValidateManagedResourceConfig implements [tofuprovider.Provider].
func (Provider) ValidateManagedResourceConfig(ctx context.Context, req *providerops.ValidateManagedResourceConfigRequest) (providerops.ValidateManagedResourceConfigResponse, error) {
	return nil, common.ErrUnimplemented
}

// ValidateProviderConfig implements [tofuprovider.Provider].
func (Provider) ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error) {
	return nil, common.ErrUnimplemented
}

// ListManagedResources implements [tofuprovider.Provider].
func (Provider) ListManagedResources(ctx context.Context, req *providerops.ListManagedResourcesRequest) (providerops.ListManagedResourcesResponse, error) {
	return nil, common.ErrUnimplemented
}

// GRPCPluginProvider is a base implementation of
// [tofuprovider.GRPCPluginProvider] intended to be embedded in external
// implementations of that interface, providing default implementations of
// all of the required methods.
//
// This type embeds [Provider], inheriting its implementation of
// [tofuprovider.Provider]. Use this base type only if your external
// implementation for some reason needs to act like a gRPC provider in
// particular, such as if you're mocking to test some code that manages
// gRPC-based providers. If you're implementing some other kind of provider then
// you should just embed [Provider] to implement the base interface, rather than
// including is type's useless implementations of
// [tofuprovider.GRPCPluginProvider].
type GRPCPluginProvider struct {
	Provider
}

var _ tofuprovider.GRPCPluginProvider = GRPCPluginProvider{}

// ClientProxy implements [tofuprovider.GRPCPluginProvider] by returning a
// value that is not assertable to any known gRPC client proxy interface.
func (GRPCPluginProvider) ClientProxy() any {
	return struct{}{}
}

// Close implements [tofuprovider.GRPCPluginProvider] by doing absolutely
// nothing.
func (GRPCPluginProvider) Close() error {
	return nil
}
