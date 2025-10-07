package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type ConfigureProviderRequest struct {
	// TerraformVersion specifies which version of Terraform is configuring the
	// provider.
	//
	// This part of the protocol was obviously not designed with non-Terraform
	// callers in mind and so there is no defined way for other callers to
	// populate this field. Callers may wish to specify a version of Terraform
	// that roughly matches the set of functionality they intend to use, or
	// can leave this completely unspecified which providers will treat as
	// being called by older versions of Terraform from before this part of
	// the protocol was added.
	TerraformVersion string

	// Config is a dynamic value representation of the provider configuration.
	//
	// The serialization type of the value must be the type implied by the
	// provider's configuration schema, as returned in
	// [providerschema.ProviderSchema.ProviderConfigSchema()].
	//
	// The value passed here should have previously been passed in the similar
	// field of [ValidateProviderConfigRequest] to ensure that the provider
	// considers it valid. Providers will not necessarily repeat the same
	// validation when asked to configure, and so passing an unvalidated
	// configuration causes unspecified behavior.
	Config providerschema.DynamicValueIn

	// ClientCapabilities allows the caller to declare that it is capable of
	// handling certain response data that was added to the protocol after
	// it was initially defined, and thus which the provider must disable
	// by default to avoid confusing older clients.
	ClientCapabilities *ClientCapabilities
}

type ConfigureProviderResponse interface {
	// Diagnostics describe any problems the provider reported with the
	// provided configuration.
	//
	// If this includes any error diagnostics then passing the configuration
	// object is somehow invalid and so passing it to ConfigureProvider causes
	// unspecified behavior.
	Diagnostics() Diagnostics

	common.Sealed
}
