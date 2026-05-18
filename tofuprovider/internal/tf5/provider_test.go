package tf5_test

import (
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf5"
)

//go:generate go tool go.uber.org/mock/mockgen -destination client_mock_test.go -package tf5_test github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5 ProviderClient

func TestProviderBasics(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockProviderClient(ctrl)
	provider, err := tf5.NewProvider(t.Context(), nil, client)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got, want := provider.ClientProxy(), client; got != want {
		t.Error("ClientProxy did not return the originally-provided client")
	}

	client.EXPECT().Stop(gomock.Any(), mockutil.Eq(&tfplugin5.Stop_Request{})).Return(&tfplugin5.Stop_Response{}, nil)
	if err := provider.GracefulStop(t.Context()); err != nil {
		t.Errorf("GracefulStop failed: %s", err)
	}
	if err := provider.Close(); err != nil {
		// NOTE: For testing we set the "plugin" for the provider to be nil,
		// so this is expected to succeed without attempting to close it.
		t.Errorf("Close failed: %s", err)
	}
}
