package providerops

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type GetResourceIdentitySchemasRequest struct {
}

type GetResourceIdentitySchemasResponse interface {
	IdentitySchemas() iter.Seq2[string, providerschema.IdentitySchema]
	Diagnostics() Diagnostics

	common.Sealed
}
