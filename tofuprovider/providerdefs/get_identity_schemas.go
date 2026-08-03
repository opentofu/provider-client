package providerdefs

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// GetIdentitySchemasResponse is a base implementation of
// [providerops.GetIdentitySchemasResponse].
type GetIdentitySchemasResponse struct {
	sealedImpl
}

var _ providerops.GetIdentitySchemasResponse = GetIdentitySchemasResponse{}

// Diagnostics implements [providerops.GetIdentitySchemasResponse] by returning
// no diagnostics at all.
func (g GetIdentitySchemasResponse) Diagnostics() providerops.Diagnostics {
	return emptyDiagnostics{}
}

// IdentitySchemas implements [providerops.GetIdentitySchemasResponse] by
// reporting that no managed resource types have identity schemas.
func (g GetIdentitySchemasResponse) IdentitySchemas() iter.Seq2[string, providerschema.IdentitySchema] {
	return func(func(string, providerschema.IdentitySchema) bool) {}
}
