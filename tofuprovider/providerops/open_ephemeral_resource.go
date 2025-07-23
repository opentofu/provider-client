package providerops

import (
	"time"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type OpenEphemeralResourceRequest struct {
	// ResourceType is the name of the type of resource the client wants to
	// open.
	ResourceType string

	// Config is a dynamic value representing the resource configuration.
	//
	// Providers expect that the given will conform to the current schema
	// for the given resource type.
	Config providerschema.DynamicValueIn

	// ClientCapabilities allows the caller to declare that it is capable of
	// handling certain response data that was added to the protocol after
	// it was initially defined, and thus which the provider must disable
	// by default to avoid confusing older clients.
	ClientCapabilities *ClientCapabilities
}

type OpenEphemeralResourceResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// Result returns an object describing what was opened.
	//
	// This must be decoded using the type implied by the schema of the
	// resource type.
	Result() providerschema.DynamicValueOut

	// RenewTime returns a timestamp along with true if the object associated
	// with this ephemeral resource needs to be renewed in future, or a
	// meaningless value with false if renewal is not needed.
	//
	// If renewal is needed then the client must make sure to call
	// RenewEphemeralResource at some point before the given timestamp to
	// ensure that the associated object remains valid.
	RenewTime() (time.Time, bool)

	// ProviderInternal returns an opaque blob that must be sent back verbatim
	// in any subsequent RenewEphemeralResource or CloseEphemeralResource
	// call.
	ProviderInternal() []byte

	// Deferred returns a non-nil value if the provider does not have enough
	// information to satisfy this request, such as if the provider
	// configuration is not yet known enough to know which API endpoint to
	// connect to.
	//
	// If this returns nil then other methods return updated data that
	// should replace the previous values that were saved in the prior state.
	Deferred() Deferred

	common.Sealed
}
