package common

// UnimplementedError is one of the [error] types that we recognize as
// representing "not implemented" in providerops.IsUnimplementedErr, though
// we use it mainly for the default implementations in package providerdefs
// since other provider implementations can return error types specific to
// the protocol used for communication with the provider, like gRPC's own
// "Unimplemented" error code.
//
// [ErrUnimplemented] is the only value of this type.
type UnimplementedError struct{}

func (UnimplementedError) Error() string {
	return "not implemented"
}

// ErrUnimplemented is the only value of type [UnimplementedError].
var ErrUnimplemented UnimplementedError
