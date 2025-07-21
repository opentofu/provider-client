package providerops

import (
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// IsUnimplementedErr returns true if the given error represents "operation not
// implemented".
//
// Callers might use this to trigger fallback behavior using an older protocol
// feature that the provider might implement instead.
//
// It's only meaningful to call this with errors returned by the methods of
// [tofuprovider.Provider]. Errors obtained from other locations produce
// unspecified results.
func IsUnimplementedErr(err error) bool {
	// We'll try to see if this is a gRPC "unimplemented" error. This
	// function returns codes.Unimplemented only if the given error
	// is based on a gRPC status with that code.
	switch {
	case grpcStatus.Code(err) == grpcCodes.Unimplemented:
		return true

		// (if we have protocol implementations that are not gRPC-based in
		// future then we should add additional cases to catch whatever
		// they return to represent "unimplemented" here.)

	default:
		return false
	}
}
