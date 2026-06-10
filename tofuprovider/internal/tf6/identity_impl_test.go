package tf6

import (
	"slices"
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin6"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func TestIdentitySchemaImpl(t *testing.T) {
	s := providerschema.IdentitySchema(identitySchema{
		proto: &tfplugin6.ResourceIdentitySchema{
			Version: 42,
			IdentityAttributes: []*tfplugin6.ResourceIdentitySchema_IdentityAttribute{
				{
					Name:              "foo",
					Type:              mockutil.JSON("string"),
					RequiredForImport: true,
					OptionalForImport: true,
					Description:       "foo descripion",
				},
				{
					Name:              "bar",
					Type:              mockutil.JSON("number"),
					RequiredForImport: false,
					OptionalForImport: false,
					Description:       "bar description",
				},
			},
		},
	})

	if got, want := s.Version(), int64(42); got != want {
		t.Errorf("wrong Version\ngot:  %d\nwant: %d", got, want)
	}

	attrs := slices.Collect(s.IdentityAttributes())
	if len(attrs) != 2 {
		t.Errorf("wrong number of attributes\ngot:  %d\nwant: %d", len(attrs), 2)
	}

	{
		fooAttr := attrs[0]
		if got, want := fooAttr.Name(), "foo"; got != want {
			t.Errorf("wrong attribute name\ngot:  %q\nwant: %q", got, want)
		}
		if !fooAttr.RequiredForImport() {
			t.Errorf("wrong RequiredForImport\ngot:  %t\nwant: %t", false, true)
		}
		if !fooAttr.OptionalForImport() {
			t.Errorf("wrong OptionalForImport\ngot:  %t\nwant: %t", false, true)
		}
		if got, want := fooAttr.Description(), "foo descripion"; got != want {
			t.Errorf("wrong description\ngot:  %q\nwant: %q", got, want)
		}
		ty, err := fooAttr.Type().AsCtyType()
		if err != nil {
			t.Fatalf("invalid attribute type: %s", err)
		}
		if got, want := ty, cty.String; got != want {
			t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
		}
	}

	{
		barAttr := attrs[1]
		if got, want := barAttr.Name(), "bar"; got != want {
			t.Errorf("wrong attribute name\ngot:  %q\nwant: %q", got, want)
		}
		if barAttr.RequiredForImport() {
			t.Errorf("wrong RequiredForImport\ngot:  %t\nwant: %t", true, false)
		}
		if barAttr.OptionalForImport() {
			t.Errorf("wrong OptionalForImport\ngot:  %t\nwant: %t", true, false)
		}
		if got, want := barAttr.Description(), "bar description"; got != want {
			t.Errorf("wrong description\ngot:  %q\nwant: %q", got, want)
		}
		ty, err := barAttr.Type().AsCtyType()
		if err != nil {
			t.Fatalf("invalid attribute type: %s", err)
		}
		if got, want := ty, cty.Number; got != want {
			t.Errorf("wrong attribute type\ngot:  %#v\nwant: %#v", got, want)
		}
	}

}

func TestResourceIdentityDataImpl(t *testing.T) {
	s := providerschema.ResourceIdentityData(resourceIdentityData{
		proto: &tfplugin6.ResourceIdentityData{
			IdentityData: &tfplugin6.DynamicValue{
				Json: mockutil.JSON(map[string]any{
					"foo": "bar",
					"bar": int64(42),
				}),
			},
		},
	})

	wantType := cty.Object(map[string]cty.Type{
		"foo": cty.String,
		"bar": cty.Number,
	})
	ctyVal, err := s.IdentityData().AsCtyValue(wantType)
	if err != nil {
		t.Fatalf("invalid identity data: %s", err)
	}

	if got, want := ctyVal.GetAttr("foo").AsString(), "bar"; got != want {
		t.Errorf("wrong attribute value\ngot:  %q\nwant: %q", got, want)
	}

	if got, _ := ctyVal.GetAttr("bar").AsBigFloat().Int64(); got != 42 {
		t.Errorf("wrong attribute value\ngot:  %d\nwant: %d", got, int64(42))
	}
}
