package tf6

import (
	"context"
	"errors"

	"go.rpcplugin.org/rpcplugin"
	"google.golang.org/grpc"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

type Provider struct {
	client tfplugin6.ProviderClient
	plugin *rpcplugin.Plugin

	common.SealedImpl
}

// NewProvider builds a new [Provider] wrapping the given plugin process and
// client proxy.
//
// Pass a nil plugin if the given client proxy is for an in-process
// implementation, such as a mock used for testing. If you pass a non-nil plugin
// then it will be closed automatically when the provider is closed but won't
// be used for any other purpose.
func NewProvider(ctx context.Context, plugin *rpcplugin.Plugin, clientProxy any) (*Provider, error) {
	return &Provider{
		client: clientProxy.(tfplugin6.ProviderClient),
		plugin: plugin,
	}, nil
}

func (p *Provider) ProtocolMajorVersion() int {
	return 6
}

func (p *Provider) ClientProxy() any {
	return p.client
}

func (p *Provider) Close() error {
	plugin := p.plugin
	p.plugin = nil
	p.client = nil // subsequent usage of the client will panic
	if plugin == nil {
		return nil // it's okay to call Close multiple times on the same provider instance
	}
	return plugin.Close()
}

func (p *Provider) GracefulStop(ctx context.Context) error {
	resp, err := p.client.StopProvider(ctx, &tfplugin6.StopProvider_Request{})
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}

// PluginClient is an adapter used by the main package to obtain the low-level
// gRPC client proxy when protocol version 6 is selected.
type PluginClient struct{}

func (c PluginClient) ClientProxy(ctx context.Context, conn *grpc.ClientConn) (any, error) {
	return tfplugin6.NewProviderClient(conn), nil
}
