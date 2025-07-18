package providerschema

import (
	"iter"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"

	// For links in documentation comments:
	_ "slices"
)

// FunctionSignature describes the signature of a single function that a
// provider offers for use in expressions in OpenTofu configuration.
type FunctionSignature interface {
	// Parameters returns an iterable sequence of all of the non-variadic
	// parameters of the function, in the order they must appear in a
	// function call.
	//
	// Use [slices.Collect] with the result if you need a slice of parameters.
	Parameters() iter.Seq[FunctionParameter]

	// VariadicParameter returns a description of a function's optional
	// final variadic parameter, which accepts a caller-decided number
	// of additional arguments after the arguments corresponding to
	// the main parameters described by [FunctionSignature.Parameters].
	//
	// The result is nil for a function that does not accept any variadic
	// arguments at all.
	VariadicParameter() FunctionParameter

	// DocSummary returns the provider's human-readable short summary of the
	// purpose of this function, intended for use in a listing or table of
	// functions available for the provider, or an empty string if the provider
	// does not provide a summary.
	DocSummary() string

	// DocDescription returns the provider's human-readable long description
	// of the function, intended for use in a dedicated documentation page
	// about the function, or an empty string if the provider does not provide
	// a summary. The second result describes the intended format for the
	// the description string.
	DocDescription() (string, DocStringFormat)

	// DeprecationMessage returns a string containing a description of the
	// deprecation of this function, or an empty string if the function is
	// not deprecated.
	DeprecationMessage() string

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// FunctionParameter describes a single parameter as part of [FunctionSignature].
type FunctionParameter interface {
	// Name returns a name for this parameter which is intended primarily for
	// describing the parameter to a human author, such as in an error message
	// reporting that the argument was invalid.
	//
	// It should typically be a string that would be a valid HCL identifier,
	// but there is nothing stopping a provider from returning something
	// that does not meet that constraint.
	Name() string

	// Type returns the type constraint for values passed for this parameter.
	Type() TypeConstraint

	// NullValueAllowed returns true if callers are allowed to pass a null
	// value for this parameter, or false if not.
	//
	// Callers are expected to avoid passing null values to any parameter
	// for which this flag is unset; violating that constraint causes unspecified
	// behavior due to violating the provider's assumptions.
	//
	// This null value constraint applies only to the direct value passed
	// to the parameter. It does not constrain nested values, such as
	// passing null elements in a non-null list.
	NullValueAllowed() bool

	// UnknownValuesAllowed returns true if callers are allowed to pass
	// unknown values to this parameter, or false if not.
	//
	// Callers are expected to avoid passing unknown values in any part of
	// a parameter for which this flag is unset, including nested unknown
	// values inside known collections or structural values. Violating this
	// constraint causes unspecified behavior due to violating the provider's
	// assumptions.
	UnknownValuesAllowed() bool

	// DocDescription returns the provider's human-readable description
	// of the parameter. The second result describes the intended format for the
	// the description string.
	DocDescription() (string, DocStringFormat)

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
