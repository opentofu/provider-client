package tf5

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin5"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

type deferred struct {
	proto *tfplugin5.Deferred
	common.Sealed
}

// Reason implements providerops.Deferred.
func (d deferred) Reason() providerops.DeferredReason {
	return deferredReason(d.proto.Reason)
}

func deferredReason(proto tfplugin5.Deferred_Reason) providerops.DeferredReason {
	switch proto {
	case tfplugin5.Deferred_PROVIDER_CONFIG_UNKNOWN:
		return providerops.DeferredBecauseProviderConfigUnknown
	case tfplugin5.Deferred_RESOURCE_CONFIG_UNKNOWN:
		return providerops.DeferredBecauseResourceConfigUnknown
	case tfplugin5.Deferred_ABSENT_PREREQ:
		return providerops.DeferredBecauseAbsentPrereq
	default:
		return providerops.DeferredUnsupportedReason
	}
}
