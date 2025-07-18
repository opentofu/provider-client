package providerschema

import (
	"iter"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"

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
	// to drive their "schema upgrade" process.
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
}

// ObjectType describes an object type to be used with an [Attribute],
// or possibly a collection of objects of that type depending on the nesting
// mode.
type ObjectType interface {
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

// NestedBlockType describes a nested block type that can appear inside
// another [BlockType].
type NestedBlockType interface {

	// NestedBlockType is a kind of [BlockType], which is a
	// configuration-oriented description of an object type in the OpenTofu
	// language.
	//
	// For NestedBlockType this describes the object type of each instance
	// of this block type. When multiple blocks of the same type are supported
	// they are collected into some sort of aggregate type depending on
	// [NestedBlockType.NestingMode].
	BlockType

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
