# Design Notes About This Library

The following are some general notes about the design goals and decisions that
informed the API and implementation of this library.

When interacting with OpenTofu providers there are a number of different layers
of abstraction to think about, including but probably not limited to:

1. The mechanical process of launching a provider plugin and the raw wire
   protocol used to interact with it.

     At this level, the raw protocol details are exposed and so a separate
     implementation is needed for each variation of the protocol.
2. The conceptual operations that the provider protocol provides, regardless
   of the specific execution process and wire protocol used to achieve those
   operations.

     At this level a caller does need to be aware of the low-level conceptual
     details of the protocol, such as what order operations should be called in
     to achieve a given effect, but they can (at least to some extent) share
     a single caller implementation across a number of different wire
     protocols and execution strategies.
3. The higher-level operations that OpenTofu implements in terms of the
   provider protocol concepts, such as the validate, plan, apply sequence
   typically used to make changes to managed resource instances.

     This layer includes all of OpenTofu's opinions about how these different
     steps ought to interact, such as the consistency rules providers are
     required to follow as objects pass through the series of operations and
     get gradually closer to being the "new state".

This library is aimed at abstraction level 2 from the above list: it aims to
hide the details of executing a provider plugin and the wire protocol used to
send requests to it, but it has very little abstraction beyond that. Callers
are expected to be aware of the provider protocol's usage rules and to call
an appropriate sequence of methods to properly meet a provider's expectations
as guaranteed by the protocol.

Along with the primary goal of abstracting away the differences between
wire protocols and execution models, we have a secondary goal of introducing
as little translation overhead as possible: a caller should only pay a
translation cost for operations and data they actually use. The API of this
library therefore relies a lot on dynamic dispatch and Go iterators to delay
the conversion of protocol-version-specific data types into this library's
data types until the very last moment, which makes some sacrifices on
ergonomics but hopefully strikes a good balance nonetheless.

## Naming Conventions

As the lowest-level wire protocols have evolved, the terminology used in them
has shifted depending on the fashion of the day. There is therefore no single
set of terminology that all protocol versions have in common.

Because this library is focused on the client side of the provider protocol,
and because OpenTofu itself is the primary and most important client of the
protocol, this library chooses to follow the naming conventions most commonly
used in the OpenTofu Core codebase, even if the current wire protocols all use
different terminology.

This is admittedly a rather arbitrary decision, but its goal is to prioritize
this library fitting well as a callee of the OpenTofu Core codebase, and
secondarily so that those who refer to code in the OpenTofu Core codebase can
hopefully more easily understand how the core functionality corresponds to the
provider protocol functionality when reimplementing similar behaviors in
their own client code.

## Types and Values

OpenTofu has its own type system which forms a core part of the protocol in
representing any data whose type is decided dynamically by the provider rather
than fixed as part of the protocol itself.

To avoid re-modelling that entire type system again in this library, we use
the same `go-cty` library that OpenTofu itself uses to implement the foundations
of its type system. In particular, the API of this library uses `cty.Type` and
`cty.Value` as the primary abstraction for representing types and values,
translating those to and from the serialization formats used by the underlying
wire protocols internally as necessary.

However, this library only uses `cty`'s type and value representations and its
serializer/deserializer functions, and does not perform any `cty`-level value
operations such as arithmetic or type conversions; callers are expected to
perform any semantic transformations required by the protocol themselves, so
that the given values can be sent verbatim to provider operations.
