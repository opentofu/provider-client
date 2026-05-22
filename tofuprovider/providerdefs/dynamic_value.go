package providerdefs

import (
	"errors"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

// DynamicValueOut is a base implementation of [providerschema.DynamicValueOut].
type DynamicValueOut struct {
	sealedImpl
}

var _ providerschema.DynamicValueOut = DynamicValueOut{}

// AsCtyValue implements [providerschema.DynamicValueOut] by immediately
// returning an error.
//
// External implementers that embed [DynamicValueOut] should override this
// to do somthing useful.
func (d DynamicValueOut) AsCtyValue(withType cty.Type) (cty.Value, error) {
	return cty.NilVal, errors.New("no value available")
}
