package tf5

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// ListManagedResources implements tofuprovider.ListManagedResources.
func (p *Provider) ListManagedResources(ctx context.Context, req *providerops.ListManagedResourcesRequest) (providerops.ListManagedResourcesResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin5.ListResource_Request{
		TypeName:              req.TypeName,
		Config:                configVal,
		IncludeResourceObject: req.IncludeResourceObject,
		Limit:                 req.Limit,
	}

	streamCtx, cancel := context.WithCancel(ctx)
	protoResp, err := p.client.ListResource(streamCtx, protoReq)
	if err != nil {
		cancel()
		return nil, err
	}

	return &listManagedResourcesResponse{proto: protoResp, cancel: cancel}, nil
}

type listManagedResourcesResponse struct {
	proto  grpc.ServerStreamingClient[tfplugin5.ListResource_Event]
	cancel context.CancelFunc

	common.SealedImpl
}

func (r *listManagedResourcesResponse) ReadResult(_ context.Context) (providerops.ListManagedResourcesEvent, error) {
	res, err := r.proto.Recv()
	if err != nil {
		return nil, err
	}

	item := listManagedResourcesEvent{
		diagnostics: diagnostics{proto: res.GetDiagnostic()},
		displayName: res.GetDisplayName(),
	}

	// TODO Fetch identity data

	if res := res.GetResourceObject(); res != nil {
		item.resource = dynamicValue{proto: res}
	}

	return item, nil
}

func (r *listManagedResourcesResponse) Close(_ context.Context) error {
	r.cancel()
	return nil
}

type listManagedResourcesEvent struct {
	displayName string
	resource    providerschema.DynamicValueOut
	diagnostics providerops.Diagnostics

	common.SealedImpl
}

func (i listManagedResourcesEvent) DisplayName() string                      { return i.displayName }
func (i listManagedResourcesEvent) Resource() providerschema.DynamicValueOut { return i.resource }
func (i listManagedResourcesEvent) Diagnostics() providerops.Diagnostics     { return i.diagnostics }
