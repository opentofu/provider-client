package tf5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"google.golang.org/grpc"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
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
	protoReq := &tfplugin5.ListResource_Request{
		TypeName:              req.TypeName,
		Config:                configVal,
		IncludeResourceObject: req.IncludeResourceObject,
		Limit:                 req.Limit,
	}

	protoResp, err := p.client.ListResource(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return listResourceResponse{proto: protoResp}, nil
}

type listResourceResponse struct {
	proto grpc.ServerStreamingClient[tfplugin5.ListResource_Event]

	common.SealedImpl
}

func (r listResourceResponse) Resources() iter.Seq2[providerops.ListResourceEvent, error] {
	return func(yield func(providerops.ListResourceEvent, error) bool) {
		for {
			res, err := r.proto.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				yield(nil, err)
				return
			}

			item := listResourceEvent{
				diagnostics: diagnostics{proto: res.GetDiagnostic()},
				displayName: res.GetDisplayName(),
			}

			// TODO Fetch identity data

			if res := res.GetResourceObject(); res != nil {
				item.resource = dynamicValue{proto: res}
			}

			if !yield(item, nil) {
				return
			}
		}
	}
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
