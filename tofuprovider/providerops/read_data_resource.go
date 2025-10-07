package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type ReadDataResourceRequest struct {
	// ResourceType is the name of the type of resource the client wants to
	// read.
	ResourceType string

	// Config is a dynamic value representing the resource configuration.
	//
	// Providers expect that the given will conform to the current schema
	// for the given resource type.
	Config providerschema.DynamicValueIn

	// ProviderMeta is some additional metadata declared in the module where
	// this resource was declared. This is a rarely-used feature that most
	// callers should ignore, leaving this field completely unassigned.
	// When populated the value must be of the type implied by the provider's
	// ProviderMetaSchema.
	ProviderMeta providerschema.DynamicValueIn

	// ClientCapabilities allows the caller to declare that it is capable of
	// handling certain response data that was added to the protocol after
	// it was initially defined, and thus which the provider must disable
	// by default to avoid confusing older clients.
	ClientCapabilities *ClientCapabilities
}

type ReadDataResourceResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// State is the resulting state data.
	//
	// This must be decoded using the type implied by the schema of the
	// resource type.
	State() providerschema.DynamicValueOut

	// Deferred returns a non-nil value if the provider does not have enough
	// information to satisfy this request, such as if the provider
	// configuration is not yet known enough to know which API endpoint to
	// connect to.
	//
	// If this returns nil then other methods return updated data that
	// should replace the previous values that were saved in the prior state.
	Deferred() Deferred

	common.Sealed
}
