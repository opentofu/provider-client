package common

// Sealed is an interface that we embed in various other interface types to
// ensure that they can't be implemented by types outside of this module,
// because we need to be able to grow those interfaces over time as the
// underlying protocols evolve.
//
// This module uses interfaces only for dynamic dispatch over a fixed
// set of implementations covering different versions of the plugin protocol,
// and not to support third-party implementations.
type Sealed interface {
	// This unexported method can be implemented only by types in this
	// package. Embed [SealedImpl] into another struct type to make it
	// implement this interface.
	sealed()
}

// SealedImpl is a zero-sized type that implements [Sealed], intended to
// be embedded into other struct types to make them implement that interface
// even though they can't directly implement the "sealed()" method.
type SealedImpl struct{}

// sealed implements [Sealed].
func (s SealedImpl) sealed() {}
