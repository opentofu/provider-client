package providerops

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

type Deferred interface {
	Reason() DeferredReason

	common.Sealed
}

type DeferredReason int

const (
	// DeferredUnsupportedReason is a placeholder for when a provider returns
	// a reason code that this client library does not recognize.
	DeferredUnsupportedReason DeferredReason = 0

	// DeferredBecauseResourceConfigUnknown suggests that the provider needed
	// to defer the operation because of unknown values in the resource
	// instance's own configuration values.
	DeferredBecauseResourceConfigUnknown DeferredReason = 1

	// DeferredBecauseProviderConfigUnknown suggests that the provider needed
	// to defer the operation because of unknown values in the provider
	// configuration, such as when the target network endpoint is not yet
	// known.
	DeferredBecauseProviderConfigUnknown DeferredReason = 2

	// DeferredBecauseAbsentPrereq suggests that the provider needed to defer
	// the operation because a hard dependency has not been satisfied.
	DeferredBecauseAbsentPrereq DeferredReason = 3
)
