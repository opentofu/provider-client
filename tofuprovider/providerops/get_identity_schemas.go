package providerops

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type GetIdentitySchemasRequest struct {
	// There are currently no arguments in a resource identity schemas request.
}

type GetIdentitySchemasResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// IdentitySchemas returns an iterable sequence of identity schemas,
	// one for each managed resource type whose identity the provider
	// is able to report.
	//
	// The first result of each item is the unique managed resource type
	// name. Use [maps.Collect] to produce a map from resource type name
	// to identity schema.
	IdentitySchemas() iter.Seq2[string, providerschema.IdentitySchema]

	common.Sealed
}
