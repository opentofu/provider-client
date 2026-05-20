package tf6_test

import (
	"context"
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf6"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func TestCallFunction(t *testing.T) {
	testProviderCalls(t,
		map[string]providerCallTest[*providerops.CallFunctionRequest, providerops.CallFunctionResponse]{
			"no arguments": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.CallFunction(
						mockutil.ContextWithValue("propagation_test", true),
						mockutil.Eq(&tfplugin6.CallFunction_Request{
							Name: "do_something",
						}),
					).Return(
						&tfplugin6.CallFunction_Response{
							Result: &tfplugin6.DynamicValue{
								Msgpack: mockutil.MsgPack("done_something"),
							},
						}, nil,
					)
				},
				Request: &providerops.CallFunctionRequest{
					FunctionName: "do_something",
				},
				Check: func(t *testing.T, resp providerops.CallFunctionResponse) {
					gotResult := resp.Result()
					gotV, err := gotResult.AsCtyValue(cty.String)
					if err != nil {
						t.Fatalf("result has invalid syntax: %s", err)
					}
					if wantV := cty.StringVal("done_something"); !wantV.RawEquals(gotV) {
						t.Errorf("wrong result\ngot:  %#v\nwant: %#v", gotV, wantV)
					}
					if gotErr := resp.Error(); gotErr != nil {
						t.Errorf("unexpected error: %#v", gotErr)
					}
				},
			},
			"one argument": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.CallFunction(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.CallFunction_Request{
							Name: "do_something",
							Arguments: []*tfplugin6.DynamicValue{
								{Msgpack: mockutil.MsgPack("hello")},
							},
						}),
					).Return(
						&tfplugin6.CallFunction_Response{
							Result: &tfplugin6.DynamicValue{
								Msgpack: mockutil.MsgPack("done_something"),
							},
						}, nil,
					)
				},
				Request: &providerops.CallFunctionRequest{
					FunctionName: "do_something",
					Arguments: []providerschema.DynamicValueIn{
						providerschema.NewDynamicValue(cty.StringVal("hello"), cty.String),
					},
				},
				Check: func(t *testing.T, resp providerops.CallFunctionResponse) {
					gotResult := resp.Result()
					gotV, err := gotResult.AsCtyValue(cty.String)
					if err != nil {
						t.Fatalf("result has invalid syntax: %s", err)
					}
					if wantV := cty.StringVal("done_something"); !wantV.RawEquals(gotV) {
						t.Errorf("wrong result\ngot:  %#v\nwant: %#v", gotV, wantV)
					}
					if gotErr := resp.Error(); gotErr != nil {
						t.Errorf("unexpected error: %#v", gotErr)
					}
				},
			},
			"two arguments": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.CallFunction(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.CallFunction_Request{
							Name: "do_something",
							Arguments: []*tfplugin6.DynamicValue{
								{Msgpack: mockutil.MsgPack("hello")},
								{Msgpack: mockutil.MsgPack("world")},
							},
						}),
					).Return(
						&tfplugin6.CallFunction_Response{
							Result: &tfplugin6.DynamicValue{
								Msgpack: mockutil.MsgPack("done_something"),
							},
						}, nil,
					)
				},
				Request: &providerops.CallFunctionRequest{
					FunctionName: "do_something",
					Arguments: []providerschema.DynamicValueIn{
						providerschema.NewDynamicValue(cty.StringVal("hello"), cty.String),
						providerschema.NewDynamicValue(cty.StringVal("world"), cty.String),
					},
				},
				Check: func(t *testing.T, resp providerops.CallFunctionResponse) {
					gotResult := resp.Result()
					gotV, err := gotResult.AsCtyValue(cty.String)
					if err != nil {
						t.Fatalf("result has invalid syntax: %s", err)
					}
					if wantV := cty.StringVal("done_something"); !wantV.RawEquals(gotV) {
						t.Errorf("wrong result\ngot:  %#v\nwant: %#v", gotV, wantV)
					}
					if gotErr := resp.Error(); gotErr != nil {
						t.Errorf("unexpected error: %#v", gotErr)
					}
				},
			},
			"failed call with general error": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.CallFunction(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.CallFunction_Request{
							Name: "do_something",
						}),
					).Return(
						&tfplugin6.CallFunction_Response{
							Error: &tfplugin6.FunctionError{
								Text: "failed to do something",
							},
						}, nil,
					)
				},
				Request: &providerops.CallFunctionRequest{
					FunctionName: "do_something",
				},
				Check: func(t *testing.T, resp providerops.CallFunctionResponse) {
					gotErr := resp.Error()
					wantErr := &mockutil.ComparableFunctionError{
						Text_: "failed to do something",
					}
					if diff := mockutil.Diff(wantErr, gotErr); diff != "" {
						t.Error("wrong error\n" + diff)
					}
				},
			},
			"failed call with argument-specific error": {
				Mock: func(expect *MockProviderClientMockRecorder) {
					expect.CallFunction(
						mockutil.AnyContext(),
						mockutil.Eq(&tfplugin6.CallFunction_Request{
							Name: "do_something",
							Arguments: []*tfplugin6.DynamicValue{
								{Msgpack: mockutil.MsgPack("hey dude")},
							},
						}),
					).Return(
						&tfplugin6.CallFunction_Response{
							Error: &tfplugin6.FunctionError{
								Text:             "insufficiently polite greeting",
								FunctionArgument: new(int64(0)),
							},
						}, nil,
					)
				},
				Request: &providerops.CallFunctionRequest{
					FunctionName: "do_something",
					Arguments: []providerschema.DynamicValueIn{
						providerschema.NewDynamicValue(cty.StringVal("hey dude"), cty.String),
					},
				},
				Check: func(t *testing.T, resp providerops.CallFunctionResponse) {
					gotErr := resp.Error()
					wantErr := &mockutil.ComparableFunctionError{
						Text_:          "insufficiently polite greeting",
						ArgumentIndex_: new(0),
					}
					if diff := mockutil.Diff(wantErr, gotErr); diff != "" {
						t.Error("wrong error\n" + diff)
					}
				},
			},
		},
		func(t *testing.T, provider *tf6.Provider, req *providerops.CallFunctionRequest) (providerops.CallFunctionResponse, error) {
			// So we can test trace span propagation, some of the tests
			// use mocks that require a context with this key/value pair:
			ctx := context.WithValue(t.Context(), "propagation_test", true)
			return provider.CallFunction(ctx, req)
		},
	)
}
