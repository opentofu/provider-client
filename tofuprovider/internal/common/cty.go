package common

import (
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	ctymsgpack "github.com/zclconf/go-cty/cty/msgpack"
)

// CtyTypeJSON makes a byte array containing a standard JSON serialization of
// [cty.Type] implement [providerschema.TypeConstraint].
type CtyTypeJSON []byte

// AsCtyType implements providerschema.TypeConstraint.
func (c CtyTypeJSON) AsCtyType() (cty.Type, error) {
	return ctyjson.UnmarshalType(c)
}

func (c CtyTypeJSON) sealed() {}

type CtyValueJSON []byte

type CtyValueMsgpack []byte

func CtyValueAsJSON(v cty.Value, ty cty.Type) ([]byte, error) {
	return ctyjson.Marshal(v, ty)
}

func CtyValueAsMsgpack(v cty.Value, ty cty.Type) ([]byte, error) {
	return ctymsgpack.Marshal(v, ty)
}
