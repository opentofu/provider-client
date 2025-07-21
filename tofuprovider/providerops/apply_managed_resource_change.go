package providerops

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type ApplyManagedResourceChangeRequest struct {
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

	// Config is a dynamic value representation of the object value
	// representing the resource instance's configuration.
	//
	// A value must be provided and its serialization type must be the implied
	// type of the schema given by this provider's
	// [providerschema.ProviderSchema.ManagedResourceTypeSchemas] method.
	Config providerschema.DynamicValueIn

	// PlannedNewState is the same value returned from the PlannedNewState
	// method in the response from PlanManagedResourceChange that this call
	// is intending to apply.
	//
	// Callers must not modify the value that the provider planned, or
	// provider behavior is unspecified.
	PlannedNewState providerschema.DynamicValueIn

	// PlannedProviderInternal is the opaque blob that the provider returned in
	// the PlanManagedResourceChange response that this call is intending to
	// apply.
	//
	// Callers must preserve this value exactly or provider behavior is
	// unspecified.
	PlannedProviderInternal []byte

	// ProviderMeta is some additional metadata declared in the module where
	// this resource was declared. This is a rarely-used feature that most
	// callers should ignore, leaving this field completely unassigned.
	// When populated the value must be of the type implied by the provider's
	// ProviderMetaSchema.
	ProviderMeta providerschema.DynamicValueIn

	// TODO: PlannedNewIdentity
}

type ApplyManagedResourceChangeResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// NewState describes the updated state of the corresponding remote object.
	//
	// This must be decoded using the type implied by the schema of the
	// resource type. If a caller intends to work with the same object again
	// on a future upgrade+refresh+plan+apply round then it must pass the decoded
	// result to [providerschema.NewRawState] and store that result alongside
	// the result of the ProviderInternal method and the current schema version
	// number associated with the requested resource type.
	PlannedNewState() providerschema.DynamicValueOut

	// ProviderInternal returns an opaque blob that must be sent
	// verbatim as part of a subsequent call to ReadManagedResource or
	// PlanManagedResourceChange, containing any additional internal data the
	// provider needs to track.
	ProviderInternal() []byte

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

	// TODO: NewIdentity

	common.Sealed
}
