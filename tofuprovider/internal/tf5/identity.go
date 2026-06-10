package tf5

import (
	"context"
	"iter"
	"maps"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func (p *Provider) GetResourceIdentitySchemas(ctx context.Context, _ *providerops.GetResourceIdentitySchemasRequest) (providerops.GetResourceIdentitySchemasResponse, error) {
	protoReq := &tfplugin5.GetResourceIdentitySchemas_Request{
		// There are currently no fields in providerops.GetResourceIdentitySchemas,
		// so nothing to populate here.
	}
	protoResp, err := p.client.GetResourceIdentitySchemas(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return getResourceIdentitySchemas{proto: protoResp}, nil
}

type getResourceIdentitySchemas struct {
	proto *tfplugin5.GetResourceIdentitySchemas_Response
	common.SealedImpl
}

func (g getResourceIdentitySchemas) IdentitySchemas() iter.Seq2[string, providerschema.IdentitySchema] {
	return common.MapSeq2(maps.All(g.proto.IdentitySchemas), func(name string, protoSchema *tfplugin5.ResourceIdentitySchema) (string, providerschema.IdentitySchema) {
		return name, identitySchema{proto: protoSchema}
	})
}

func (g getResourceIdentitySchemas) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: g.proto.Diagnostics}
}

type identitySchema struct {
	proto *tfplugin5.ResourceIdentitySchema
	common.SealedImpl
}

func (i identitySchema) Version() int64 {
	return i.proto.Version
}

func (i identitySchema) IdentityAttributes() iter.Seq[providerschema.IdentityAttribute] {
	return func(yield func(providerschema.IdentityAttribute) bool) {
		for _, attr := range i.proto.IdentityAttributes {
			if !yield(identityAttribute{proto: attr}) {
				return
			}
		}
	}
}

type identityAttribute struct {
	proto *tfplugin5.ResourceIdentitySchema_IdentityAttribute
	common.SealedImpl
}

func (i identityAttribute) Name() string {
	return i.proto.GetName()
}

func (i identityAttribute) Type() providerschema.TypeConstraint {
	if len(i.proto.Type) == 0 {
		return nil
	}
	return common.CtyTypeJSON(i.proto.Type)
}

func (i identityAttribute) RequiredForImport() bool {
	return i.proto.GetRequiredForImport()
}

func (i identityAttribute) OptionalForImport() bool {
	return i.proto.GetOptionalForImport()
}

func (i identityAttribute) Description() string {
	return i.proto.GetDescription()
}

type resourceIdentityData struct {
	proto *tfplugin5.ResourceIdentityData
	common.SealedImpl
}

func (r resourceIdentityData) IdentityData() providerschema.DynamicValueOut {
	return dynamicValue{proto: r.proto.GetIdentityData()}
}
