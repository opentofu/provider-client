package common

// Sealed is an interface that we embed in various other interface types to
// ensure that they can't be directly implemented by types outside of this
// module, because we need to be able to grow those interfaces over time as the
// underlying protocols evolve.
//
// This module uses interfaces mainly for dynamic dispatch over a fixed
// set of implementations covering different versions of the plugin protocol,
// and not to support third-party implementations. If you need to write a
// third-party implementation e.g. for mocking in a calling application, you
// can embed the zero-size base types from
// [github.com/opentofu/provider-client/tofuproviders/providerdefs] to ensure
// that your subtypes will automatically get default implementations of any
// new methods added in later versions of this library.
type Sealed interface {
	// This unexported method can be implemented only by types in this
	// package. Embed [SealedImpl] into another struct type to make it
	// implement this interface.
	sealed()
}

// SealedImpl is a zero-sized type that implements [Sealed], intended to
// be embedded into other struct types to make them implement that interface
// even though they can't directly implement the "sealed()" method.
//
// If you've found this while attempting to implement an interface from this
// library in an external Go module: the supported way to achieve that is
// to embed the appropriate base type from
// [github.com/opentofu/provider-client/tofuproviders/providerdefs] into your
// implementation type, which will then implement [Sealed] while also providing
// default implementations for any methods added in later versions.
type SealedImpl struct{}

// sealed implements [Sealed].
func (s SealedImpl) sealed() {}
