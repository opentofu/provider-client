package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// ServerCapabilities is a base implementation of
// [providerops.ServerCapabilities] that is intended to be embedded in external
// implementations of that interface.
//
// It implements all methods of the interface, and will be expanded to return
// reasonable default values for any other methods added in future. Refer to
// the documentation of each method to learn what the default results are, and
// override in your subtype as needed.
type ServerCapabilities struct {
	sealedImpl
}

var _ providerops.ServerCapabilities = ServerCapabilities{}

// CanMoveManagedResourceState implements [providerops.ServerCapabilities] by
// returning false, because it only makes sense to return true here for
// providers that have a useful implementation of migrating between resource
// types.
func (ServerCapabilities) CanMoveManagedResourceState() bool {
	return false
}

// CanPlanDestroy implements [providerops.ServerCapabilities] by returning
// true, because this capability was added to accommodate some legacy
// assumptions made by older providers but newly-written providers should
// always support destroy-planning.
func (ServerCapabilities) CanPlanDestroy() bool {
	return true
}

// GetProviderSchemaIsOptional implements [providerops.ServerCapabilities]
// by returning true, because this capability was added only to support some
// quirky misbehavior in one specific legacy provider and so newly-written
// providers should avoid relying on it and should behave correctly even if
// the caller never requests provider schema.
func (ServerCapabilities) GetProviderSchemaIsOptional() bool {
	return true
}
