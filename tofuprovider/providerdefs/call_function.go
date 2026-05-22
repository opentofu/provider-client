package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// CallFunctionResponse is a base implementation of
// [providerops.CallFunctionResponse].
//
// External implementations should embed this for forward-compatibility with
// future extensions to the interface, but must still provide their own
// implementations of [providerops.CallFunctionResponse.Result] and
// [providerops.CallFunctionResponse.Error]
type CallFunctionResponse struct {
	sealedImpl
}

// NOTE: Right now there aren't actually any other methods of
// CallFunctionResponse to implement here, but external implementers are
// expected to embed this type anyway so that we can potentially add new methods
// to the interface in future and then write their default implementations here.
// If you're here reading this comment because that's what you're about to do
// then you can remove this comment at the same time!

// ===================== EXAMPLE IMPLEMENTATIONS ==============================
//
// The remaining unexported types are here just to make sure that a third-party
// implementation that already has these methods can be made to implement
// the corresponding interfaces just by embedding our default types from above.
//
// If future additions to the interface types cause failures here then the
// appropriate fix is to add default implementations of those methods to the
// exported types, _not_ to add those methods to these unexported types.

type callFunctionResponseExample struct{ CallFunctionResponse }

// If future additions to [providerops.Diagnostics] cause a compile-time error
// here, refer to the comment labeled "EXAMPLE IMPLEMENTATIONS" above.
var _ providerops.CallFunctionResponse = callFunctionResponseExample{}

func (callFunctionResponseExample) Error() providerops.FunctionError       { panic("unimplemented") }
func (callFunctionResponseExample) Result() providerschema.DynamicValueOut { panic("unimplemented") }
