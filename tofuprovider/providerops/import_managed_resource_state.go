package providerops

import (
	"iter"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

type ImportManagedResourceStateRequest struct {
	// ResourceType is the name of the type of resource the object should
	// be imported as.
	ResourceType string

	// ID is some sort of unique identifier for the object to be imported,
	// in a format decided by the provider.
	ID string

	// ClientCapabilities allows the caller to declare that it is capable of
	// handling certain response data that was added to the protocol after
	// it was initially defined, and thus which the provider must disable
	// by default to avoid confusing older clients.
	ClientCapabilities *ClientCapabilities

	// TODO: Identity
}

type ImportManagedResourceStateResponse interface {
	// Diagnostics are any diagnostics included in the provider's response.
	//
	// If the result's [Diagnostics.HasErrors] method returns true then
	// the results of all other methods are unspecified and meaningless.
	Diagnostics() Diagnostics

	// ImportedResources returns an interable sequence of the managed resource
	// objects that were imported.
	//
	// Although ImportManagedResourceState specifies just a single object
	// to import, providers are allowed to propose importing a number of
	// separate objects. It's never been very well defined when providers
	// ought to do this and what it means when they do, so different providers
	// use this ability in different ways. Some clients might prefer to
	// filter the result to include only objects whose resource type matches
	// the request, and possibly to fail if there isn't exactly one left
	// after that filtering, to approximate the effect of importing only
	// one specific object.
	ImportedResources() iter.Seq[ImportedManagedResource]

	common.Sealed
}

type ImportedManagedResource interface {
	ResourceType() string

	// Refer to [ApplyManagedResourceChangeResponse.NewState] for details on
	// how to use this.
	State() providerschema.DynamicValueOut

	// Refer to [ApplyManagedResourceChangeResponse.ProviderInternal] for
	// details on how to use this.
	ProviderInternal() []byte

	// TODO: Identity

	common.Sealed
}
