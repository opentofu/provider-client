package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type CallFunctionRequest struct {
	// FunctionName is the name of the function to call.
	FunctionName string

	// Arguments are the values to use as the function's arguments.
	//
	// The acceptable values depend on the function signature as described
	// by the corresponding [providerschema.FunctionSignature] object:
	//
	// - There must be at least as many elements as there are results
	//   from [providerschema.FunctionSignature.Parameters]. If
	//   [providerschema.FunctionSignature.VariadicParameter] returns
	//   nil then there must be _exactly_ as many elements as normal
	//   parameters.
	// - The serialization type of each parameter must exactly match
	//   what was returned by the [FunctionParameter.Type] method of
	//   the corresponding parameter.
	Arguments []providerschema.DynamicValueIn
}

type CallFunctionResponse interface {
	// Error returns nil if the function call was successful, or a non-nil
	// object describing why the function call failed.
	//
	// If this returns a non-nil result then the results from other methods
	// are meaningless.
	Error() FunctionError

	// Result returns the result of the function call.
	//
	// This must be decoded with the type that was returned from this
	// function's [providerschema.FunctionSignature.ResultType] method.
	Result() providerschema.DynamicValueOut

	common.Sealed
}
