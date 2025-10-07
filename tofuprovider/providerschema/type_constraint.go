package providerschema

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

// TypeConstraint describes a type constraint that some value is required
// to conform to.
type TypeConstraint interface {
	// AsCtyType returns a [cty.Type] representation of the type constraint.
	//
	// Some implementations delay parsing a wire-format representation of the
	// type until this method is called, and so an error from this method
	// represents that the returned wire representation was somehow invalid.
	AsCtyType() (cty.Type, error)

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}
