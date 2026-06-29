package tf6_test

import (
	"context"
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"google.golang.org/grpc"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf6"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type fakeListReponseStream struct {
	grpc.ServerStreamingClient[tfplugin6.ListResource_Event]

	events []*tfplugin6.ListResource_Event
	err    error
}

func (f *fakeListReponseStream) Recv() (*tfplugin6.ListResource_Event, error) {
	if len(f.events) == 0 {
		return nil, io.EOF
	}
	ev := f.events[0]
	f.events = f.events[1:]
	return ev, nil
}

func TestListManagedResourcesImpl(t *testing.T) {
	testProviderCalls(t,
		map[string]providerCallTest[*providerops.ListManagedResourcesRequest, providerops.ListManagedResourcesResponse]{
			"list resource": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.ListResource(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.ListResource_Request{
							TypeName:              "test_type",
							Limit:                 1234,
							IncludeResourceObject: true,
							Config:                &tfplugin6.DynamicValue{Msgpack: mockutil.MsgPack(map[string]string{"foo": "bar"})},
						}),
					).Return(&fakeListReponseStream{
						events: []*tfplugin6.ListResource_Event{
							{
								DisplayName: "resource1",
								Diagnostic: []*tfplugin6.Diagnostic{
									{
										Severity: tfplugin6.Diagnostic_ERROR,
									},
								},
							},
							{
								DisplayName:    "resource2",
								ResourceObject: &tfplugin6.DynamicValue{Msgpack: mockutil.MsgPack(map[string]string{"id": "123"})},
							},
						},
					}, nil)
				},
				Request: &providerops.ListManagedResourcesRequest{
					TypeName:              "test_type",
					Limit:                 1234,
					IncludeResourceObject: true,
					Config: providerschema.NewDynamicValue(
						cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("bar")}),
						cty.Object(map[string]cty.Type{"foo": cty.String}),
					),
				},
				Check: func(t *testing.T, resp providerops.ListManagedResourcesResponse) {
					defer resp.Close(nil)
					var resources []providerops.ListManagedResourcesEvent
					for {
						ev, err := resp.ReadResult(nil)
						if errors.Is(err, io.EOF) {
							break
						}
						if err != nil {
							t.Fatalf("error reading resources: %s", err)
						}
						resources = append(resources, ev)
					}
					if got, expected := len(resources), 2; got != expected {
						t.Fatalf("wrong number of resources: got %d, want %d", got, expected)
					}

					if got, expected := resources[0].DisplayName(), "resource1"; got != expected {
						t.Errorf("wrong resource name: got %q, want %q", got, expected)
					}
					diags := slices.Collect(resources[0].Diagnostics().All())
					if got, expected := len(diags), 1; got != expected {
						t.Fatalf("wrong number of diagnostics: got %d, want %d", got, expected)
					}
					if got, expected := diags[0].Severity(), providerops.DiagnosticError; got != expected {
						t.Errorf("wrong diagnostic severity: got %s, want %s", got, expected)
					}

					if got, expected := resources[1].DisplayName(), "resource2"; got != expected {
						t.Errorf("wrong resource name: got %q, want %q", got, expected)
					}
					resourceObject, err := resources[1].Resource().AsCtyValue(cty.Object(map[string]cty.Type{"id": cty.String}))
					if err != nil {
						t.Fatalf("resource object invalid: %s", err)
					}
					if got, expected := resourceObject.GetAttr("id").AsString(), "123"; got != expected {
						t.Errorf("wrong resource object foo field: got %q, want %q", got, expected)
					}
				},
			},
		},
		func(t *testing.T, provider *tf6.Provider, req *providerops.ListManagedResourcesRequest) (providerops.ListManagedResourcesResponse, error) {
			// So we can test trace span propagation, some of the tests
			// use mocks that require a context with this key/value pair:
			ctx := context.WithValue(t.Context(), "propagation_test", true)
			return provider.ListManagedResources(ctx, req)
		})
}
