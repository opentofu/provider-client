package providerops

import (
	"context"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

type ListManagedResourcesRequest struct {
	// Config is the list ConfigSchema-based configuration data.
	Config providerschema.DynamicValueIn
	// TypeName is the list resource type name.
	TypeName string
	// IncludeResourceObject when set to true, the provider should
	// include the full resource object for each result
	IncludeResourceObject bool
	// Limit is the maximum number of results that Terraform is expecting.
	// The stream will stop, once this limit is reached.
	Limit int64
}

type ListManagedResourcesEvent interface {
	DisplayName() string
	Resource() providerschema.DynamicValueOut
	Diagnostics() Diagnostics

	common.Sealed
}

type ListManagedResourcesResponse interface {
	// ReadResult reads the next result from the stream. It returns io.EOF
	// once all results have been read.
	//
	// The context is accepted so that non-streaming implementations (e.g.
	// ones backed by traditional pagination) can perform fallible, possibly
	// cancellable work per read.
	ReadResult(ctx context.Context) (ListManagedResourcesEvent, error)

	// Close terminates the underlying stream. Callers must call Close once
	// they are done reading, especially if they stopped before reaching io.EOF,
	// otherwise the stream may stay active until the context passed to
	// ListManagedResources is cancelled.
	//
	// The context and error result are accepted so that non-streaming
	// implementations can do fallible, possibly cancellable cleanup.
	Close(ctx context.Context) error

	common.Sealed
}
