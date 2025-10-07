package providerschema

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"

	// For links in documentation comments:
	_ "maps"
)

// Schema describes a dynamic object schema, as used by various different
// provider features to represent configuration, returned results, or both.
//
// The schema structure exposes some details specific to OpenTofu's HCL-based
// language, including the distinction between "attributes" and "nested block
// types" that only really applies when using a schema to model configuration.
// When sending a value of a dynamic type derived from schema the attribute
// vs. block distinction is erased and that value is just a single object
// whose attributes are the superset of all of the attribute names and all
// of the nested block types.
type Schema interface {
	// SchemaVersion is the schema version number reported by the provider.
	//
	// Not all schema-based objects in the protocol actually make use of
	// schema information. It's primarily used for managed resource types
	// to drive their "schema upgrade" process. For schemas of "unversioned"
	// objects the result of this method is unspecified and meaningless.
	SchemaVersion() int64

	// Schema is a kind of [BlockType], which is a configuration-oriented
	// description of an object type in the OpenTofu language.
	BlockType

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// Attribute describes a single attribute within a [BlockType].
type Attribute interface {
	// Usage returns an enumeration value describing how this attribute can
	// be used across both the configuration and provider responses.
	Usage() AttributeUsage

	// Type returns the type constraint that any value assigned to this
	// attribute must conform to.
	//
	// Callers should typically call [Attribute.NestedType] first and
	// use this method only as a fallback if that function returns nil.
	// A correct provider should implement exactly one of these two
	// methods.
	Type() TypeConstraint

	// NestedType describes the required shape any value assigned to this
	// attribute, potentially including differing behavioral constraints
	// for each nested attribute.
	//
	// If this function returns nil, callers should typically call
	// [Attribute.Type] as a fallback to obtain a less sophiticated
	// representation of the requirements as a single type constraint.
	// A correct provider should implement exactly one of these two
	// methods.
	NestedType() ObjectType

	// If IsWriteOnly returns true then this attribute should appear only
	// in objects describing configuration, and should be null for
	// objects representing the object's state.
	//
	// This is meaningful only for the schema of a resource, where there
	// is an explicit distinction between "configuration" and "state" values.
	// Its result is meaningless and unspecified in other contexts.
	IsWriteOnly() bool

	// If IsSensitive returns true then the provider suggests that values
	// associated with this attribute not be displayed by default in any
	// human-oriented UI.
	IsSensitive() bool

	// DocDescription returns the provider's human-readable description
	// of the attribute. The second result describes the intended format for the
	// the description string.
	DocDescription() (string, DocStringFormat)

	// IsDeprecated returns true if the provider considers this attribute to
	// be deprecated.
	IsDeprecated() bool
}

// AttributeUsage is an enumeration describing how a particular attribute can
// be used across both the configuration and provider responses.
type AttributeUsage int

const (
	// AttributeUsageUnsupported represents that the provider returned an
	// attribute usage that this library does not understand.
	AttributeUsageUnsupported AttributeUsage = 0

	// AttributeRequired means that the attribute must be set to a non-null
	// value in the configuration and cannot be overridden by the provider
	// at all.
	AttributeRequired AttributeUsage = 1

	// AttributeRequired means that the attribute may be set in the
	// configuration, but if it is not set then it defaults to null.
	AttributeOptional AttributeUsage = 2

	// AttributeOptionalComputed means that the attribute may be set in the
	// configuration, but if (and only if) it is not set then the provider
	// will provide a value to be used as a default.
	//
	// This usage is allowed only for schemas representing object types that
	// can be returned in the result of a provider call.
	AttributeOptionalComputed AttributeUsage = 3

	// AttributeComputed means that the attribute may not be set in the
	// configuration and its value is decided exclusively by the provider in
	// all cases.
	//
	// This usage is allowed only for schemas representing object types that
	// can be returned in the result of a provider call.
	AttributeComputed AttributeUsage = 4
)

// ObjectType describes an object type to be used with an [Attribute],
// or possibly a collection of objects of that type depending on the nesting
// mode.
type ObjectType interface {
	// Nesting returns the nesting mode for this nested object type.
	//
	// The values returned from this function influence whether and how
	// the described object type is wrapped into a collection type.
	Nesting() NestingMode

	// Attributes returns an iterable sequence of the expected attributes
	// in this object type.
	//
	// The first result of each item is the unique attribute name. Use
	// [maps.Collect] to produce a map from attribute name to definition.
	Attributes() iter.Seq2[string, Attribute]

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// Block is implemented by both [Schema] and [NestedBlockType] to describe the
// features that both top-level blocks and nested blocks have in common.
type BlockType interface {
	// Attributes returns an iterable sequence of the expected attributes
	// in this block type.
	//
	// The first result of each item is the unique attribute name. Use
	// [maps.Collect] to produce a map from attribute name to definition.
	Attributes() iter.Seq2[string, Attribute]

	// NestedBlockTypes returns an iterable sequence of child block types
	// that are allowed to nest inside this block type.
	//
	// The first result of each item is the unique block type name. Use
	// [maps.Collect] to produce a map from attribute name to definition.
	//
	// In a valid provider schema the keys from [BlockType.Attributes] and
	// the keys from [BlockType.NestedBlockTypes] are disjoint; there should
	// never be an attribute and a nested block type of the same name.
	NestedBlockTypes() iter.Seq2[string, NestedBlockType]

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

type NestingMode int

const (
	// NestingInvalid is the zero value of [NestingMode], used when a provider
	// returns a value that this library does not support.
	NestingInvalid NestingMode = 0

	// NestingSingle means that the nested block or object is just a single
	// object that may be null.
	NestingSingle NestingMode = 1

	// NestingList means that nested block or object is represented as a
	// list of the object type implied by the nested block's schema.
	NestingList NestingMode = 2

	// NestingList means that nested block or object is represented as a
	// set of the object type implied by the nested block's schema.
	NestingSet NestingMode = 3

	// NestingList means that nested block or object is represented as a
	// map of the object type implied by the nested block's schema.
	NestingMap NestingMode = 4

	// NestingGroup is a slight variant of [NestingSingle] where the client
	// is expected to never send a null object directly and should instead
	// construct a non-null object whose attributes are all set to null.
	//
	// When constructing the synthetic object to represent the not-present
	// state, the client may ignore the "requiredness" of nested attributes,
	// because nested attributes are required only when a block is written
	// out explicitly in the configuration.
	NestingGroup NestingMode = 5
)

// NestedBlockType describes a nested block type that can appear inside
// another [BlockType].
type NestedBlockType interface {
	// Nesting returns the nesting mode for this nested block type.
	//
	// The values returned from this function influence whether and how
	// the described object type is wrapped into a collection type.
	Nesting() NestingMode

	// NestedBlockType is a kind of [BlockType], which is a
	// configuration-oriented description of an object type in the OpenTofu
	// language.
	//
	// For NestedBlockType this describes the object type of each instance
	// of this block type. When multiple blocks of the same type are supported
	// they are collected into some sort of aggregate type depending on
	// [NestedBlockType.NestingMode].
	BlockType

	// ItemLimits returns the minimum and maximum number of nested objects
	// that are allowed, respectively.
	//
	// Item limits are a legacy concept not actually enforced by modern OpenTofu.
	// Providers are instead expected to enforce any limits as part of the
	// normal validation process, just like any other non-type-based
	// constraint on which values can be provided. However, some providers
	// may still return this information for documentation purposes. Most
	// callers should completely disregard this method, since its results are
	// not reliable.
	//
	// If the maximum number is returned as zero, that represents that there
	// is no limit on the number of items allowed. If the minimum number is
	// one then it represents that at least one item is required. Minimum
	// values other than zero or one are not meaningful in the modern protocol
	// and should be treated the same as returning one.
	//
	// Item limits are not meaningful for [NestingSingle] and [NestingGroup],
	// because those nesting modes inherently imply a maximum of one item.
	ItemLimits() (int64, int64)

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
