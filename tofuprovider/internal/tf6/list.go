package tf6

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// ListResource implements tofuprovider.ListResource.
func (p *Provider) ListResource(ctx context.Context, req *providerops.ListResourceRequest) (providerops.ListResourceResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin6.ListResource_Request{
		TypeName:              req.TypeName,
		Config:                configVal,
		IncludeResourceObject: req.IncludeResourceObject,
		Limit:                 req.Limit,
	}

	// Use a cancelable context so Close can terminate the stream even if the
	// caller stops reading before reaching the end.
	streamCtx, cancel := context.WithCancel(ctx)
	protoResp, err := p.client.ListResource(streamCtx, protoReq)
	if err != nil {
		cancel()
		return nil, err
	}

	return &listResourceResponse{proto: protoResp, cancel: cancel}, nil
}

type listResourceResponse struct {
	proto  grpc.ServerStreamingClient[tfplugin6.ListResource_Event]
	cancel context.CancelFunc

	common.SealedImpl
}

// ReadResult ignores its context because the gRPC stream is already bound to
// the context passed to ListResource; the parameter exists only for
// non-streaming implementations of the interface.
func (r *listResourceResponse) ReadResult(_ context.Context) (providerops.ListResourceEvent, error) {
	// Recv returns io.EOF at the end of the stream, which we pass through to
	// the caller as the loop-termination signal.
	res, err := r.proto.Recv()
	if err != nil {
		return nil, err
	}

	item := listResourceEvent{
		diagnostics: diagnostics{proto: res.GetDiagnostic()},
		displayName: res.GetDisplayName(),
	}

	// TODO Fetch identity data

	if res := res.GetResourceObject(); res != nil {
		item.resource = dynamicValue{proto: res}
	}

	return item, nil
}

// Close ignores its context and never errors because terminating a gRPC
// stream is just a local context cancellation; the signature exists for
// non-streaming implementations that may do fallible cleanup.
func (r *listResourceResponse) Close(_ context.Context) error {
	r.cancel()
	return nil
}

type listResourceEvent struct {
	displayName string
	resource    providerschema.DynamicValueOut
	diagnostics providerops.Diagnostics

	common.SealedImpl
}

func (i listResourceEvent) DisplayName() string                      { return i.displayName }
func (i listResourceEvent) Resource() providerschema.DynamicValueOut { return i.resource }
func (i listResourceEvent) Diagnostics() providerops.Diagnostics     { return i.diagnostics }
