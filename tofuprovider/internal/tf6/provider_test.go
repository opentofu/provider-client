package tf6_test

import (
	"iter"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf6"
)

//go:generate go tool go.uber.org/mock/mockgen -destination client_mock_test.go -package tf6_test github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6 ProviderClient

func TestProviderBasics(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockProviderClient(ctrl)
	provider, err := tf6.NewProvider(t.Context(), nil, client)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got, want := provider.ClientProxy(), client; got != want {
		t.Error("ClientProxy did not return the originally-provided client")
	}

	client.EXPECT().StopProvider(gomock.Any(), mockutil.Eq(&tfplugin6.StopProvider_Request{})).Return(&tfplugin6.StopProvider_Response{}, nil)
	if err := provider.GracefulStop(t.Context()); err != nil {
		t.Errorf("GracefulStop failed: %s", err)
	}
	if err := provider.Close(); err != nil {
		// NOTE: For testing we set the "plugin" for the provider to be nil,
		// so this is expected to succeed without attempting to close it.
		t.Errorf("Close failed: %s", err)
	}
}

type providerCallTest[Req, Resp any] struct {
	Mock      func(expect *MockProviderClientMockRecorder)
	Request   Req
	Check     func(t *testing.T, resp Resp)
	WantError string
}

func testProviderCalls[Req, Resp any](
	t *testing.T,
	tests map[string]providerCallTest[Req, Resp],
	call func(t *testing.T, provider *tf6.Provider, req Req) (Resp, error),
) {
	t.Helper()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := NewMockProviderClient(ctrl)
			test.Mock(client.EXPECT())
			provider, err := tf6.NewProvider(t.Context(), nil, client)
			if err != nil {
				t.Fatalf("failed to create provider object: %s", err)
			}

			resp, err := call(t, provider, test.Request)
			if test.WantError != "" {
				if err == nil {
					t.Fatalf("unexpected success\nwant error: %s", test.WantError)
				}
				if got, want := err.Error(), test.WantError; got != want {
					t.Fatalf("wrong error\ngot:  %s\nwant: %s", got, want)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}
			if test.Check != nil {
				test.Check(t, resp)
			}
		})
	}
}

// checkEmptySeq is a testing helper which produces a log line and marks the
// test as failed if the given sequence has any items in it.
//
// This does not immediately halt the test if items are found. Check for return
// value of false to detect that case.
func checkEmptySeq[T any](t *testing.T, noun string, seq iter.Seq[T]) bool {
	t.Helper()
	isEmpty := true
	for v := range seq {
		t.Errorf("unexpected %s: %#v", noun, v)
		isEmpty = false
	}
	return isEmpty
}

// checkEmptySeq2 is a testing helper which produces a log line and marks the
// test as failed if the given sequence has any items in it.
//
// This does not immediately halt the test if items are found. Check for return
// value of false to detect that case.
func checkEmptySeq2[K, V any](t *testing.T, noun string, seq iter.Seq2[K, V]) bool {
	t.Helper()
	isEmpty := true
	for k, v := range seq {
		t.Errorf("unexpected %s: %#v -> %#v", noun, k, v)
		isEmpty = false
	}
	return isEmpty
}
