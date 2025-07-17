// Package tofuprovider is a low-level client library for the OpenTofu provider
// plugin API, allowing Go programs to call into OpenTofu provider plugins
// without using code from OpenTofu itself.
//
// The scope of this library is intentionally limited: it focuses only on
// abstracting away the protocol major version negotiation so that caller
// code does not need to be duplicated to support multiple protocol versions.
// Otherwise, the API closely matches how the underlying protocol is designed
// with as little additional abstraction as possible.
//
// This package currently implements clients for protocol major versions 5 and
// 6.
package tofuprovider

import (
	"context"
	"fmt"
	"os/exec"

	"go.rpcplugin.org/rpcplugin"

	// The following is required to force google.golang.org/genproto to
	// appear in our go.mod, which is in turn needed to resolve ambiguous
	// package imports in google.golang.org/grpc which can potentially
	// match two different module layouts as the module boundaries
	// under this prefix have changed over time.
	_ "google.golang.org/genproto/protobuf/ptype"
)

// Provider represents a running provider plugin.
type Provider interface {
	Close() error
}

// Start executes the given command line as an OpenTofu provider plugin
// and returns an object representing it.
//
// The provider is initially unconfigured, meaning that it can only be used
// for object validation tasks. It must be configured (that is, it must be
// provided with a valid configuration object) before it can take any
// non-validation actions.
//
// OpenTofu providers run as child processes, so if this function returns
// successfully there will be a new child process beneath the calling process
// waiting to recieve provider commands. Be sure to call Close on the returned
// object when you no longer need the provider, so that the child process
// can be terminated.
func Start(ctx context.Context, exe string, args ...string) (Provider, error) {
	plugin, err := rpcplugin.New(ctx, &rpcplugin.ClientConfig{
		Handshake: rpcplugin.HandshakeConfig{
			CookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
			CookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
		},
		Cmd:           exec.Command(exe, args...),
		ProtoVersions: map[int]rpcplugin.ClientVersion{
			//5: protocol5.PluginClient{},
			//6: protocol6.PluginClient{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to launch provider plugin: %s", err)
	}

	protoVersion, _ /*clientProxy*/, err := plugin.Client(ctx)
	if err != nil {
		plugin.Close()
		return nil, fmt.Errorf("failed to create plugin client: %s", err)
	}

	switch protoVersion {
	//case 5:
	//	return protocol5.NewProvider(ctx, plugin, clientProxy)
	//case 6:
	//	return protocol6.NewProvider(ctx, plugin, clientProxy)
	default:
		// Should not be possible to get here because the above cases cover
		// all of the versions we listed in ProtoVersions; rpcplugin bug?
		panic(fmt.Sprintf("unsupported protocol version %d", protoVersion))
	}
}
