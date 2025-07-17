package tofuprovider

import (
	"context"
	"fmt"
	"os/exec"

	"go.rpcplugin.org/rpcplugin"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/tf5"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/tf6"

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

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
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
		Cmd: exec.Command(exe, args...),
		ProtoVersions: map[int]rpcplugin.ClientVersion{
			5: tf5.PluginClient{}, // clientProxy is tfplugin5.ProviderClient
			6: tf6.PluginClient{}, // clientProxy is tfplugin6.ProviderClient
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to launch provider plugin: %s", err)
	}

	// If plugin init and handshake is successful then clientProxy is
	// of the type described in the comments associated with the
	// matching rpcplugin.ClientConfig.ProtoVersions element above,
	// for returned protoVersion.
	protoVersion, clientProxy, err := plugin.Client(ctx)
	if err != nil {
		plugin.Close()
		return nil, fmt.Errorf("failed to create plugin client: %s", err)
	}

	var ret Provider
	switch protoVersion {
	case 5:
		// These extra steps are to avoid returning a "typed nil" if
		// NewProvider returns (*tf6.Provider)(nil).
		impl, err := tf5.NewProvider(ctx, plugin, clientProxy)
		if impl != nil {
			ret = impl
		}
		return ret, err
	case 6:
		// These extra steps are to avoid returning a "typed nil" if
		// NewProvider returns (*tf6.Provider)(nil).
		impl, err := tf6.NewProvider(ctx, plugin, clientProxy)
		if impl != nil {
			ret = impl
		}
		return ret, err
	default:
		// Should not be possible to get here because the above cases cover
		// all of the versions we listed in ProtoVersions; rpcplugin bug?
		panic(fmt.Sprintf("unsupported protocol version %d", protoVersion))
	}
}
