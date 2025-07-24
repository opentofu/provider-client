package tf6

import (
	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin6"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
)

type serverCapabilities struct {
	proto *tfplugin6.ServerCapabilities

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
