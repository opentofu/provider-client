package tf6

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

// CallFunction implements tofuprovider.GRPCPluginProvider.
func (p *Provider) CallFunction(ctx context.Context, req *providerops.CallFunctionRequest) (providerops.CallFunctionResponse, error) {
	args, err := prepareFunctionArgs(req.Arguments)
	if err != nil {
		return callFunctionResponse{
			proto: &tfplugin6.CallFunction_Response{
				Error: &tfplugin6.FunctionError{
					Text: err.Error(),
				},
			},
		}, nil
	}

	resp, err := p.client.CallFunction(ctx, &tfplugin6.CallFunction_Request{
		Name:      req.FunctionName,
		Arguments: args,
	})
	return callFunctionResponse{proto: resp}, err
}

func prepareFunctionArgs(args []providerschema.DynamicValueIn) ([]*tfplugin6.DynamicValue, error) {
	if len(args) == 0 {
		return nil, nil
	}
	ret := make([]*tfplugin6.DynamicValue, len(args))
	for i, arg := range args {
		src, err := common.CtyValueAsMsgpack(arg.Value(), arg.SerializationType())
		if err != nil {
			// This indicates a bug in our caller, rather than a problem caused
			// by our caller's end-user input, so we accept a relatively
			// low-quality error message here.
			return nil, fmt.Errorf("invalid value for argument %d: %w", i, err)
		}
		ret[i] = &tfplugin6.DynamicValue{
			Msgpack: src,
		}
	}
	return ret, nil
}

type callFunctionResponse struct {
	proto *tfplugin6.CallFunction_Response
	common.Sealed
}

// Error implements providerops.CallFunctionResponse.
func (c callFunctionResponse) Error() providerops.FunctionError {
	if c.proto.Error == nil {
		return nil
	}
	return functionError{proto: c.proto.Error}
}

// Result implements providerops.CallFunctionResponse.
func (c callFunctionResponse) Result() providerschema.DynamicValueOut {
	if c.proto.Result == nil {
		return nil
	}
	return dynamicValue{proto: c.proto.Result}
}

type functionError struct {
	proto *tfplugin6.FunctionError
}

// ArgumentIndex implements providerops.FunctionError.
func (f functionError) ArgumentIndex() (int, bool) {
	if f.proto.FunctionArgument == nil {
		return 0, false
	}
	// This conversion to int is okay because in practice there can't be
	// more arguments than int can store because Go uses int to represent
	// the length of a slice. Anything out of bounds here cannot possibly
	// be valid.
	return int(*f.proto.FunctionArgument), true
}

// Text implements providerops.FunctionError.
func (f functionError) Text() string {
	return f.proto.Text
}
