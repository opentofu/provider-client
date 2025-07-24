package tf6

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
)

type deferred struct {
	proto *tfplugin6.Deferred
	common.Sealed
}

// Reason implements providerops.Deferred.
func (d deferred) Reason() providerops.DeferredReason {
	return deferredReason(d.proto.Reason)
}

func deferredReason(proto tfplugin6.Deferred_Reason) providerops.DeferredReason {
	switch proto {
	case tfplugin6.Deferred_PROVIDER_CONFIG_UNKNOWN:
		return providerops.DeferredBecauseProviderConfigUnknown
	case tfplugin6.Deferred_RESOURCE_CONFIG_UNKNOWN:
		return providerops.DeferredBecauseResourceConfigUnknown
	case tfplugin6.Deferred_ABSENT_PREREQ:
		return providerops.DeferredBecauseAbsentPrereq
	default:
		return providerops.DeferredUnsupportedReason
	}
}
