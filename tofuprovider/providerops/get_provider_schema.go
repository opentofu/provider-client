package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type GetProviderSchemaRequest struct {
	// There are currently no arguments in a provider schema request.
}

type GetProviderSchemaResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// ServerCapabilities returns an object describing various special
	// capabilities the provider claims to have, intended for use as part
	// of client/server capability negotiation.
	//
	// Callers MUST respect the server's reported capabilities or else
	// the provider is likely to malfunction in unexpected ways.
	ServerCapabilities() ServerCapabilities

	// ProviderSchema returns the overall provider schema, describing all
	// of the features that the provider claims to offer and how clients
	// are expected to interact with those features.
	ProviderSchema() providerschema.ProviderSchema

	common.Sealed
}
