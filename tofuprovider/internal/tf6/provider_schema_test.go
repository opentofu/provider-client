package tf6_test

import (
	"context"
	"maps"
	"slices"
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf6"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func TestGetProviderSchema(t *testing.T) {
	testProviderCalls(t,
		map[string]providerCallTest[*providerops.GetProviderSchemaRequest, providerops.GetProviderSchemaResponse]{
			"empty schema": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.ContextWithValue("propagation_test", true), // set in the "call" function below
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					// Our main purpose here is just to make sure we can call
					// everything and get the expected "nothing" result instead
					// of crashing when the provider returns nothing at all.
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					schema := resp.ProviderSchema()
					if got := schema.ProviderConfigSchema(); got != nil {
						t.Errorf("unexpected provider config schema: %#v", got)
					}
					if got := schema.ProviderMetaSchema(); got != nil {
						t.Errorf("unexpected provider meta schema: %#v", got)
					}
					checkEmptySeq2(t, "managed resource type", schema.ManagedResourceTypeSchemas())
					checkEmptySeq2(t, "data resource type", schema.DataResourceTypeSchemas())
					checkEmptySeq2(t, "ephemeral resource type", schema.EphemeralResourceTypeSchemas())
					checkEmptySeq2(t, "function", schema.FunctionSignatures())
				},
			},
			"diagnostics": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							Diagnostics: []*tfplugin6.Diagnostic{
								{
									Severity: tfplugin6.Diagnostic_ERROR,
									Summary:  "Failed to do something",
									Detail:   "Couldn't do it.",
								},
								{
									Severity: tfplugin6.Diagnostic_WARNING,
									Summary:  "Doing things is deprecated",
									Detail:   "Consider just generating the most likely outcome based on a statistical model.",
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					gotDiags := slices.Collect(resp.Diagnostics().All())
					wantDiags := []providerops.Diagnostic{
						&mockutil.ComparableDiagnostic{
							Severity_: providerops.DiagnosticError,
							Summary_:  "Failed to do something",
							Detail_:   "Couldn't do it.",
						},
						&mockutil.ComparableDiagnostic{
							Severity_: providerops.DiagnosticWarning,
							Summary_:  "Doing things is deprecated",
							Detail_:   "Consider just generating the most likely outcome based on a statistical model.",
						},
					}
					if diff := mockutil.Diff(wantDiags, gotDiags); diff != "" {
						t.Error("wrong diagnostics\n" + diff)
					}
				},
			},
			"provider config schema": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							Provider: &tfplugin6.Schema{
								Block: &tfplugin6.Schema_Block{
									Attributes: []*tfplugin6.Schema_Attribute{
										{
											Name:     "url",
											Type:     mockutil.JSON("string"),
											Optional: true,
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					schema := resp.ProviderSchema().ProviderConfigSchema()
					if schema == nil {
						t.Fatal("no provider config schema was returned")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["url"]
					if !ok {
						t.Fatal("no attribute named 'url' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeOptional; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"provider meta schema": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							ProviderMeta: &tfplugin6.Schema{
								Block: &tfplugin6.Schema_Block{
									Attributes: []*tfplugin6.Schema_Attribute{
										{
											Name:     "module_addr",
											Type:     mockutil.JSON("string"),
											Optional: true,
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					schema := resp.ProviderSchema().ProviderMetaSchema()
					if schema == nil {
						t.Fatal("no provider meta schema was returned")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["module_addr"]
					if !ok {
						t.Fatal("no attribute named 'module_addr' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeOptional; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"managed resource types": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							ResourceSchemas: map[string]*tfplugin6.Schema{
								"foo": {
									Block: &tfplugin6.Schema_Block{
										Attributes: []*tfplugin6.Schema_Attribute{
											{
												Name:     "name",
												Type:     mockutil.JSON("string"),
												Required: true,
											},
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					got := maps.Collect(resp.ProviderSchema().ManagedResourceTypeSchemas())
					schema, ok := got["foo"]
					if !ok {
						t.Fatal("no managed resource type 'foo' in response")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["name"]
					if !ok {
						t.Fatal("no attribute named 'name' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeRequired; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"data resource types": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							DataSourceSchemas: map[string]*tfplugin6.Schema{
								"foo": {
									Block: &tfplugin6.Schema_Block{
										Attributes: []*tfplugin6.Schema_Attribute{
											{
												Name:     "name",
												Type:     mockutil.JSON("string"),
												Required: true,
											},
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					got := maps.Collect(resp.ProviderSchema().DataResourceTypeSchemas())
					schema, ok := got["foo"]
					if !ok {
						t.Fatal("no data resource type 'foo' in response")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["name"]
					if !ok {
						t.Fatal("no attribute named 'name' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeRequired; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"ephemeral resource types": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							EphemeralResourceSchemas: map[string]*tfplugin6.Schema{
								"foo": {
									Block: &tfplugin6.Schema_Block{
										Attributes: []*tfplugin6.Schema_Attribute{
											{
												Name:     "name",
												Type:     mockutil.JSON("string"),
												Required: true,
											},
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					got := maps.Collect(resp.ProviderSchema().EphemeralResourceTypeSchemas())
					schema, ok := got["foo"]
					if !ok {
						t.Fatal("no ephemeral resource type 'foo' in response")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["name"]
					if !ok {
						t.Fatal("no attribute named 'name' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeRequired; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"list resource types": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in provider_schema_impl_test.go. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							ListResourceSchemas: map[string]*tfplugin6.Schema{
								"foo": {
									Block: &tfplugin6.Schema_Block{
										Attributes: []*tfplugin6.Schema_Attribute{
											{
												Name:     "name",
												Type:     mockutil.JSON("string"),
												Required: true,
											},
										},
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					got := maps.Collect(resp.ProviderSchema().ManagedResourceTypeListSchemas())
					schema, ok := got["foo"]
					if !ok {
						t.Fatal("no list resource type 'foo' in response")
					}
					gotAttrs := maps.Collect(schema.Attributes())
					attrS, ok := gotAttrs["name"]
					if !ok {
						t.Fatal("no attribute named 'name' in response")
					}
					if got, want := attrS.Usage(), providerschema.AttributeRequired; got != want {
						t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
					}
					ty, err := attrS.Type().AsCtyType()
					if err != nil {
						t.Fatalf("invalid attribute type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
			"functions": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.GetProviderSchema(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.GetProviderSchema_Request{}),
					).Return(
						&tfplugin6.GetProviderSchema_Response{
							// NOTE: This doesn't need to have comprehensive
							// coverage of every possible part of schema because
							// we have separate tests for the mapping from
							// tfplugin6.Schema to the version-agnostic
							// interfaces in other functions in this file. We're
							// just checking whether resource types from the
							// response make it into the return value at all.
							Functions: map[string]*tfplugin6.Function{
								"boop": {
									Return: &tfplugin6.Function_Return{
										Type: mockutil.JSON("string"),
									},
								},
							},
						}, nil,
					)
				},
				Request: &providerops.GetProviderSchemaRequest{},
				Check: func(t *testing.T, resp providerops.GetProviderSchemaResponse) {
					mockutil.AssertNoDiags(t, resp.Diagnostics())
					got := maps.Collect(resp.ProviderSchema().FunctionSignatures())
					sig, ok := got["boop"]
					if !ok {
						t.Fatal("no function 'boop' in response")
					}
					checkEmptySeq(t, "function parameter", sig.Parameters())
					ty, err := sig.ResultType().AsCtyType()
					if err != nil {
						t.Fatalf("invalid result type: %s", err)
					}
					if got, want := ty, cty.String; got != want {
						t.Errorf("wrong result type\ngot:  %#v\nwant: %#v", got, want)
					}
				},
			},
		},
		func(t *testing.T, provider *tf6.Provider, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error) {
			// So we can test trace span propagation, some of the tests
			// use mocks that require a context with this key/value pair:
			ctx := context.WithValue(t.Context(), "propagation_test", true)
			return provider.GetProviderSchema(ctx, req)
		},
	)
}
