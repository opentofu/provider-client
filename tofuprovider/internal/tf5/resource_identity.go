package tf5

import (
	"context"
	"iter"
	"maps"
	"slices"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func (p *Provider) GetIdentitySchemas(ctx context.Context, req *providerops.GetIdentitySchemasRequest) (providerops.GetIdentitySchemasResponse, error) {
	protoReq := &tfplugin5.GetResourceIdentitySchemas_Request{
		// There are currently no fields in providerops.GetIdentitySchemasRequest,
		// so nothing to populate here.
	}
	protoResp, err := p.client.GetResourceIdentitySchemas(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return getIdentitySchemasResponse{proto: protoResp}, nil
}

type getIdentitySchemasResponse struct {
	proto *tfplugin5.GetResourceIdentitySchemas_Response

	common.SealedImpl
}

var _ providerops.GetIdentitySchemasResponse = getIdentitySchemasResponse{}

func (g getIdentitySchemasResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: g.proto.Diagnostics}
}

func (g getIdentitySchemasResponse) IdentitySchemas() iter.Seq2[string, providerschema.IdentitySchema] {
	return namedIdentitySchemasSeq(g.proto.IdentitySchemas)
}

type identitySchema struct {
	proto *tfplugin5.ResourceIdentitySchema

	common.SealedImpl
}

var _ providerschema.IdentitySchema = identitySchema{}

func namedIdentitySchemasSeq(proto map[string]*tfplugin5.ResourceIdentitySchema) iter.Seq2[string, providerschema.IdentitySchema] {
	return common.MapSeq2(maps.All(proto), func(name string, protoSchema *tfplugin5.ResourceIdentitySchema) (string, providerschema.IdentitySchema) {
		return name, identitySchema{proto: protoSchema}
	})
}

func (s identitySchema) SchemaVersion() int64 {
	return s.proto.Version
}

func (s identitySchema) Attributes() iter.Seq2[string, providerschema.IdentityAttribute] {
	return identityAttributesSeq(s.proto.IdentityAttributes)
}

type identityAttribute struct {
	proto *tfplugin5.ResourceIdentitySchema_IdentityAttribute

	common.SealedImpl
}

var _ providerschema.IdentityAttribute = identityAttribute{}

func identityAttributesSeq(proto []*tfplugin5.ResourceIdentitySchema_IdentityAttribute) iter.Seq2[string, providerschema.IdentityAttribute] {
	return common.MapSeqToSeq2(slices.Values(proto), func(protoAttr *tfplugin5.ResourceIdentitySchema_IdentityAttribute) (string, providerschema.IdentityAttribute) {
		return protoAttr.Name, identityAttribute{proto: protoAttr}
	})
}

func (a identityAttribute) Type() providerschema.TypeConstraint {
	if len(a.proto.Type) == 0 {
		return nil
	}
	return common.CtyTypeJSON(a.proto.Type)
}

func (a identityAttribute) DocDescription() (string, providerschema.DocStringFormat) {
	return a.proto.Description, providerschema.DocStringMarkdown
}

func (a identityAttribute) ImportUsage() providerschema.AttributeUsage {
	switch {
	case a.proto.RequiredForImport && a.proto.OptionalForImport:
		return providerschema.AttributeUsageUnsupported
	case a.proto.RequiredForImport:
		return providerschema.AttributeRequired
	case a.proto.OptionalForImport:
		return providerschema.AttributeOptional
	default:
		return providerschema.AttributeUsageUnsupported
	}
}

func (p *Provider) UpgradeIdentity(ctx context.Context, req *providerops.UpgradeIdentityRequest) (providerops.UpgradeIdentityResponse, error) {
	protoReq := &tfplugin5.UpgradeResourceIdentity_Request{
		TypeName: req.ResourceType,
		Version:  req.SchemaVersion,
		RawIdentity: &tfplugin5.RawState{
			Json: req.PrevIdentityJSON,
		},
	}
	protoResp, err := p.client.UpgradeResourceIdentity(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return upgradeIdentityResponse{proto: protoResp}, nil
}

type upgradeIdentityResponse struct {
	proto *tfplugin5.UpgradeResourceIdentity_Response

	common.SealedImpl
}

var _ providerops.UpgradeIdentityResponse = upgradeIdentityResponse{}

func (u upgradeIdentityResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: u.proto.Diagnostics}
}

func (u upgradeIdentityResponse) UpgradedIdentity() providerschema.DynamicValueOut {
	if u.proto.UpgradedIdentity == nil || u.proto.UpgradedIdentity.IdentityData == nil {
		return nil
	}
	return dynamicValue{proto: u.proto.UpgradedIdentity.IdentityData}
}
