// Package tofuprovider is a low-level client library for the OpenTofu provider
// plugin API, allowing Go programs to call into OpenTofu provider plugins
// without using code from OpenTofu itself.
//
// The scope of this library is intentionally limited: it focuses only on
// abstracting away the protocol major version negotiation so that caller
// code does not need to be duplicated to support multiple protocol versions.
// Otherwise, the API closely matches how the underlying protocol is designed
// with as little additional abstraction as possible.
//
// This package currently implements clients for protocol major versions 5 and
// 6.
package tofuprovider
