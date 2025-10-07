package tofuprovider

import (
	"context"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"

	// The following is required to force google.golang.org/genproto to
	// appear in our go.mod, which is in turn needed to resolve ambiguous
	// package imports in google.golang.org/grpc which can potentially
	// match two different module layouts as the module boundaries
	// under this prefix have changed over time.
	_ "google.golang.org/genproto/protobuf/ptype"
)

// Provider represents operations on a provider that cause requests to be
// sent to the running provider plugin, regardless of the specific execution
// model used for the provider.
type Provider interface {
	// GetProviderSchema requests the full provider schema, as a single
	// large object. A successful result describes all features that the
	// provider offers, and how clients are expected to interact with those
	// features.
	//
	// The response also includes a [ServerCapabilities] object which the
	// caller should use to recognize certain limitations in a particular
	// provider's support of the provider protocol.
	//
	// This method should be called before calling [ConfigureProvider].
	GetProviderSchema(ctx context.Context, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error)

	// ValidateProviderConfig tests whether a given provider configuration
	// object is acceptable per the provider's internally-implemented
	// validation rules.
	//
	// This method should be called before calling [ConfigureProvider].
	ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error)

	// ConfigureProvider asks the provider to transition from the "unconfigured"
	// state to the "configured" state, using a given configuration value.
	//
	// This should be called only once per provider instance. Repeated calls
	// cause unspecified behavior. Most other methods of this type are only
	// valid to call after a provider has been successfully configured.
	//
	// There is no way to revert from "configured" to "unconfigured", so callers
	// that need ongoing access to the unconfigured operations should retain
	// a separate instance for which ConfigureProvider is never called.
	ConfigureProvider(ctx context.Context, req *providerops.ConfigureProviderRequest) (providerops.ConfigureProviderResponse, error)

	// ValidateManagedResourceConfig tests whether a given managed resource
	// configuration object is acceptable per the provider's internally-implemented
	// validation rules.
	ValidateManagedResourceConfig(ctx context.Context, req *providerops.ValidateManagedResourceConfigRequest) (providerops.ValidateManagedResourceConfigResponse, error)

	// UpgradeManagedResourceState asks the provider to prepare some raw data
	// previously saved for a managed resource instance to suit the schema
	// of its resource type in the current version of the provider.
	//
	// This operation deals with the problem that the client has no way to
	// automatically find the schema for an earlier version of a provider.
	// Instead, new provider versions are expected to know how earlier versions
	// of the same provider represented each resource type and know how to
	// transform that raw data into a valid representation for the current
	// schema.
	UpgradeManagedResourceState(ctx context.Context, req *providerops.UpgradeManagedResourceStateRequest) (providerops.UpgradeManagedResourceStateResponse, error)

	// ReadManagedResource trades a previously-saved state object of a
	// managed resource type for a new object updated to match the current
	// configuration of the remote object.
	//
	// This is intended for detecting changes that were made outside of
	// OpenTofu, although some providers have bugs that cause spurious
	// differences so callers should avoid assuming that all differences
	// are problematic.
	ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error)

	// ImportManagedResourceState attempts to produce a suitable resource state
	// representation of an object that was originally created outside of
	// OpenTofu but that will be managed by OpenTofu in future.
	//
	// This is essentially an alternative to ReadManagedResource that uses
	// only identification information to locate the intended object, rather
	// than relying on full prior state data.
	ImportManagedResourceState(ctx context.Context, req *providerops.ImportManagedResourceStateRequest) (providerops.ImportManagedResourceStateResponse, error)

	// PlanManagedResourceChange asks the provider to compare prior state
	// and configuration and produce a merged "planned new state" that
	// could be reached by a subsequent call to ApplyManagedResourceChange.
	PlanManagedResourceChange(ctx context.Context, req *providerops.PlanManagedResourceChangeRequest) (providerops.PlanManagedResourceChangeResponse, error)

	// ApplyManagedResourceChange asks the provider to apply a change
	// previously planned by PlanManagedResourceChange, thereby modifying
	// an associated remote object to match the desired state.
	ApplyManagedResourceChange(ctx context.Context, req *providerops.ApplyManagedResourceChangeRequest) (providerops.ApplyManagedResourceChangeResponse, error)

	// MoveManagedResourceState attempts to transform some state data originally
	// created for another resource type and possibly in another provider to
	// suit a resource type in this provider.
	//
	// OpenTofu uses this when a "moved" block declares that an object moved
	// between two resource addresses of different types. It's the target
	// provider's responsibility to understand the source provider's state
	// storage format, so that the source provider isn't actually needed to
	// complete the action. (This functionality is sometimes used when the
	// source provider is no longer usable for some reason, such as if it's
	// deprecated has no releases available for the current platform.)
	MoveManagedResourceState(ctx context.Context, req *providerops.MoveManagedResourceStateRequest) (providerops.MoveManagedResourceStateResponse, error)

	// ValidateDataResourceConfig tests whether a given data resource
	// configuration object is acceptable per the provider's internally-implemented
	// validation rules.
	//
	// This method should be called before calling [ConfigureProvider].
	ValidateDataResourceConfig(ctx context.Context, req *providerops.ValidateDataResourceConfigRequest) (providerops.ValidateDataResourceConfigResponse, error)

	// ReadDataResource asks the provider to read from a data source
	// corresponding to a particular data resource type.
	ReadDataResource(ctx context.Context, req *providerops.ReadDataResourceRequest) (providerops.ReadDataResourceResponse, error)

	// ValidateEphemeralResourceConfig tests whether a given ephemeral resource
	// configuration object is acceptable per the provider's internally-implemented
	// validation rules.
	//
	// This method should be called before calling [ConfigureProvider].
	ValidateEphemeralResourceConfig(ctx context.Context, req *providerops.ValidateEphemeralResourceConfigRequest) (providerops.ValidateEphemeralResourceConfigResponse, error)

	// OpenEphemeralResource asks the provider to "open" (create, acquire, etc)
	// something represented by an ephemeral resource type.
	OpenEphemeralResource(ctx context.Context, req *providerops.OpenEphemeralResourceRequest) (providerops.OpenEphemeralResourceResponse, error)

	// RenewEphemeralResource asks the provider to renew whatever is associated
	// with an ephemeral resource that was previously opened and that indicated
	// that it needs periodic renewal.
	RenewEphemeralResource(ctx context.Context, req *providerops.RenewEphemeralResourceRequest) (providerops.RenewEphemeralResourceResponse, error)

	// CloseEphemeralResource asks the provider to "close" (delete, release, etc)
	// whatever is associated with an ephemeral resource that was previously
	// opened.
	CloseEphemeralResource(ctx context.Context, req *providerops.CloseEphemeralResourceRequest) (providerops.CloseEphemeralResourceResponse, error)

	// GetFunctions is essentially a lighter version of GetProviderSchema
	// that describes only the provider's exported functions.
	//
	// Most callers should fall back on using GetProviderSchema if this method
	// returns an error that causes [providerops.IsUnimplementedErr] to
	// return true.
	//
	// This method should be called before calling [ConfigureProvider].
	GetFunctions(ctx context.Context, req *providerops.GetFunctionsRequest) (providerops.GetFunctionsResponse, error)

	// CallFunction calls one of the functions supported by the provider using
	// a given set of arguments, returning the function's result.
	CallFunction(ctx context.Context, req *providerops.CallFunctionRequest) (providerops.CallFunctionResponse, error)

	// GracefulStop asks the provider to gracefully abort any active
	// calls that are running concurrently, causing them to return
	// with a cancellation-related error as soon as it's safe to do so.
	//
	// Not all providers actually support cancellation for all of their
	// resource types, so a caller must not assume that concurrent calls
	// will definitely return promptly after calling this method.
	//
	// It's safe to call GracefulStop multiple times on the same provider,
	// although for most providers the additional calls have no additional
	// effect. Calling other methods after calling GracefulStop causes
	// unspecified behavior: the provider might choose to immediately return
	// a cancellation-related error, or it might exhibit strange behavior
	// such as only partially completing a request.
	GracefulStop(ctx context.Context) error

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
