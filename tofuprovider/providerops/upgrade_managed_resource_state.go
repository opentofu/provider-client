package providerops

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type UpgradeManagedResourceStateRequest struct {
	// ResourceType is the name of the type of resource the given state data
	// was created by.
	ResourceType string

	// SchemaVersion is the schema version for the given resource type that was
	// current when the provided raw state data was created.
	//
	// Callers must save that version number as part of their state snapshot
	// storage representation and provide it when upgrading in case the provider
	// needs to vary its behavior based on the starting schema version.
	//
	// Note however that not all upgrades are represented by schema version
	// changes: clients are required to call UpgradeManagedResourceState
	// even when the prior state's schema version matches the provider's
	// current schema version for this resource type.
	SchemaVersion int64

	// PrevStateRaw is the raw representation of the previously-saved state.
	//
	// Clients are not expected to have access to schema information for older
	// versions of a provider and so for this operation the client skips
	// trying to decode the data itself and instead assumes that the provider
	// knows how to decode data created by earlier versions of itself.
	PrevStateRaw providerschema.RawState
}

type UpgradeManagedResourceStateResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// UpgradedState is the upgraded state data.
	//
	// This must be decoded using the type implied by the current schema of the
	// resource type. Note that this operation does not return a
	// "ProviderInternal" blob itself, and clients are instead expected to
	// preserve whatever blob was associated with the state data that was
	// provided in the request.
	UpgradedState() providerschema.DynamicValueOut

	common.Sealed
}
