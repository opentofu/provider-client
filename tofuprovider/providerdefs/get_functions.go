package providerdefs

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// GetFunctionsResponse is a base implementation of
// [providerops.GetFunctionsResponse].
type GetFunctionsResponse struct {
	sealedImpl
}

var _ providerops.GetFunctionsResponse = GetFunctionsResponse{}

// Diagnostics implements [providerops.GetFunctionsResponse] by returning
// no diagnostics at all.
func (g GetFunctionsResponse) Diagnostics() providerops.Diagnostics {
	return emptyDiagnostics{}
}

// FunctionSignatures implements [providerops.GetFunctionsResponse] by reporting
// that no functions exist.
func (g GetFunctionsResponse) FunctionSignatures() iter.Seq2[string, providerschema.FunctionSignature] {
	return func(func(string, providerschema.FunctionSignature) bool) {}
}
