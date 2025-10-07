package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type ValidateProviderConfigRequest struct {
	// Config is a dynamic value representation of the object value
	// representing the provider configuration.
	//
	// A value must be provided and its serialization type must be the implied
	// type of the schema given in by this provider's
	// [providerschema.ProviderSchema.ProviderConfigSchema] method.
	Config providerschema.DynamicValueIn
}

type ValidateProviderConfigResponse interface {
	// Diagnostics describe any problems the provider reported with the
	// provided configuration.
	//
	// If this includes any error diagnostics then passing the configuration
	// object is somehow invalid and so passing it to ConfigureProvider causes
	// unspecified behavior.
	Diagnostics() Diagnostics

	common.Sealed
}
