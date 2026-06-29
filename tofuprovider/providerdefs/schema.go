package providerdefs

import (
	"iter"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// Schema is a base implementation of [providerschema.Schema] whose default
// methods describe a completely empty schema of version zero.
type Schema struct {
	BlockType
}

var _ providerschema.Schema = Schema{}

// DocDescription implements [providerschema.Schema].
func (Schema) DocDescription() (string, providerschema.DocStringFormat) {
	return "", providerschema.DocStringPlain
}

// SchemaVersion implements [providerschema.Schema].
func (Schema) SchemaVersion() int64 {
	return 0
}

// Attribute is a base implementation of [providerschema.Attribute].
//
// If you embed this in your own implementation of the interface, you must still
// implement [providerschema.Attribute.Usage] and exactly either
// [providerschema.Attribute.Type] or [providerschema.Attribute.NestedType]
// (but not both) to get a correctly-functioning implementation of the
// interface.
type Attribute struct {
	sealedImpl
}

// NestedType implements [providerschema.Attribute] by returning nil.
//
// Note that nil is only a valid result if the "Type" method returns a valid,
// non-nil type constraint, so downstream implementers MUST implement either
// NestedType.
func (Attribute) NestedType() providerschema.ObjectType {
	return nil
}

// Type implements [providerschema.Attribute] by returning nil.
//
// Note that nil is only a valid result if the "NestedType" method returns a
// valid, non-nil object type, so downstream implementers MUST implement either
// Type or NestedType.
func (Attribute) Type() providerschema.TypeConstraint {
	return nil
}

// DocDescription implements [providerschema.Attribute].
func (Attribute) DocDescription() (string, providerschema.DocStringFormat) {
	return "", providerschema.DocStringPlain
}

// IsDeprecated implements [providerschema.Attribute].
func (Attribute) IsDeprecated() bool {
	return false
}

// IsSensitive implements [providerschema.Attribute].
func (Attribute) IsSensitive() bool {
	return false
}

// IsWriteOnly implements [providerschema.Attribute].
func (Attribute) IsWriteOnly() bool {
	return false
}

// BlockType is a base implementation of [providerschema.BlockType].
type BlockType struct {
	sealedImpl
}

var _ providerschema.BlockType = BlockType{}

// Attributes implements [providerschema.BlockType] by reporting no attributes
// at all.
func (b BlockType) Attributes() iter.Seq2[string, providerschema.Attribute] {
	return func(func(string, providerschema.Attribute) bool) {}
}

// NestedBlockTypes implements [providerschema.BlockType] by reporting no
// nested block types at all.
func (b BlockType) NestedBlockTypes() iter.Seq2[string, providerschema.NestedBlockType] {
	return func(func(string, providerschema.NestedBlockType) bool) {}
}

// NestedBlockType is a base implementation of [providerschema.NestedBlockType].
//
// Implementers embedding this type must still provide their own implementation
// of [providerschema.NestedBlockType.Nesting] to describe the nested block
// type's nesting mode.
type NestedBlockType struct {
	BlockType
}

// ItemLimits implements [providerschema.NestedBlockType] by reporting no
// item limits at all.
func (n NestedBlockType) ItemLimits() (int64, int64) {
	return 0, 0
}

// ObjectType is a base implementation of [providerschema.ObjectType].
//
// Implementers embedding this type must still provide their own implementation
// of [providerschema.ObjectType.Nesting] to describe the nesting mode.
type ObjectType struct {
	sealedImpl
}

// Attributes implements [providerschema.ObjectType] by reporting no attributes
// at all.
func (o ObjectType) Attributes() iter.Seq2[string, providerschema.Attribute] {
	return func(func(string, providerschema.Attribute) bool) {}
}

// TypeConstraint is a base implementation of [providerschema.TypeConstraint].
//
// Implementers embedding this type must still provide their own implementation
// of [providerschema.TypeConstraint.AsCtyType].
type TypeConstraint struct {
	sealedImpl
}

// NOTE: Right now there aren't actually any other methods of TypeConstraint
// to implement here, but external implementers are expected to embed this
// type anyway so that we can potentially add new methods to the interface
// in future and then write their default implementations here. If you're here
// reading this comment because that's what you're about to do then you can
// remove this comment at the same time!

// TypeConstraint is a base implementation of [providerschema.ProviderSchema].
type ProviderSchema struct {
	sealedImpl
}

var _ providerschema.ProviderSchema = ProviderSchema{}

// DataResourceTypeSchemas implements [providerschema.ProviderSchema] by
// reporting no resource types at all.
func (p ProviderSchema) DataResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return func(func(string, providerschema.Schema) bool) {}
}

// EphemeralResourceTypeSchemas implements [providerschema.ProviderSchema] by
// reporting no resource types at all.
func (p ProviderSchema) EphemeralResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return func(func(string, providerschema.Schema) bool) {}
}

// FunctionSignatures implements [providerschema.ProviderSchema] by reporting
// no callable functions at all.
func (p ProviderSchema) FunctionSignatures() iter.Seq2[string, providerschema.FunctionSignature] {
	return func(func(string, providerschema.FunctionSignature) bool) {}
}

// ManagedResourceTypeSchemas implements [providerschema.ProviderSchema] by
// reporting no resource types at all.
func (p ProviderSchema) ManagedResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return func(func(string, providerschema.Schema) bool) {}
}

// ManagedResourceTypeListSchemas implements [providerschema.ProviderSchema] by
// reporting no resource types at all.
func (p ProviderSchema) ManagedResourceTypeListSchemas() iter.Seq2[string, providerschema.Schema] {
	return func(func(string, providerschema.Schema) bool) {}
}

// ProviderConfigSchema implements [providerschema.ProviderSchema] by reporting
// an empty schema.
func (p ProviderSchema) ProviderConfigSchema() providerschema.Schema {
	return Schema{}
}

// ProviderMetaSchema implements [providerschema.ProviderSchema] by returning
// nil to indicate that this provider doesn't support "provider meta" at all.
//
// (Most providers should not implement this, because "provider meta" is a niche
// feature that only serves the very narrow case of a module being written by
// the same author as the main provider it uses and using that provider as a
// way to collect usage data for the module.)
func (p ProviderSchema) ProviderMetaSchema() providerschema.Schema {
	return nil
}

// FunctionSignature is a base implementation of [providerschema.FunctionSignature].
//
// External implementers must still provide their own implementation of
// [providerschema.FunctionSignature.ResultType].
type FunctionSignature struct {
	sealedImpl
}

// DeprecationMessage implements [providerschema.FunctionSignature] by reporting
// no deprecation at all.
func (f FunctionSignature) DeprecationMessage() string {
	return ""
}

// DocDescription implements [providerschema.FunctionSignature] by returning
// no documentation.
func (f FunctionSignature) DocDescription() (string, providerschema.DocStringFormat) {
	return "", providerschema.DocStringPlain
}

// DocSummary implements [providerschema.FunctionSignature] by returning no
// summary.
func (f FunctionSignature) DocSummary() string {
	return ""
}

// Parameters implements [providerschema.FunctionSignature] by reporting that
// there aren't any parameters.
func (f FunctionSignature) Parameters() iter.Seq[providerschema.FunctionParameter] {
	return func(func(providerschema.FunctionParameter) bool) {}
}

// VariadicParameter implements [providerschema.FunctionSignature] by reporting
// that there is no variadic parameter.
func (f FunctionSignature) VariadicParameter() providerschema.FunctionParameter {
	return nil
}

// ===================== EXAMPLE IMPLEMENTATIONS ==============================
//
// The remaining unexported types are here just to make sure that a third-party
// implementation that already has these methods can be made to implement
// the corresponding interfaces just by embedding our default types from above.
//
// If future additions to the interface types cause failures here then the
// appropriate fix is to add default implementations of those methods to the
// exported types, _not_ to add those methods to these unexported types.

type attributeExample struct{ Attribute }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerschema.Attribute = attributeExample{}

func (attributeExample) Usage() providerschema.AttributeUsage { panic("unimplemented") }

type nestedBlockTypeExample struct{ NestedBlockType }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerschema.NestedBlockType = nestedBlockTypeExample{}

func (nestedBlockTypeExample) Nesting() providerschema.NestingMode { panic("unimplemented") }

type objectTypeExample struct{ ObjectType }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerschema.ObjectType = objectTypeExample{}

func (objectTypeExample) Nesting() providerschema.NestingMode { panic("unimplemented") }

type typeConstraintExample struct{ TypeConstraint }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerschema.TypeConstraint = typeConstraintExample{}

func (typeConstraintExample) AsCtyType() (cty.Type, error) { panic("unimplemented") }

type functionSignatureExample struct{ FunctionSignature }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerschema.FunctionSignature = functionSignatureExample{}

func (functionSignatureExample) ResultType() providerschema.TypeConstraint { panic("unimplemented") }
