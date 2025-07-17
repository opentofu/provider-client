package tf6

import (
	"context"
	"fmt"

	"go.rpcplugin.org/rpcplugin"
	"google.golang.org/grpc"

	"github.com/apparentlymart/opentofu-providers/internal/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
)

type Provider struct {
	common.SealedImpl
}

func NewProvider(ctx context.Context, plugin *rpcplugin.Plugin, clientProxy any) (*Provider, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (p *Provider) Close() error {
	return nil
}

// PluginClient is an adapter used by the main package to obtain the low-level
// gRPC client proxy when protocol version 6 is selected.
type PluginClient struct{}

func (c PluginClient) ClientProxy(ctx context.Context, conn *grpc.ClientConn) (any, error) {
	return tfplugin6.NewProviderClient(conn), nil
}
