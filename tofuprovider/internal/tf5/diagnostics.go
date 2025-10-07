package tf5

import (
	"iter"
	"slices"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

type diagnostics struct {
	proto []*tfplugin5.Diagnostic

	common.SealedImpl
}

// All implements providerops.Diagnostics.
func (d diagnostics) All() iter.Seq[providerops.Diagnostic] {
	return common.MapSeq(slices.Values(d.proto), func(proto *tfplugin5.Diagnostic) providerops.Diagnostic {
		return diagnostic{proto: proto}
	})
}

// HasErrors implements providerops.Diagnostics.
func (d diagnostics) HasErrors() bool {
	for _, diag := range d.proto {
		if diag.Severity == tfplugin5.Diagnostic_ERROR {
			return true
		}
	}
	return false
}

type diagnostic struct {
	proto *tfplugin5.Diagnostic

	common.SealedImpl
}

// Detail implements providerops.Diagnostic.
func (d diagnostic) Detail() string {
	return d.proto.Detail
}

// Severity implements providerops.Diagnostic.
func (d diagnostic) Severity() providerops.DiagnosticSeverity {
	switch d.proto.Severity {
	case tfplugin5.Diagnostic_ERROR:
		return providerops.DiagnosticError
	case tfplugin5.Diagnostic_WARNING:
		return providerops.DiagnosticWarning
	default:
		return providerops.DiagnosticUnsupported
	}
}

// Summary implements providerops.Diagnostic.
func (d diagnostic) Summary() string {
	return d.proto.Summary
}
