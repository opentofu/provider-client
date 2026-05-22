package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// ValidateProviderConfigResponse is a base implementation of
// [providerops.ValidateProviderConfigResponse].
type ValidateProviderConfigResponse struct {
	sealedImpl
}

var _ providerops.ValidateProviderConfigResponse = ValidateProviderConfigResponse{}

// Diagnostics implements [providerops.ValidateProviderConfigResponse] by
// returning no diagnostics at all.
func (v ValidateProviderConfigResponse) Diagnostics() providerops.Diagnostics {
	return emptyDiagnostics{}
}
