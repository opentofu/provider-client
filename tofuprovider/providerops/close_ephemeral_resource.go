package providerops

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
)

type CloseEphemeralResourceRequest struct {
	// ResourceType is the name of the type of resource the client wants to
	// close.
	ResourceType string

	// ProviderInternal is the exact blob returned the most recent call to
	// either [OpenEphemeralResourceResponse.ProviderInternal] or
	// [RenewEphemeralResourceResponse.ProviderInternal] relating to this
	// object.
	//
	// The provider uses this to determine which object is being closed.
	ProviderInternal []byte
}

type CloseEphemeralResourceResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	common.Sealed
}
