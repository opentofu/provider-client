package tf6

import (
	"context"
	"errors"

	"go.rpcplugin.org/rpcplugin"
	"google.golang.org/grpc"

	"github.com/apparentlymart/opentofu-providers/internal/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
)

type Provider struct {
	client tfplugin6.ProviderClient
	plugin *rpcplugin.Plugin

	common.SealedImpl
}

func NewProvider(ctx context.Context, plugin *rpcplugin.Plugin, clientProxy any) (*Provider, error) {
	return &Provider{
		client: clientProxy.(tfplugin6.ProviderClient),
		plugin: plugin,
	}, nil
}

func (p *Provider) Close() error {
	if p.plugin == nil {
		return nil // it's okay to call Close multiple times on the same provider instance
	}
	plugin := p.plugin
	p.plugin = nil
	p.client = nil // subsequent usage of the client will panic
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
