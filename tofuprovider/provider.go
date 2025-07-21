package tofuprovider

import (
	"context"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"

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
	GetProviderSchema(ctx context.Context, req *providerops.GetProviderSchemaRequest) (providerops.GetProviderSchemaResponse, error)

	// ValidateProviderConfig tests whether a given provider configuration
	// object is acceptable per the provider's internally-implemented
	// validation rules.
	ValidateProviderConfig(ctx context.Context, req *providerops.ValidateProviderConfigRequest) (providerops.ValidateProviderConfigResponse, error)

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

	// GracefulStop asks the provider to gracefully abort any active
	// calls that are running concurrently, causing them to return
	// with a cancellation-related error as soon as it's safe to do so.
	//
	// Not all providers actually support cancellation for all of their
	// resource types, so a caller must not assume that concurrent calls
	// definitly will return promptly after calling this method.
	//
	// It's safe to call GracefulStop multiple times on the same provider,
	// although for most providers the additional calls have no additional
	// effect.
	GracefulStop(ctx context.Context) error

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
