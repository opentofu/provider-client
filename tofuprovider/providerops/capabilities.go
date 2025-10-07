package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

type ServerCapabilities interface {
	// CanPlanDestroy returns true if the provider's PlanManagedResourceChange
	// implementation can support a request where the configuration object
	// and proposed new state are null, representing that the associated
	// object is to be deleted.
	//
	// If this returns false then the provider expects the caller to produce
	// a synthetic plan where the planned new state is a null object; the
	// provider gets no opportunity to return plan-time errors or warnings
	// about the proposed deletion in this case.
	CanPlanDestroy() bool

	// CanMoveManagedResourceState returns true if the provider has a working
	// implementation of the MoveManagedResourceState operation, or false
	// otherwise.
	CanMoveManagedResourceState() bool

	// GetProviderSchemaOptional returns true if the provider can function
	// correctly even when there has not been a call to its GetProviderSchema
	// operation. If this returns false then callers MUST call
	// GetProviderSchema next, before calling any other method of [Provider].
	//
	// Because Terraform and OpenTofu historically always called
	// GetProviderSchema immediately after launching a provider plugin process,
	// some provider developers inadvertently came to rely on some side-effects
	// of their schema generation process. Failing to call GetProviderSchema
	// when this method returns true can therefore cause strange malfunctions
	// in provider behavior for some existing provider plugins.
	//
	// Note that [ServerCapabilities] is itself exposed as part of the response
	// to GetProviderSchema, and a caller that obtained a [ServerCapabilities]
	// object through a valid [GetProviderSchemaResponse] object does not need
	// to call this method and does not need to make any additional call to
	// GetProviderSchema. This capability applies only to [ServerCapabilities]
	// objects obtained via other provider calls.
	GetProviderSchemaIsOptional() bool

	common.Sealed
}

type ClientCapabilities struct {
	// If SupportsDeferral is set, providers are allowed to return a non-nil
	// value from the "Deferred" method of the corresponding response type,
	// if present. When that field is set the caller must treat the response
	// as a possibly-incomplete prediction of what the resource will be once
	// more information is available, and so must not plan any immediate actions
	// based on the response.
	//
	// If this is not set then a provider is required to treat a deferral
	// situation as an error, describing the problem using one or more
	// error diagnostics.
	SupportsDeferral bool

	// If SupportsWriteOnlyAttributes is set, the client is able to correctly
	// handle write-only attributes.
	//
	// Providers typically use this to determine whether it's safe to accept
	// non-null values for write-only attributes, because older clients don't
	// understand that the values must not be persisted to plan files or state
	// snapshots. Therefore not setting this will typically cause providers
	// to return errors or warnings when write-only attributes are set.
	SupportsWriteOnlyAttributes bool
}
