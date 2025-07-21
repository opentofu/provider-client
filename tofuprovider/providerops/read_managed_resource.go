package providerops

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type ReadManagedResourceRequest struct {
	// ResourceType is the name of the type of resource the given state data
	// was created by.
	ResourceType string

	// CurrentState is a dynamic value representing the state that was most
	// recently returned for this remote object.
	//
	// Providers expect that the given will conform to the current schema
	// for the given resource type. If the value might have been created
	// by an earlier version of the provider then callers must first use
	// the UpgradeResourceState operation to allow the provider to upgrade
	// the data to conform to the current version's schema.
	CurrentState providerschema.DynamicValueIn

	// ProviderInternal is the opaque blob that the provider returned in
	// the same response that produced the CurrentState value, which the
	// provider might use to track some out-of-band information needed to
	// work with the object. Callers MUST populate this with exactly what
	// the provider most recently returned, or behavior is unspecified.
	ProviderInternal []byte

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

	// TODO: CurrentIdentity
}

type ReadManagedResourceResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// NewState is the updated state data.
	//
	// This must be decoded using the type implied by the schema of the
	// resource type.
	NewState() providerschema.DynamicValueOut

	// ProviderInternal is the updated opaque blob to be saved alongside
	// the new state data and provided alongside it in any future provider
	// call.
	ProviderInternal() []byte

	// Deferred returns a non-nil value if the provider does not have enough
	// information to satisfy this request, such as if the provider
	// configuration is not yet known enough to know which API endpoint to
	// connect to.
	//
	// If this returns nil then other methods return updated data that
	// should replace the previous values that were saved in the prior state.
	Deferred() Deferred

	// TODO: NewIdentity

	common.Sealed
}
