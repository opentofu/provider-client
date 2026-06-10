package providerschema

import (
	"iter"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

type IdentitySchema interface {
	Version() int64
	IdentityAttributes() iter.Seq[IdentityAttribute]

	common.Sealed
}

type IdentityAttribute interface {
	Name() string
	Type() TypeConstraint
	RequiredForImport() bool
	OptionalForImport() bool
	Description() string

	common.Sealed
}

type ResourceIdentityData interface {
	IdentityData() DynamicValueOut

	common.Sealed
}
