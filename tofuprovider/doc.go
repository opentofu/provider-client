// Package tofuprovider is a low-level client library for the OpenTofu provider
// plugin API, allowing Go programs to call into OpenTofu provider plugins
// without using code from OpenTofu itself.
//
// The scope of this library is intentionally limited: it focuses only on
// hiding low-level wire protocol details and differences between supported
// protocols, while otherwise directly modeling the conceptual protocol with
// minimal additional abstraction.
package tofuprovider
