package mockutil

import (
	"testing"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// ComparableDiagnostic is an implementation of [providerops.Diagnostic] that
// is friendly to being compared using functions from [cmp].
//
// If you use [Eq] (or other similar [gomock.Matcher] implementations based on
// [Comparer]) or compare using [Diff], then any implementations of
// [providerops.Diagnostic] and [providerops.Diagnostic] are automatically
// converted to this type so that diagnostics are compared by their user-facing
// messages rather than by the types used to implement them.
//
// This type is exported as a convenient way to describe "expected diagnostics"
// in a test, using composite literal syntax.
type ComparableDiagnostic struct {
	Severity_ providerops.DiagnosticSeverity
	Summary_  string
	Detail_   string

	common.SealedImpl
}

var _ providerops.Diagnostic = (*ComparableDiagnostic)(nil)

// Severity implements [providerops.Diagnostic].
func (diag *ComparableDiagnostic) Severity() providerops.DiagnosticSeverity {
	return diag.Severity_
}

// Summary implements [providerops.Diagnostic].
func (diag *ComparableDiagnostic) Summary() string {
	return diag.Summary_
}

// Detail implements [providerops.Diagnostic].
func (diag *ComparableDiagnostic) Detail() string {
	return diag.Detail_
}

func cmpTransformDiagnostic(diag providerops.Diagnostic) *ComparableDiagnostic {
	return &ComparableDiagnostic{
		Severity_: diag.Severity(),
		Summary_:  diag.Summary(),
		Detail_:   diag.Detail(),

		// TODO: AttributePath, once [providerops.Diagnostic] has that method.
	}
}

func cmpTransformDiagnostics(diags providerops.Diagnostics) []*ComparableDiagnostic {
	var ret []*ComparableDiagnostic
	for diag := range diags.All() {
		ret = append(ret, cmpTransformDiagnostic(diag))
	}
	return ret
}

// AssertNoDiags is a test helper that immediately fails the test (halting
// further execution) if the given [providerops.Diagnostics] has any diagnostics
// in it, regardless of severity.
func AssertNoDiags(t *testing.T, diags providerops.Diagnostics) {
	t.Helper()
	foundDiags := false
	for diag := range diags.All() {
		t.Errorf("unexpected %s diagnostic\nsummary: %s\ndetail:  %s", diag.Severity(), diag.Summary(), diag.Detail())
		foundDiags = true
	}
	if foundDiags {
		t.FailNow()
	}
}
