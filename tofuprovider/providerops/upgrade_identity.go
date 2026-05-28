package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type UpgradeIdentityRequest struct {
	// ResourceType is the name of the managed resource type whose identity
	// data is being upgraded.
	ResourceType string

	// SchemaVersion is the identity schema version that was current when
	// the provided raw identity data was created.
	//
	// Note that this is the identity schema version, distinct from the
	// resource's main schema version. Callers must persist this number
	// alongside any stored raw identity data so that the provider can
	// interpret it correctly when upgrading.
	SchemaVersion int64

	// PrevIdentityJSON is the previously-saved identity data as JSON bytes,
	// in whatever format the provider produced under an earlier version of
	// its identity schema.
	//
	// The client cannot decode this value because it does not have access
	// to the older identity schema; only the provider knows how to
	// interpret it. Identity data is always JSON-encoded over the wire,
	// so this is exposed as raw bytes rather than via [providerschema.RawState].
	PrevIdentityJSON []byte
}

type UpgradeIdentityResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// UpgradedIdentity is the upgraded resource identity data, suitable
	// for decoding against the provider's current identity schema (as
	// returned by [tofuprovider.Provider.GetIdentitySchemas]).
	//
	// In a successful response this should always be non-nil. A nil
	// result without an accompanying error diagnostic indicates a
	// misbehaving provider.
	UpgradedIdentity() providerschema.DynamicValueOut

	common.Sealed
}
