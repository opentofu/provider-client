package providerschema

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"

	// For links in documentation comments:
	_ "maps"
)

// IdentitySchema describes the structure of the data used to identify a
// single instance of a managed resource type.
//
// Resource identity is a versioned, JSON-encoded object that uniquely
// identifies a remote object managed by the provider. Unlike [Schema],
// an identity schema only supports a flat set of named attributes:
// there are no nested blocks or nested object types.
type IdentitySchema interface {
	// SchemaVersion is the identity schema version number reported by the
	// provider.
	//
	// This is distinct from the resource type's main schema version
	// (as returned by [Schema.SchemaVersion]). Callers must persist this
	// number alongside any stored raw identity data so that the provider
	// can interpret it correctly when upgrading.
	SchemaVersion() int64

	// Attributes returns an iterable sequence of the attributes that
	// make up this identity schema.
	//
	// The first result of each item is the unique attribute name. Use
	// [maps.Collect] to produce a map from attribute name to definition.
	Attributes() iter.Seq2[string, IdentityAttribute]

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// IdentityAttribute describes a single attribute within an [IdentitySchema].
type IdentityAttribute interface {
	// Type returns the type constraint that any value assigned to this
	// attribute must conform to.
	Type() TypeConstraint

	// DocDescription returns the provider's human-readable description
	// of the attribute. The second result describes the intended format
	// for the description string.
	DocDescription() (string, DocStringFormat)

	// IsRequiredForImport returns true if a caller using this attribute
	// to identify a resource for import must supply a value for it.
	IsRequiredForImport() bool

	// IsOptionalForImport returns true if a caller using this attribute
	// to identify a resource for import may supply a value for it but
	// may also omit it.
	IsOptionalForImport() bool

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
