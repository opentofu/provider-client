package providerops

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type GetFunctionsRequest struct {
	// There are currently no arguments in a functions request.
}

type GetFunctionsResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// FunctionSignatures returns an iterable sequence of the signature of
	// each "provider-defined function" supported by this provider.
	//
	// The first result in each pair is the unique function type name that
	// the signature belongs to. Use [maps.Collect] to gather the result into a
	// map from name to signature if you expect to need signatures for more than
	// one function.
	FunctionSignatures() iter.Seq2[string, providerschema.FunctionSignature]

	common.Sealed
}
