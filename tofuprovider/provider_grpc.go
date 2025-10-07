package tofuprovider

import (
	"context"
	"fmt"
	"os/exec"

	"go.rpcplugin.org/rpcplugin"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf5"
	"github.com/opentofu/provider-client/tofuprovider/internal/tf6"
	"github.com/opentofu/provider-client/tofuprovider/providertrace"

	// The following is required to force google.golang.org/genproto to
	// appear in our go.mod, which is in turn needed to resolve ambiguous
	// package imports in google.golang.org/grpc which can potentially
	// match two different module layouts as the module boundaries
	// under this prefix have changed over time.
	_ "google.golang.org/genproto/protobuf/ptype"
)

// GRPCPluginProvider represents a running provider plugin that was started by
// [StartRPCPlugin].
type GRPCPluginProvider interface {
	// GRPCPluginProvider is a subtype of [Provider], which represents the
	// methods supported by all providers regardless of underlying execution
	// model and protocol transport.
	//
	// The other methods of [GRPCPluginProvider] below interact with the child
	// process that the provider plugin runs inside, handled locally
	// inside this library rather than remotely in the plugin.
	Provider

	// ProtocolMajorVersion returns the major version number of the wire
	// protocol that was negotiated during startup.
	//
	// In most cases the [Provider] abstraction should avoid callers needing
	// to vary their behavior by major version. This is exposed primarily
	// to allow callers to include it as diagnostic information in logs/etc
	ProtocolMajorVersion() int

	// ClientProxy returns the underlying gRPC client proxy object that this
	// provider is using to make the lower-level protocol requests.
	//
	// The result implements a different interface depending on which major
	// protocol version was negotiated, as returned by
	// [GRPCPluginProvider.ProtocolMajorVersion]:
	// - Version 6 client proxy implements [tfplugin6.ProviderClient]
	// - Version 5 client proxy implements [tfplugin5.ProviderClient]
	//
	// The set of supported versions could change in future, so callers using
	// this lower-level API should robustly handle recieving a client proxy
	// that implements neither of these interfaces. If the goal is to use
	// protocol features that are not part of this module's abstraction then
	// it may be better to use the tfplugin5/tfplugin6 APIs directly with
	// the underlying rpcplugin or go-plugin libraries and skip this
	// abstraction altogether.
	ClientProxy() any

	// Close terminates the child process representing the provider.
	//
	// After calling this function, the client object enters an invalid state
	// where all other methods have unspecified behavior. However, it's
	// acceptable to call close multiple times, with subsequent calls having
	// no effect.
	Close() error

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// StartGRPCPlugin executes the given command line, expecting it to behave
// as a "gRPC-style" provider plugin, and returns a [GRPCPluginProvider] object
// representing it.
//
// This function handles providers which use the gRPC-based protocols that
// originated in HashiCorp Terraform, automatically handling protocol
// schema version negotiation to provide a common interface over all
// supported protocol versions. This function currently supports protocol
// major versions 5 and 6.
//
// The provider is initially unconfigured, meaning that it can only be used
// for object validation tasks. It must be configured (that is, it must be
// provided with a valid configuration object) before it can take any
// non-validation actions.
//
// Provider plugins run as child processes, so if this function returns
// successfully there will be a new child process beneath the calling process
// waiting to receive provider commands. Be sure to call Close on the returned
// object when you no longer need the provider, so that the child process
// can be terminated.
func StartGRPCPlugin(ctx context.Context, exe string, args ...string) (GRPCPluginProvider, error) {
	tracer := providertrace.TracerFromContext(ctx)

	plugin, err := rpcplugin.New(ctx, &rpcplugin.ClientConfig{
		Handshake: rpcplugin.HandshakeConfig{
			CookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
			CookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
		},
		Cmd:    exec.Command(exe, args...),
		Stderr: tracer.ChildStderr,
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

	var ret GRPCPluginProvider
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
