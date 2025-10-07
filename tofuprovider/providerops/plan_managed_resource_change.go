package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type PlanManagedResourceChangeRequest struct {
	// ResourceType is the name of the type of resource the given state data
	// was created by.
	ResourceType string

	// PriorState is a dynamic value representing the state that was most
	// recently returned for this remote object.
	//
	// Providers expect that the given will conform to the current schema
	// for the given resource type. If the value might have been created
	// by an earlier version of the provider then callers must first use
	// the UpgradeResourceState operation to allow the provider to upgrade
	// the data to conform to the current version's schema.
	//
	// Clients should also typically use ReadManagedResource to detect any
	// changes that might have been made out of band. Although it's not
	// supposed to be required, some real-world providers malfunction if
	// given prior state that has not been refreshed in this way.
	PriorState providerschema.DynamicValueIn

	// PriorProviderInternal is the opaque blob that the provider returned in
	// the same response that produced the PriorState value, which the
	// provider might use to track some out-of-band information needed to
	// work with the object. Callers MUST populate this with exactly what
	// the provider most recently returned, or behavior is unspecified.
	PriorProviderInternal []byte

	// Config is a dynamic value representation of the object value
	// representing the resource instance's configuration.
	//
	// A value must be provided and its serialization type must be the implied
	// type of the schema given by this provider's
	// [providerschema.ProviderSchema.ManagedResourceTypeSchemas] method.
	Config providerschema.DynamicValueIn

	// ProposedNewState is an object representing the client's suggested
	// merge of PriorState and Config, which the provider is free to use
	// or ignore as it wishes.
	//
	// Unfortunately the presence of this field in the protocol means that
	// clients other than OpenTofu must emulate some business logic inside
	// OpenTofu to match how it would've populated this field. That behavior
	// is beyond the scope of this library because it lives at a higher level
	// of abstraction.
	ProposedNewState providerschema.DynamicValueIn

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

	// TODO: PriorIdentity
}

type PlanManagedResourceChangeResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// PlannedNewState returns an approximation of the new state value the
	// provider expects to produce if this change is applied.
	//
	// This will include "unknown value" placeholders in any location where
	// the final value cannot be predicted until after the change has been
	// applied.
	//
	// This must be decoded using the type implied by the schema of the
	// resource type.
	PlannedNewState() providerschema.DynamicValueOut

	// PlannedProviderInternal returns an opaque blob that must be sent
	// verbatim as part of a subsequent call to ApplyResourceChange,
	// containing any additional internal data the provider needs to track.
	PlannedProviderInternal() []byte

	// LegacyTypeSystem returns true if this provider is implemented using
	// the legacy SDK originally intended for now-obsolete versions of
	// Terraform, which cannot properly satisfy the requirements of the
	// modern OpenTofu resource instance change lifecycle.
	//
	// OpenTofu uses this to transform certain consistency errors into
	// log warnings instead, as a pragmatic way to keep old providers
	// working as well as possible until they have been updated. Other
	// clients can ignore this unless they intend to implement similar
	// provider behavior consistency checks.
	LegacyTypeSystem() bool

	// Deferred returns a non-nil value if the provider does not have enough
	// information to satisfy this request, such as if the provider
	// configuration is not yet known enough to know which API endpoint to
	// connect to.
	//
	// If this returns nil then other methods describe a change that could
	// potentially be applied.
	Deferred() Deferred

	// TODO: PlannedNewIdentity

	common.Sealed
}
