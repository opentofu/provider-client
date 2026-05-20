package mockutil

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
)

// ComparableFunctionError is an implementation of [providerops.FunctionError]
// that is friendly to being compared using functions from [cmp].
//
// If you use [Eq] (or other similar [gomock.Matcher] implementations based on
// [Comparer]) or compare using [Diff], then any implementations of
// [providerops.FunctionError] are automatically converted to this type so that
// errors are compared by their user-facing content rather than by the types
// used to implement them.
//
// This type is exported as a convenient way to describe "expected error"
// in a test, using composite literal syntax.
type ComparableFunctionError struct {
	Text_          string
	ArgumentIndex_ *int

	common.SealedImpl
}

var _ providerops.FunctionError = (*ComparableFunctionError)(nil)

// Text implements [providerops.FunctionError].
func (c *ComparableFunctionError) Text() string {
	return c.Text_
}

// ArgumentIndex implements [providerops.FunctionError].
func (c *ComparableFunctionError) ArgumentIndex() (int, bool) {
	if c.ArgumentIndex_ == nil {
		return 0, false
	}
	return *c.ArgumentIndex_, true
}

func cmpTransformFunctionError(diag providerops.FunctionError) *ComparableFunctionError {
	ret := &ComparableFunctionError{
		Text_: diag.Text(),
	}
	if argIdx, ok := diag.ArgumentIndex(); ok {
		ret.ArgumentIndex_ = new(argIdx)
	}
	return ret
}
