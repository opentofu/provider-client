package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// ConfigureProviderResponse is a base implementation of
// [providerops.ConfigureProviderResponse].
type ConfigureProviderResponse struct {
	sealedImpl
}

var _ providerops.ConfigureProviderResponse = ConfigureProviderResponse{}

// Diagnostics implements [providerops.ConfigureProviderResponse] by returning
// no diagnostics at all.
func (c ConfigureProviderResponse) Diagnostics() providerops.Diagnostics {
	return emptyDiagnostics{}
}
