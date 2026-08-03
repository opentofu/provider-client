package providerdefs

import (
	"context"
	"io"

	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// ListManagedResourcesEvent is a base implementation of
// [providerops.ListManagedResourcesEvent].
type ListManagedResourcesEvent struct {
	sealedImpl
}

var _ providerops.ListManagedResourcesEvent = ListManagedResourcesEvent{}

// DisplayName implements [providerops.ListManagedResourcesEvent] by returning
// an empty string.
func (e ListManagedResourcesEvent) DisplayName() string {
	return ""
}

// Resource implements [providerops.ListManagedResourcesEvent] by returning nil,
// meaning the full resource object was not included in the event.
func (e ListManagedResourcesEvent) Resource() providerschema.DynamicValueOut {
	return nil
}

// Identity implements [providerops.ListManagedResourcesEvent] by returning nil,
// meaning no resource identity was included in the event.
func (e ListManagedResourcesEvent) Identity() providerschema.IdentityData {
	return nil
}

// Diagnostics implements [providerops.ListManagedResourcesEvent] by returning
// no diagnostics at all.
func (e ListManagedResourcesEvent) Diagnostics() providerops.Diagnostics {
	return emptyDiagnostics{}
}

// ListManagedResourcesResponse is a base implementation of
// [providerops.ListManagedResourcesResponse].
type ListManagedResourcesResponse struct {
	sealedImpl
}

var _ providerops.ListManagedResourcesResponse = ListManagedResourcesResponse{}

// ReadResult implements [providerops.ListManagedResourcesResponse] by
// immediately reporting the end of the stream.
func (r ListManagedResourcesResponse) ReadResult(ctx context.Context) (providerops.ListManagedResourcesEvent, error) {
	return nil, io.EOF
}

// Close implements [providerops.ListManagedResourcesResponse] by doing nothing.
func (r ListManagedResourcesResponse) Close(ctx context.Context) error {
	return nil
}
