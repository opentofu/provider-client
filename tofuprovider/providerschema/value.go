package providerschema

import (
	"fmt"

	"github.com/opentofu/provider-client/tofuprovider/internal/common"

	"github.com/zclconf/go-cty/cty"
)

// DynamicValueOut is a dynamically-typed value returned in a provider response,
// in an opaque form that requires schema information to decode.
//
// This indirection exists both to allow callers to avoid paying an
// unmarshalling cost for values they don't use and because the protocol
// expects the caller to "know" which schema to use when decoding a value,
// based on broader context.
type DynamicValueOut interface {
	// AsCtyValue attempts to interpret the data as a value of the given
	// type, returning the result if successful.
	//
	// Some implementations delay parsing the protocol's wire format until
	// this method is called, and so an error for this function could represent
	// either that the provider returned an incorrectly-serialized value or
	// that the value is incompatible with the given type.
	//
	// Explicit type information is needed because the wire formats used in
	// the provider protocol expect the caller to have retrieved schema
	// information out of band in order to avoid redundantly retransmitting
	// type information alongside every value. The correct type to pass
	// depends on the context in which the value was returned.
	AsCtyValue(withType cty.Type) (cty.Value, error)

	// This interface cannot be implemented outside of this module, because
	// future versions might extend the interface to include new protocol
	// features.
	common.Sealed
}

// DynamicValueIn is the equivalent of [DynamicValueOut] for values being
// sent _to_ a provider request by a caller.
//
// [DynamicValueOut] represents a value of a known wire format but of unknown
// serialization type, while [DynamicValueIn] represents
// a value of a known serialization type but as-yet-undecided wire format.
//
// Use [NewDynamicValue] to construct values of this type. The zero value
// of this type is [NoDynamicValue], representing the absense of a value.
type DynamicValueIn struct {
	// v is the value to be serialized
	v cty.Value

	// ty is the type constraint used to serialize it.
	ty cty.Type
}

var NoDynamicValue DynamicValueIn

// NewDynamicValue constructs a [DynamicValueIn] with the given value and
// serialization type.
//
// The "serialization type" of a value is separate from the value's own
// dynamic type because serialization formats use it to decide whether
// they need to serialize type information in-band in the value or whether
// the serialization type is sufficient. The correct type to use depends on
// where the value is being sent, but it's often a type implied by the
// schema of a resource type, where the schema might include attributes whose
// types are decided dynamically at runtime rather than fixed in the schema.
func NewDynamicValue(v cty.Value, ty cty.Type) DynamicValueIn {
	if v == cty.NilVal {
		// The total absense of a dynamic value value should be represented as
		// [NoDynamicValue], which is the zero value of [DynamicValueIn].
		//
		// Note that this is different from representing a _null_ value. The
		// zero value of [DynamicValueIn] represents the value being absent
		// from the perspective of Go code and the protocol, whereas null
		// represents the absense of a value at the OpenTofu language level.
		// Yes, this is quite confusing, but is unfortunately how the system
		// already works. Sorry.
		panic("cannot use cty.NilVal as dynamic value")
	}
	if ty == cty.NilType {
		panic("cannot use cty.NilType for dynamic value")
	}
	return DynamicValueIn{
		v:  v,
		ty: ty,
	}
}

// Value returns the value to be serialized, or [cty.NilVal] if called on
// [NoDynamicValue].
func (dv DynamicValueIn) Value() cty.Value {
	return dv.v
}

// SerializationType returns the type that the value should be serialized as,
// or [cty.NilType] if called on [NoDynamicValue].
func (dv DynamicValueIn) SerializationType() cty.Type {
	return dv.ty
}

// RawState is a low-level, raw representation of resource instance state
// as it would be saved by OpenTofu in a state snapshot.
//
// This is used only when preparing previously-saved state data for use
// in a new version, where the client cannot know what schema was used
// in an earlier version of the provider and thus how to decode the
// data saved in the state.
type RawState struct {
	// JSON is the JSON representation of the resource instance state data
	// exactly as saved in OpenTofu state snapshot files.
	//
	// This part of the protocol is troublesome for use by clients other than
	// OpenTofu because it relies on internal details of OpenTofu's state
	// snapshot format. Use NewRawState with data state data returned from
	// some other provider operation to produce a suitable value to save
	// and provide in a future call to UpgradeeManagedResourceState.
	JSON []byte

	// Flatmap is a historical legacy representation of resource instance
	// state data created by obsolete versions of Terraform. Clients should
	// use this only if they are trying to work with data from
	// OpenTofu/Terraform state snapshots and when the corresponding field
	// is populated in those state snapshots. None of the provider protocols
	// supported by this library can produce new data in flatmap format,
	// so clients that are only saving data they created using other parts
	// of this library can ignore this field completely.
	Flatmap map[string]string
}

// NewRawState emulates OpenTofu's behavior when preparing resource instance
// state data to save as part of a state snapshot file, producing some data
// that can be saved and used in future calls to UpgradeeManagedResourceState
// to prepare a value for use in possibly-different version of the same
// provider.
//
// Because [RawState.Flatmap] is a legacy field only used for state snapshots
// created in obsolete versions of Terraform, that field is guaranteed to
// always be nil in the result: this function always populates only the
// JSON field, and so it's valid for a caller to just save the content of
// that field alone.
//
// This is a small pragmatic exception to the usual rule that this library
// does not provide any abstraction of the provider protocol other than
// hiding the protocol version negotiation, since this behavior is relatively
// straightforward to implement and it's impossible to implement the resource
// instance change lifecycle correctly without it.
func NewRawState(v cty.Value, ty cty.Type) (RawState, error) {
	// OpenTofu's current state snapshot representation at the time of writing
	// is to use cty's JSON serialization of the value in the JSON field,
	// and leave Flatmap completely unset.
	src, err := common.CtyValueAsJSON(v, ty)
	if err != nil {
		return RawState{}, fmt.Errorf("can't serialize as JSON: %w", err)
	}
	return RawState{
		JSON: src,
	}, nil
}
