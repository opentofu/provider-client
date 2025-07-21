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
