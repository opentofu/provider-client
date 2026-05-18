package tf6_test

import (
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
