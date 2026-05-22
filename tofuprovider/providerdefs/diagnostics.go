package providerdefs

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// Diagnostics is a base implementation of [providerops.Diagnostics]
// intended to be embedded into implementations outside of this module.
//
// Implementations embedding this type must still provide their own
// implementations of [providerops.Diagnostics.HasErrors] and
// [providerops.Diagnostics.All], but any new methods added to the interface
// type in future will have default implementations added to this base
// type to ensure forward-compatibility.
type Diagnostics struct {
	sealedImpl
}

// emptyDiagnostics is an implementation of [providerops.Diagnostics] that
// we use as part of default implementations of methods elsewhere that return
// diagnostics, causting them to default to returning no diagnostics at all.
//
// This is unexported to discourage external implementers from just using
// this "empty" implementation instead of implementing one that actually works.
type emptyDiagnostics struct {
	Diagnostics
}

var _ providerops.Diagnostics = emptyDiagnostics{}

// All implements [providerops.Diagnostics].
func (e emptyDiagnostics) All() iter.Seq[providerops.Diagnostic] {
	return func(func(providerops.Diagnostic) bool) {}
}

// HasErrors implements [providerops.Diagnostics].
func (e emptyDiagnostics) HasErrors() bool {
	return false
}

// Diagnostics is a base implementation of [providerops.Diagnostic]
// intended to be embedded into implementations outside of this module.
//
// Implementations embedding this type must still provide their own
// implementations of [providerops.Diagnostic.Severity],
// [providerops.Diagnostic.Summary], and
// [providerops.Diagnostic.Detail], but any new methods added to the interface
// type in future will have default implementations added to this base
// type to ensure forward-compatibility.
type Diagnostic struct {
	sealedImpl
}

// Diagnostics is a base implementation of [providerops.FunctionError]
// intended to be embedded into implementations outside of this module.
//
// Implementations embedding this type must still provide their own
// implementation of [providerops.FunctionError.Text], but other methods
// of the interface have default implementations in this base type.
type FunctionError struct {
	sealedImpl
}

// ArgumentIndex implements [providerops.FunctionError.ArgumentIndex] by just
// always reporting that the error is not related to any single argument.
func (FunctionError) ArgumentIndex() (int, bool) {
	return 0, false
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

type diagnosticsExample struct{ Diagnostics }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerops.Diagnostics = diagnosticsExample{}

func (diagnosticsExample) All() iter.Seq[providerops.Diagnostic] { panic("unimplemented") }
func (diagnosticsExample) HasErrors() bool                       { panic("unimplemented") }

type diagnosticExample struct{ Diagnostic }

// If future additions to [providerops.Diagnostic] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerops.Diagnostic = diagnosticExample{}

func (diagnosticExample) Severity() providerops.DiagnosticSeverity { panic("unimplemented") }
func (diagnosticExample) Summary() string                          { panic("unimplemented") }
func (diagnosticExample) Detail() string                           { panic("unimplemented") }

type functionErrorExample struct{ FunctionError }

// If future additions to [providerops.Diagnostic] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerops.FunctionError = functionErrorExample{}

func (functionErrorExample) Text() string { panic("unimplemented") }
