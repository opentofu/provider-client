# OpenTofu Provider Client Library

This is a library for the Go programming language which allows callers to act
as clients to various provider plugin protocols that OpenTofu supports,
published separately from OpenTofu in the hope that it's useful to other
software that wants to act as a client to provider plugins.

> [!CAUTION]
>
> This is currently in very early development and not yet ready to use.
> Anything about this library's API and behavior could potentially change before
> it reaches a stable release.

This is a relatively low-level library that hides some differences between
different protocol variants but does not offer significant abstraction beyond
the direct operations and data types from the protocol. In particular, it does
not include any of the additional behavior OpenTofu implements in terms of
the provider protocol.
