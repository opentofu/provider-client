package providerops

import (
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"
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
	// The following matches all of the error values we expect that our own
	// tofuprovider.Provider implementations could return.
	switch {
	case grpcStatus.Code(err) == grpcCodes.Unimplemented:
		return true
	case err == error(common.ErrUnimplemented):
		return true
	default:
		return false
	}
}
