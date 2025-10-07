package providerschema

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"

	// For links in documentation comments:
	_ "maps"
)

// ProviderSchema represents the overall schema for an entire provider,
// describing all features that the provider offers for use by clients.
type ProviderSchema interface {
	// ProviderConfigSchema returns the schema for the provider's own
	// overal configuration, as used with the ConfigureProvider method.
	ProviderConfigSchema() Schema

	// ManagedResourceTypeSchemas returns an iterable sequence of the
	// schema for each managed resource type supported by this provider.
	//
	// The first result in each pair is the unique resource type name that
	// the schema belongs to. Use [maps.Collect] to gather the result into a
	// map from name to schema if you expect to need schemas for more than one
	// resource type.
	ManagedResourceTypeSchemas() iter.Seq2[string, Schema]

	// DataResourceTypeSchemas returns an iterable sequence of the
	// schema for each data resource type supported by this provider.
	//
	// The first result in each pair is the unique resource type name that
	// the schema belongs to. Use [maps.Collect] to gather the result into a
	// map from name to schema if you expect to need schemas for more than one
	// resource type.
	DataResourceTypeSchemas() iter.Seq2[string, Schema]

	// EphemeralResourceTypeSchemas returns an iterable sequence of the
	// schema for each ephemeral resource type supported by this provider.
	//
	// The first result in each pair is the unique resource type name that
	// the schema belongs to. Use [maps.Collect] to gather the result into a
	// map from name to schema if you expect to need schemas for more than one
	// resource type.
	EphemeralResourceTypeSchemas() iter.Seq2[string, Schema]

	// FunctionSignatures returns an iterable sequence of the signature of
	// each "provider-defined function" supported by this provider.
	//
	// The first result in each pair is the unique function type name that
	// the signature belongs to. Use [maps.Collect] to gather the result into a
	// map from name to signature if you expect to need signatures for more than
	// one function.
	FunctionSignatures() iter.Seq2[string, FunctionSignature]

	// ProviderMetaSchema returns the schema used for the rarely-used
	// "provider_meta" block type in the OpenTofu language, which allows
	// a module author to send module-related metadata with many different
	// provider requests related to objects in their module.
	//
	// Most callers should disregard this method. "Provider meta" is not
	// a widely-used provider protocol feature, and its corresponding
	// OpenTofu language features are not widely known in the community.
	ProviderMetaSchema() Schema

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
