package providerops

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"

	// For links in documentation comments:
	_ "slices"
)

type Diagnostics interface {
	// HasErrors returns true if any of the diagnostics have severity
	// [DiagnosticError], which often suggests that other parts of a response
	// are invalid or incomplete and so a caller ought to return early or
	// otherwise ignore the results, even if it doesn't make direct use of
	// the individual diagnostics.
	//
	// (The documentation for some operations may make exceptions about what
	// caller behavior is appropriate when a response includes error
	// diagnostics. The above is just broad guidance for the common case.)
	HasErrors() bool

	// All returns an iterable sequence of each of the individual diagnostics.
	//
	// Use [slices.Collect] with the result to gather all of the diagnostics
	// into a slice, if needed.
	All() iter.Seq[Diagnostic]

	common.Sealed
}

type Diagnostic interface {
	Severity() DiagnosticSeverity
	Summary() string
	Detail() string

	// TODO: AttributePath, which has an awkward model as a sequence of
	// values of three possible types.

	common.Sealed
}

type DiagnosticSeverity int

const (
	DiagnosticUnsupported DiagnosticSeverity = 0
	DiagnosticWarning     DiagnosticSeverity = 1
	DiagnosticError       DiagnosticSeverity = 2
)

// FunctionError describes an error that occurred when calling a function.
type FunctionError interface {
	// Text returns the human-oriented description of the problem.
	Text() string

	// ArgumentIndex returns the index of an argument that is being blamed
	// for the problem along with true, or a meaningless value along with
	// false if this error is not specific to an individual argument.
	ArgumentIndex() (int, bool)
}
