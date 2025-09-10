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

func (p *Provider) GetProviderSchema(ctx context.Context, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error) {
	protoReq := &tfplugin5.GetProviderSchema_Request{
		// There are currently no fields in providerops.GetProviderSchemaRequest,
		// so nothing to populate here.
	}
	protoResp, err := p.client.GetSchema(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return getProviderSchemaResponse{proto: protoResp}, nil
}

type getProviderSchemaResponse struct {
	proto *tfplugin5.GetProviderSchema_Response

	common.SealedImpl
}

// Diagnostics implements providerops.GetProviderSchemaResponse.
func (g getProviderSchemaResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: g.proto.Diagnostics}
}

// ProviderSchema implements providerops.GetProviderSchemaResponse.
func (g getProviderSchemaResponse) ProviderSchema() providerschema.ProviderSchema {
	return providerSchema{proto: g.proto}
}

// ServerCapabilities implements providerops.GetProviderSchemaResponse.
func (g getProviderSchemaResponse) ServerCapabilities() providerops.ServerCapabilities {
	return serverCapabilities{proto: g.proto.ServerCapabilities}
}

func (p *Provider) GetFunctions(ctx context.Context, req *providerops.GetFunctionsRequest) (providerops.GetFunctionsResponse, error) {
	panic("unimplemented")
}

type providerSchema struct {
	proto *tfplugin5.GetProviderSchema_Response

	common.SealedImpl
}

// DataResourceTypeSchemas implements providerschema.ProviderSchema.
func (p providerSchema) DataResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return namedSchemasSeq(p.proto.DataSourceSchemas)
}

// EphemeralResourceTypeSchemas implements providerschema.ProviderSchema.
func (p providerSchema) EphemeralResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return namedSchemasSeq(p.proto.EphemeralResourceSchemas)
}

// FunctionSignatures implements providerschema.ProviderSchema.
func (p providerSchema) FunctionSignatures() iter.Seq2[string, providerschema.FunctionSignature] {
	return namedFunctionsSeq(p.proto.Functions)
}

// ManagedResourceTypeSchemas implements providerschema.ProviderSchema.
func (p providerSchema) ManagedResourceTypeSchemas() iter.Seq2[string, providerschema.Schema] {
	return namedSchemasSeq(p.proto.ResourceSchemas)
}

// ProviderConfigSchema implements providerschema.ProviderSchema.
func (p providerSchema) ProviderConfigSchema() providerschema.Schema {
	if p.proto.Provider == nil {
		return nil
	}
	return schema{proto: p.proto.Provider}
}

// ProviderMetaSchema implements providerschema.ProviderSchema.
func (p providerSchema) ProviderMetaSchema() providerschema.Schema {
	if p.proto.ProviderMeta == nil {
		return nil
	}
	return schema{proto: p.proto.ProviderMeta}
}

type schema struct {
	proto *tfplugin5.Schema
	common.SealedImpl
}

func namedSchemasSeq(proto map[string]*tfplugin5.Schema) iter.Seq2[string, providerschema.Schema] {
	return common.MapSeq2(maps.All(proto), func(name string, protoSchema *tfplugin5.Schema) (string, providerschema.Schema) {
		return name, schema{proto: protoSchema}
	})
}

// Attributes implements providerschema.Schema.
func (s schema) Attributes() iter.Seq2[string, providerschema.Attribute] {
	return attributesSeq(s.proto.Block.Attributes)
}

// NestedBlockTypes implements providerschema.Schema.
func (s schema) NestedBlockTypes() iter.Seq2[string, providerschema.NestedBlockType] {
	return nestedBlockTypesSeq(s.proto.Block.BlockTypes)
}

// SchemaVersion implements providerschema.Schema.
func (s schema) SchemaVersion() int64 {
	return s.proto.Version
}

// DocDescription implements providerschema.Schema.
func (s schema) DocDescription() (string, providerschema.DocStringFormat) {
	if s.proto.Block == nil {
		return "", providerschema.DocStringPlain
	}
	return s.proto.Block.Description, docStringFormat(s.proto.Block.DescriptionKind)
}

type attribute struct {
	proto *tfplugin5.Schema_Attribute
	common.SealedImpl
}

func attributesSeq(proto []*tfplugin5.Schema_Attribute) iter.Seq2[string, providerschema.Attribute] {
	return common.MapSeqToSeq2(slices.Values(proto), func(protoAttr *tfplugin5.Schema_Attribute) (string, providerschema.Attribute) {
		return protoAttr.Name, attribute{proto: protoAttr}
	})
}

// NestedType implements providerschema.Attribute.
func (a attribute) NestedType() providerschema.ObjectType {
	// Protocol 5 does not include the concept of "structural-typed" attributes
	// at all, so this is always nil.
	return nil
}

// Type implements providerschema.Attribute.
func (a attribute) Type() providerschema.TypeConstraint {
	if len(a.proto.Type) == 0 {
		return nil
	}
	return common.CtyTypeJSON(a.proto.Type)
}

// DeprecationMessage implements providerschema.Attribute.
func (a attribute) IsDeprecated() bool {
	return a.proto.Deprecated
}

// DocDescription implements providerschema.Attribute.
func (a attribute) DocDescription() (string, providerschema.DocStringFormat) {
	return a.proto.Description, docStringFormat(a.proto.DescriptionKind)
}

// IsSensitive implements providerschema.Attribute.
func (a attribute) IsSensitive() bool {
	return a.proto.Sensitive
}

// IsWriteOnly implements providerschema.Attribute.
func (a attribute) IsWriteOnly() bool {
	return a.proto.WriteOnly
}

// Usage implements providerschema.Attribute.
func (a attribute) Usage() providerschema.AttributeUsage {
	switch {
	case a.proto.Required && !a.proto.Optional && !a.proto.Computed:
		return providerschema.AttributeRequired
	case !a.proto.Required && a.proto.Optional && !a.proto.Computed:
		return providerschema.AttributeOptional
	case !a.proto.Required && a.proto.Optional && a.proto.Computed:
		return providerschema.AttributeOptionalComputed
	case !a.proto.Required && !a.proto.Optional && a.proto.Computed:
		return providerschema.AttributeComputed
	default:
		return providerschema.AttributeUsageUnsupported
	}
}

type nestedBlockType struct {
	proto *tfplugin5.Schema_NestedBlock
	common.SealedImpl
}

func nestedBlockTypesSeq(proto []*tfplugin5.Schema_NestedBlock) iter.Seq2[string, providerschema.NestedBlockType] {
	return common.MapSeqToSeq2(slices.Values(proto), func(protoBlock *tfplugin5.Schema_NestedBlock) (string, providerschema.NestedBlockType) {
		return protoBlock.TypeName, nestedBlockType{proto: protoBlock}
	})
}

// Attributes implements providerschema.NestedBlockType.
func (n nestedBlockType) Attributes() iter.Seq2[string, providerschema.Attribute] {
	return attributesSeq(n.proto.Block.Attributes)
}

// ItemLimits implements providerschema.NestedBlockType.
func (n nestedBlockType) ItemLimits() (int64, int64) {
	return n.proto.MinItems, n.proto.MaxItems
}

// NestedBlockTypes implements providerschema.NestedBlockType.
func (n nestedBlockType) NestedBlockTypes() iter.Seq2[string, providerschema.NestedBlockType] {
	return nestedBlockTypesSeq(n.proto.Block.BlockTypes)
}

// Nesting implements providerschema.NestedBlockType.
func (n nestedBlockType) Nesting() providerschema.NestingMode {
	return blockNestingMode(n.proto.Nesting)
}

type functionSignature struct {
	proto *tfplugin5.Function
	common.SealedImpl
}

func namedFunctionsSeq(proto map[string]*tfplugin5.Function) iter.Seq2[string, providerschema.FunctionSignature] {
	return common.MapSeq2(maps.All(proto), func(name string, protoFunc *tfplugin5.Function) (string, providerschema.FunctionSignature) {
		return name, functionSignature{proto: protoFunc}
	})
}

// DeprecationMessage implements providerschema.FunctionSignature.
func (f functionSignature) DeprecationMessage() string {
	return f.proto.DeprecationMessage
}

// DocDescription implements providerschema.FunctionSignature.
func (f functionSignature) DocDescription() (string, providerschema.DocStringFormat) {
	return f.proto.Description, docStringFormat(f.proto.DescriptionKind)
}

// DocSummary implements providerschema.FunctionSignature.
func (f functionSignature) DocSummary() string {
	return f.proto.Summary
}

// Parameters implements providerschema.FunctionSignature.
func (f functionSignature) Parameters() iter.Seq[providerschema.FunctionParameter] {
	return functionParametersSeq(f.proto.Parameters)
}

// VariadicParameter implements providerschema.FunctionSignature.
func (f functionSignature) VariadicParameter() providerschema.FunctionParameter {
	if f.proto.VariadicParameter == nil {
		return nil
	}
	return functionParameter{proto: f.proto.VariadicParameter}
}

// ResultType implements providerschema.FunctionSignature.
func (f functionSignature) ResultType() providerschema.TypeConstraint {
	return common.CtyTypeJSON(f.proto.Return.Type)
}

type functionParameter struct {
	proto *tfplugin5.Function_Parameter
	common.SealedImpl
}

func functionParametersSeq(proto []*tfplugin5.Function_Parameter) iter.Seq[providerschema.FunctionParameter] {
	return common.MapSeq(slices.Values(proto), func(protoParam *tfplugin5.Function_Parameter) providerschema.FunctionParameter {
		return functionParameter{proto: protoParam}
	})
}

// DocDescription implements providerschema.FunctionParameter.
func (f functionParameter) DocDescription() (string, providerschema.DocStringFormat) {
	return f.proto.Description, docStringFormat(f.proto.DescriptionKind)
}

// Name implements providerschema.FunctionParameter.
func (f functionParameter) Name() string {
	return f.proto.Name
}

// NullValueAllowed implements providerschema.FunctionParameter.
func (f functionParameter) NullValueAllowed() bool {
	return f.proto.AllowNullValue
}

// Type implements providerschema.FunctionParameter.
func (f functionParameter) Type() providerschema.TypeConstraint {
	return common.CtyTypeJSON(f.proto.Type)
}

// UnknownValuesAllowed implements providerschema.FunctionParameter.
func (f functionParameter) UnknownValuesAllowed() bool {
	return f.proto.AllowUnknownValues
}

func blockNestingMode(proto tfplugin5.Schema_NestedBlock_NestingMode) providerschema.NestingMode {
	switch proto {
	case tfplugin5.Schema_NestedBlock_SINGLE:
		return providerschema.NestingSingle
	case tfplugin5.Schema_NestedBlock_GROUP:
		return providerschema.NestingGroup
	case tfplugin5.Schema_NestedBlock_LIST:
		return providerschema.NestingList
	case tfplugin5.Schema_NestedBlock_SET:
		return providerschema.NestingSet
	case tfplugin5.Schema_NestedBlock_MAP:
		return providerschema.NestingMap
	default:
		return providerschema.NestingInvalid
	}
}

func docStringFormat(proto tfplugin5.StringKind) providerschema.DocStringFormat {
	switch proto {
	case tfplugin5.StringKind_PLAIN:
		return providerschema.DocStringPlain
	case tfplugin5.StringKind_MARKDOWN:
		return providerschema.DocStringMarkdown
	default:
		return providerschema.DocStringUnsupported
	}
}
