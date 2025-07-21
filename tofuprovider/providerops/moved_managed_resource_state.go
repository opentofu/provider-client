package providerops

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type MoveManagedResourceStateRequest struct {
	// SourceProviderAddress is the full provider source address of the
	// provider that the source resource type belonged to.
	SourceProviderAddress string

	// SourceResourceType is the name of the type of resource the given state
	// data was created by. This should be a resource type supported by some
	// existing version of the provider given in SourceProviderAddress.
	SourceResourceType string

	// SourceSchemaVersion is the schema version for the given resource type
	// that was current when the provided raw state data was created.
	//
	// Callers must save that version number as part of their state snapshot
	// storage representation and provide it when upgrading in case the provider
	// needs to vary its behavior based on the starting schema version.
	SourceSchemaVersion int64

	// SourceStateRaw is the raw representation of the source object's state.
	//
	// Clients are not expected to have access to schema information the source
	// provider and so for this operation the client skips trying to decode the
	// data itself and instead assumes that the provider knows how to decode
	// data created by any provider it intends to support migration from.
	SourceStateRaw providerschema.RawState

	// SourceProviderInternal is the ProviderInternal value that was returned
	// along with the state data passed in SourceStateRaw.
	SourceProviderInternal []byte

	// TargetResourceType is the name of a resource type supported by the
	// provider that's handling this request whose schema the source data
	// should be converted to.
	TargetResourceType string

	// TODO: SourceIdentity and SourceIdentitySchemaVersion
}

type MoveManagedResourceStateResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// TargetState returns the new state data, of a suitable structure for
	// the target resource type.
	//
	// Refer to [ApplyManagedResourceChangeResponse.NewState] for details on
	// how to use this.
	TargetState() providerschema.DynamicValueOut

	// TargetProviderInternal returns an opaque blob that must be saved
	// along with the new state data.
	//
	// Refer to [ApplyManagedResourceChangeResponse.ProviderInternal] for
	// details on how to use this.
	TargetProviderInternal() []byte

	common.Sealed
}
