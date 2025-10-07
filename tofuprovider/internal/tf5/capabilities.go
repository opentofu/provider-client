package tf5

import (
	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

type serverCapabilities struct {
	proto *tfplugin5.ServerCapabilities

	common.SealedImpl
}

// CanMoveManagedResourceState implements providerops.ServerCapabilities.
func (s serverCapabilities) CanMoveManagedResourceState() bool {
	if s.proto == nil {
		return false
	}
	return s.proto.MoveResourceState
}

// CanPlanDestroy implements providerops.ServerCapabilities.
func (s serverCapabilities) CanPlanDestroy() bool {
	if s.proto == nil {
		return false
	}
	return s.proto.PlanDestroy
}

// GetProviderSchemaIsOptional implements providerops.ServerCapabilities.
func (s serverCapabilities) GetProviderSchemaIsOptional() bool {
	if s.proto == nil {
		return false
	}
	return s.proto.GetProviderSchemaOptional
}

func prepareClientCapabilities(caps *providerops.ClientCapabilities) *tfplugin5.ClientCapabilities {
	if caps == nil {
		return nil
	}
	return &tfplugin5.ClientCapabilities{
		DeferralAllowed:            caps.SupportsDeferral,
		WriteOnlyAttributesAllowed: caps.SupportsWriteOnlyAttributes,
	}
}
