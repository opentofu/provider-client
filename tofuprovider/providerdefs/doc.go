// Package providerdefs contains zero-sized implementations for the various
// interface types in this library that can be embedded into types outside
// of this library to allow external implementations in a forward-compatible
// way.
//
// These types are all either complete or partial implementations of the
// interfaces that will grow to include implementations of any new methods
// added in future which return whatever result is appropriate to represent
// "not implemented" or "not available". For those that are documented as being
// only partial implementations, an external type relying on these will still
// need to provide the remaining mandatory methods itself, but any new
// methods added later will have default implementations that are designed to
// behave similarly to what would happen when using an older provider plugin
// that didn't yet support the new feature.
package providerdefs

import (
	"github.com/opentofu/provider-client/tofuprovider/internal/common"
)

// sealedImpl is an unexported alias of [common.SealedImpl] that is here just
// because otherwise the embedded type in each of our default types tends to
// show up as a code completion suggestion due to being exported as a field
// named "SealedImpl", and that's confusing.
//
// This type is embedded only to get its implementation of the unexported
// [common.Sealed.sealed], so there's no reason for outside callers to access
// it.
type sealedImpl = common.SealedImpl
