package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// GetProviderSchemaResponse is a base implementation of
// [providerops.GetProviderSchemaResponse].
type GetProviderSchemaResponse struct {
	sealedImpl
}

// ServerCapabilities implements [providerops.GetProviderSchemaResponse] by
// just returning the zero value of [ServerCapabilities].
//
// Relying on this implementation implies a default set of capabilities that
// reflects what providers were expected to support at the time when this
// library was originally created, in May 2026. Review the documentation of
// [ServerCapabilities] to learn what those defaults are and decide whether you
// ought to override this.
func (g GetProviderSchemaResponse) ServerCapabilities() providerops.ServerCapabilities {
	return ServerCapabilities{}
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

type getProviderSchemaResponseExample struct{ GetProviderSchemaResponse }

// If future additions to [providerops.GetProviderSchemaResponse] cause a
// compile-time error here, refer to the comment labeled
// "EXAMPLE IMPLEMENTATIONS" above.
var _ providerops.GetProviderSchemaResponse = getProviderSchemaResponseExample{}

func (getProviderSchemaResponseExample) Diagnostics() providerops.Diagnostics { panic("unimplemented") }
func (getProviderSchemaResponseExample) ProviderSchema() providerschema.ProviderSchema {
	panic("unimplemented")
}
