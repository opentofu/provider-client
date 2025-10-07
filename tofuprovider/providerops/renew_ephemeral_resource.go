package providerops

import (
	"time"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

type RenewEphemeralResourceRequest struct {
	// ResourceType is the name of the type of resource the client wants to
	// renew.
	ResourceType string

	// ProviderInternal is the exact blob returned the most recent call to
	// either [OpenEphemeralResourceResponse.ProviderInternal] or
	// [RenewEphemeralResourceResponse.ProviderInternal] relating to this
	// object.
	//
	// The provider uses this to determine which object is being renewed.
	ProviderInternal []byte
}

type RenewEphemeralResourceResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// RenewTime returns a timestamp along with true if the object associated
	// with this ephemeral resource needs to be renewed again in future, or a
	// meaningless value with false if renewal is not needed.
	//
	// If renewal is needed then the client must make sure to call
	// RenewEphemeralResource at some point before the given timestamp to
	// ensure that the associated object remains valid.
	RenewTime() (time.Time, bool)

	// ProviderInternal returns an opaque blob that must be sent back verbatim
	// in any subsequent RenewEphemeralResource or CloseEphemeralResource
	// call.
	//
	// This replaces the similar value returned by any earlier call to
	// OpenEphemeralResource or RenewEphemeralResource for the same object.
	ProviderInternal() []byte

	common.Sealed
}
