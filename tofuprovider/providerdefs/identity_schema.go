package providerdefs

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// IdentitySchema is a base implementation of [providerschema.IdentitySchema].
type IdentitySchema struct {
	sealedImpl
}

var _ providerschema.IdentitySchema = IdentitySchema{}

// SchemaVersion implements [providerschema.IdentitySchema] by reporting version
// zero.
func (IdentitySchema) SchemaVersion() int64 {
	return 0
}

// Attributes implements [providerschema.IdentitySchema] by reporting no
// attributes at all.
func (IdentitySchema) Attributes() iter.Seq2[string, providerschema.IdentityAttribute] {
	return func(func(string, providerschema.IdentityAttribute) bool) {}
}

// IdentityAttribute is a base implementation of
// [providerschema.IdentityAttribute].
//
// If you embed this in your own implementation of the interface, you must still
// override [providerschema.IdentityAttribute.Type] to return a valid, non-nil
// type constraint.
type IdentityAttribute struct {
	sealedImpl
}

var _ providerschema.IdentityAttribute = IdentityAttribute{}

// Type implements [providerschema.IdentityAttribute] by returning nil.
//
// nil is not a valid result, so embedders MUST override this method with a real
// type constraint.
func (IdentityAttribute) Type() providerschema.TypeConstraint {
	return nil
}

// DocDescription implements [providerschema.IdentityAttribute].
func (IdentityAttribute) DocDescription() (string, providerschema.DocStringFormat) {
	return "", providerschema.DocStringPlain
}

// ImportUsage implements [providerschema.IdentityAttribute] by reporting that
// the attribute is not usable during import.
func (IdentityAttribute) ImportUsage() providerschema.AttributeUsage {
	return providerschema.AttributeUsageUnsupported
}
