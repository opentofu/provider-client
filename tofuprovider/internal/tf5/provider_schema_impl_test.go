package tf5

// NOTE: Unlike most test files in this directory, this one belongs to
// "package tf5" so it can directly test the unexported implementations of the
// various providerschema interface types.

import (
	"maps"
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider/grpc/tfplugin5"
	"github.com/opentofu/provider-client/tofuprovider/internal/mockutil"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

func TestSchemaImpl(t *testing.T) {
	s := providerschema.Schema(schema{
		proto: &tfplugin5.Schema{
			Version: 42,
			Block: &tfplugin5.Schema_Block{
				Version: 24, // This isn't actually exposed anywhere, because it's redundant with the version in the parent
				Attributes: []*tfplugin5.Schema_Attribute{
					{
						Name:     "foo",
						Type:     mockutil.JSON("string"),
						Required: true,
					},
				},
				BlockTypes: []*tfplugin5.Schema_NestedBlock{
					{
						TypeName: "item",
						Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
						Block:    &tfplugin5.Schema_Block{},
					},
				},
				Description:        "Test thingy",
				DescriptionKind:    tfplugin5.StringKind_PLAIN,
				Deprecated:         true,
				DeprecationMessage: "Just for testing.",
			},
		},
	})

	if got, want := s.SchemaVersion(), int64(42); got != want {
		t.Errorf("wrong SchemaVersion\ngot:  %d\nwant: %d", got, want)
	}
	gotDesc, gotDescFormat := s.DocDescription()
	if got, want := gotDesc, "Test thingy"; got != want {
		t.Errorf("wrong description from DocDescription\ngot:  %s\nwant: %s", got, want)
	}
	if got, want := gotDescFormat, providerschema.DocStringPlain; got != want {
		t.Errorf("wrong format from DocDescription\ngot:  %s\nwant: %s", got, want)
	}

	attrs := maps.Collect(s.Attributes())
	if attr, ok := attrs["foo"]; ok {
		if got, want := attr.Usage(), providerschema.AttributeRequired; got != want {
			t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
		}
	} else {
		t.Errorf("missing attribute 'foo'")
	}

	blockTypes := maps.Collect(s.NestedBlockTypes())
	if blockType, ok := blockTypes["item"]; ok {
		if got, want := blockType.Nesting(), providerschema.NestingSingle; got != want {
			t.Errorf("wrong block nesting mode\ngot:  %s\nwant: %s", got, want)
		}
	} else {
		t.Errorf("missing nested block type 'item'")
	}
}

func TestAttributeImpl(t *testing.T) {
	tests := map[string]struct {
		Proto *tfplugin5.Schema_Attribute

		WantType       cty.Type
		WantUsage      providerschema.AttributeUsage
		WantDeprecated bool
		WantSensitive  bool
		WantWriteOnly  bool
		WantDoc        string
		WantDocFormat  providerschema.DocStringFormat
	}{
		// First we'll deal with the various "Usage" combinations
		"required string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:     "foo",
				Type:     mockutil.JSON("string"),
				Required: true,
			},
			WantType:  cty.String,
			WantUsage: providerschema.AttributeRequired,
		},
		"optional string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:     "foo",
				Type:     mockutil.JSON("string"),
				Optional: true,
			},
			WantType:  cty.String,
			WantUsage: providerschema.AttributeOptional,
		},
		"computed string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:     "foo",
				Type:     mockutil.JSON("string"),
				Computed: true,
			},
			WantType:  cty.String,
			WantUsage: providerschema.AttributeComputed,
		},
		"optional+computed string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:     "foo",
				Type:     mockutil.JSON("string"),
				Optional: true,
				Computed: true,
			},
			WantType:  cty.String,
			WantUsage: providerschema.AttributeOptionalComputed,
		},
		"unsupported usage string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:     "foo",
				Type:     mockutil.JSON("string"),
				Required: true,
				Optional: true, // intentionally-invalid combination
			},
			WantType:  cty.String,
			WantUsage: providerschema.AttributeUsageUnsupported,
		},

		// The remaining tests deal with other parts of the attribute schema,
		// aside from the "usage".
		"deprecated bool": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:       "foo",
				Type:       mockutil.JSON("bool"),
				Optional:   true,
				Deprecated: true,
			},
			WantType:       cty.Bool,
			WantDeprecated: true,
			WantUsage:      providerschema.AttributeOptional,
		},
		"sensitive map": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:      "foo",
				Type:      mockutil.JSON([]any{"map", "string"}),
				Optional:  true,
				Sensitive: true,
			},
			WantType:      cty.Map(cty.String),
			WantSensitive: true,
			WantUsage:     providerschema.AttributeOptional,
		},
		"write-only string": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:      "password",
				Type:      mockutil.JSON("string"),
				Optional:  true,
				WriteOnly: true,
			},
			WantType:      cty.String,
			WantWriteOnly: true,
			WantUsage:     providerschema.AttributeOptional,
		},
		"docstring in plain text": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:        "beep",
				Type:        mockutil.JSON("string"),
				Optional:    true,
				Description: "Beep beep!",
			},
			WantType:      cty.String,
			WantUsage:     providerschema.AttributeOptional,
			WantDoc:       "Beep beep!",
			WantDocFormat: providerschema.DocStringPlain,
		},
		"docstring in markdown": {
			Proto: &tfplugin5.Schema_Attribute{
				Name:            "beep",
				Type:            mockutil.JSON("string"),
				Optional:        true,
				Description:     "Beep, and I cannot emphasize this enough, _beep_!",
				DescriptionKind: tfplugin5.StringKind_MARKDOWN,
			},
			WantType:      cty.String,
			WantUsage:     providerschema.AttributeOptional,
			WantDoc:       "Beep, and I cannot emphasize this enough, _beep_!",
			WantDocFormat: providerschema.DocStringMarkdown,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			attr := providerschema.Attribute(attribute{
				proto: test.Proto,
			})

			if test.WantType != cty.NilType {
				if tyC := attr.Type(); tyC != nil {
					ty, err := tyC.AsCtyType()
					if err != nil {
						t.Errorf("type constraint failed to decode: %s", err)
					} else if !ty.Equals(test.WantType) {
						t.Errorf("wrong simple type\ngot:  %#v\nwant: %#v", ty, test.WantType)
					}
				} else {
					t.Errorf("missing expected simple type %#v", test.WantType)
				}
			} else {
				if tyC := attr.Type(); tyC != nil {
					ty, err := tyC.AsCtyType()
					if err == nil {
						t.Errorf("unexpected simple type: %#v", ty)
					} else {
						t.Errorf("unexpected simple type that failed to decode: %s", err)
					}
				}
			}

			if ot := attr.NestedType(); ot != nil {
				// The idea of "nested type" and structural types in general
				// was added in protocol version 6, so this should always
				// be nil in protocol 5.
				t.Errorf("unexpected NestedType: %#v", ot)
			}

			if got, want := attr.Usage(), test.WantUsage; got != want {
				t.Errorf("wrong attribute usage\ngot:  %s\nwant: %s", got, want)
			}

			if got, want := attr.IsDeprecated(), test.WantDeprecated; got != want {
				t.Errorf("wrong IsDeprecated\ngot:  %t\nwant: %t", got, want)
			}

			if got, want := attr.IsSensitive(), test.WantSensitive; got != want {
				t.Errorf("wrong IsSensitive\ngot:  %t\nwant: %t", got, want)
			}

			if got, want := attr.IsWriteOnly(), test.WantWriteOnly; got != want {
				t.Errorf("wrong IsWriteOnly\ngot:  %t\nwant: %t", got, want)
			}

			doc, docFormat := attr.DocDescription()
			if got, want := doc, test.WantDoc; got != want {
				t.Errorf("wrong docstring\ngot:  %s\nwant: %s", got, want)
			}
			if test.WantDoc != "" { // DocFormat is relevant only if doc is set
				if got, want := docFormat, test.WantDocFormat; got != want {
					t.Errorf("wrong docstring format\ngot:  %s\nwant: %s", got, want)
				}
			}
		})
	}
}

func TestNestedBlockTypeImpl(t *testing.T) {
	tests := map[string]struct {
		Proto *tfplugin5.Schema_NestedBlock

		WantNestingMode            providerschema.NestingMode
		WantMinItems, WantMaxItems int64
		WantAttrs                  map[string]*tfplugin5.Schema_Attribute
		WantBlockTypes             map[string]*tfplugin5.Schema_NestedBlock
	}{
		"single empty": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
			},
			WantNestingMode: providerschema.NestingSingle,
		},
		"group empty": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_GROUP,
			},
			WantNestingMode: providerschema.NestingGroup,
		},
		"list empty": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_LIST,
			},
			WantNestingMode: providerschema.NestingList,
		},
		"set empty": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_SET,
			},
			WantNestingMode: providerschema.NestingSet,
		},
		"map empty": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_MAP,
			},
			WantNestingMode: providerschema.NestingMap,
		},
		"item count limits": {
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block:    &tfplugin5.Schema_Block{},
				Nesting:  tfplugin5.Schema_NestedBlock_LIST,
				MinItems: 1,
				MaxItems: 10,
			},
			WantNestingMode: providerschema.NestingList,
			WantMinItems:    1,
			WantMaxItems:    10,
		},
		"attributes": {
			// NOTE: This is focused on testing whether we can enumerate
			// the attributes from the schema, and not on how individual
			// attributes are handled. [TestAttributeImpl] is where we
			// test our translation of [tfplugin5.Schema_Attribute] more
			// thoroughly.
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block: &tfplugin5.Schema_Block{
					Attributes: []*tfplugin5.Schema_Attribute{
						{
							Name:     "a",
							Type:     mockutil.JSON("string"),
							Optional: true,
						},
						{
							Name:     "b",
							Type:     mockutil.JSON("string"),
							Optional: true,
						},
					},
				},
				Nesting: tfplugin5.Schema_NestedBlock_SINGLE,
			},
			WantNestingMode: providerschema.NestingSingle,
			WantAttrs: map[string]*tfplugin5.Schema_Attribute{
				"a": {
					Name:     "a",
					Type:     mockutil.JSON("string"),
					Optional: true,
				},
				"b": {
					Name:     "b",
					Type:     mockutil.JSON("string"),
					Optional: true,
				},
			},
		},
		"nested block types": {
			// NOTE: This is focused on testing whether we can enumerate
			// the nested block types from the schema, and not on how
			// individual block types are handled. [TestNestedBlockTypeImpl] is
			// where we test our translation of [tfplugin5.Schema_NestedBlock]
			// more thoroughly.
			Proto: &tfplugin5.Schema_NestedBlock{
				TypeName: "boop",
				Block: &tfplugin5.Schema_Block{
					BlockTypes: []*tfplugin5.Schema_NestedBlock{
						{
							TypeName: "a",
							Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
							Block:    &tfplugin5.Schema_Block{},
						},
						{
							TypeName: "b",
							Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
							Block:    &tfplugin5.Schema_Block{},
						},
					},
				},
				Nesting: tfplugin5.Schema_NestedBlock_SINGLE,
			},
			WantNestingMode: providerschema.NestingSingle,
			WantBlockTypes: map[string]*tfplugin5.Schema_NestedBlock{
				"a": {
					TypeName: "a",
					Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
					Block:    &tfplugin5.Schema_Block{},
				},
				"b": {
					TypeName: "b",
					Nesting:  tfplugin5.Schema_NestedBlock_SINGLE,
					Block:    &tfplugin5.Schema_Block{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			blockType := providerschema.NestedBlockType(nestedBlockType{
				proto: test.Proto,
			})

			if got, want := blockType.Nesting(), test.WantNestingMode; got != want {
				t.Errorf("wrong nesting mode\ngot:  %s\nwant: %s", got, want)
			}
			minItems, maxItems := blockType.ItemLimits()
			if got, want := minItems, test.WantMinItems; got != want {
				t.Errorf("wrong minimum item count\ngot:  %d\nwant: %d", got, want)
			}
			if got, want := maxItems, test.WantMaxItems; got != want {
				t.Errorf("wrong maximum item count\ngot:  %d\nwant: %d", got, want)
			}

			var gotAttrs map[string]*tfplugin5.Schema_Attribute
			for name, attr := range blockType.Attributes() {
				if gotAttrs == nil {
					gotAttrs = make(map[string]*tfplugin5.Schema_Attribute)
				}
				implAttr, ok := attr.(attribute)
				if !ok {
					t.Errorf("attribute %q has wrong implementation type %T (want %T)", name, attr, implAttr)
					continue
				}
				gotAttrs[name] = implAttr.proto
			}
			if diff := mockutil.Diff(test.WantAttrs, gotAttrs); diff != "" {
				t.Error("wrong attributes\n" + diff)
			}

			var gotBlockTypes map[string]*tfplugin5.Schema_NestedBlock
			for name, blockType := range blockType.NestedBlockTypes() {
				if gotBlockTypes == nil {
					gotBlockTypes = make(map[string]*tfplugin5.Schema_NestedBlock)
				}
				implBlockType, ok := blockType.(nestedBlockType)
				if !ok {
					t.Errorf("nested block type %q has wrong implementation type %T (want %T)", name, blockType, implBlockType)
					continue
				}
				gotBlockTypes[name] = implBlockType.proto
			}
			if diff := mockutil.Diff(test.WantBlockTypes, gotBlockTypes); diff != "" {
				t.Error("wrong nested block types\n" + diff)
			}
		})
	}
}

func TestFunctionSignatureImpl(t *testing.T) {
	tests := map[string]struct {
		Proto *tfplugin5.Function

		WantParameters        []*tfplugin5.Function_Parameter
		WantVariadicParameter *tfplugin5.Function_Parameter
		WantResultType        cty.Type
		WantDoc               string
		WantDocFormat         providerschema.DocStringFormat
		WantDeprecation       string
	}{
		"no params": {
			Proto: &tfplugin5.Function{
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
			},
			WantResultType: cty.String,
		},
		"positional params": {
			Proto: &tfplugin5.Function{
				// NOTE: We're only testing whether we can enumerate the
				// parameters here. The handling of the details of parameters is
				// in the separate test [TestFunctionParameterImpl].
				Parameters: []*tfplugin5.Function_Parameter{
					{
						Name: "foo",
						Type: mockutil.JSON("bool"),
					},
					{
						Name: "bar",
						Type: mockutil.JSON("number"),
					},
				},
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
			},
			WantResultType: cty.String,
			WantParameters: []*tfplugin5.Function_Parameter{
				{
					Name: "foo",
					Type: mockutil.JSON("bool"),
				},
				{
					Name: "bar",
					Type: mockutil.JSON("number"),
				},
			},
		},
		"variadic params": {
			Proto: &tfplugin5.Function{
				// NOTE: We're only testing whether we can detect the
				// parameter here. The handling of the details of parameters is
				// in the separate test [TestFunctionParameterImpl].
				VariadicParameter: &tfplugin5.Function_Parameter{
					Name: "things",
					Type: mockutil.JSON("dynamic"),
				},
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
			},
			WantResultType: cty.String,
			WantVariadicParameter: &tfplugin5.Function_Parameter{
				Name: "things",
				Type: mockutil.JSON("dynamic"),
			},
		},
		"plain-text docstring": {
			Proto: &tfplugin5.Function{
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
				Description:     "Thingy.",
				DescriptionKind: tfplugin5.StringKind_PLAIN,
			},
			WantResultType: cty.String,
			WantDoc:        "Thingy.",
			WantDocFormat:  providerschema.DocStringPlain,
		},
		"markdown docstring": {
			Proto: &tfplugin5.Function{
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
				Description:     "Oh, _such_ thingy.",
				DescriptionKind: tfplugin5.StringKind_MARKDOWN,
			},
			WantResultType: cty.String,
			WantDoc:        "Oh, _such_ thingy.",
			WantDocFormat:  providerschema.DocStringMarkdown,
		},
		"deprecated": {
			Proto: &tfplugin5.Function{
				Return: &tfplugin5.Function_Return{
					Type: mockutil.JSON("string"),
				},
				DeprecationMessage: "Do something else instead.",
			},
			WantResultType:  cty.String,
			WantDeprecation: "Do something else instead.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sig := providerschema.FunctionSignature(functionSignature{
				proto: test.Proto,
			})

			var gotParams []*tfplugin5.Function_Parameter
			for param := range sig.Parameters() {
				implParam, ok := param.(functionParameter)
				if !ok {
					t.Errorf("parameter has wrong implementation type %T (want %T)", param, implParam)
					gotParams = append(gotParams, nil) // just so that any later params are still at the expected indices
					continue
				}
				gotParams = append(gotParams, implParam.proto)
			}
			if diff := mockutil.Diff(test.WantParameters, gotParams); diff != "" {
				t.Error("wrong parameters\n" + diff)
			}

			gotVarParam := sig.VariadicParameter()
			if gotVarParam != nil {
				if wantVarParam := test.WantVariadicParameter; wantVarParam != nil {
					if implParam, ok := gotVarParam.(functionParameter); ok {
						if diff := mockutil.Diff(wantVarParam, implParam.proto); diff != "" {
							t.Error("wrong variadic parameter\n" + diff)
						}
					} else {
						t.Errorf("variadic parameter has wrong implementation type %T (want %T)", gotVarParam, implParam)
					}
				} else {
					t.Errorf("unexpected variadic parameter %#v", gotVarParam)
				}
			} else if test.WantVariadicParameter != nil {
				t.Errorf("missing expected variadic parameter %#v", test.WantVariadicParameter)
			}

			gotResultType := sig.ResultType()
			if gotTy, err := gotResultType.AsCtyType(); err == nil {
				if wantTy := test.WantResultType; !gotTy.Equals(test.WantResultType) {
					t.Errorf("wrong result type\ngot:  %#v\nwant: %#v", gotTy, wantTy)
				}
			} else {
				t.Errorf("invalid result type serialization: %s", err)
			}

			doc, docFormat := sig.DocDescription()
			if got, want := doc, test.WantDoc; got != want {
				t.Errorf("wrong docstring\ngot:  %s\nwant: %s", got, want)
			}
			if test.WantDoc != "" { // DocFormat is relevant only if doc is set
				if got, want := docFormat, test.WantDocFormat; got != want {
					t.Errorf("wrong docstring format\ngot:  %s\nwant: %s", got, want)
				}
			}

			if got, want := sig.DeprecationMessage(), test.WantDeprecation; got != want {
				t.Errorf("wrong deprecation message\ngot:  %s\nwant: %s", got, want)
			}
		})
	}
}

func TestFunctionParameterImpl(t *testing.T) {
	tests := map[string]struct {
		Proto *tfplugin5.Function_Parameter

		WantName                 string
		WantType                 cty.Type
		WantNullValueAllowed     bool
		WantUnknownValuesAllowed bool
		WantDoc                  string
		WantDocFormat            providerschema.DocStringFormat
	}{
		"simple string": {
			Proto: &tfplugin5.Function_Parameter{
				Name: "name",
				Type: mockutil.JSON("string"),
			},
			WantName: "name",
			WantType: cty.String,
		},
		"nullable string": {
			Proto: &tfplugin5.Function_Parameter{
				Name:           "name",
				Type:           mockutil.JSON("string"),
				AllowNullValue: true,
			},
			WantName:             "name",
			WantType:             cty.String,
			WantNullValueAllowed: true,
		},
		"possibly-unknown string": {
			Proto: &tfplugin5.Function_Parameter{
				Name:               "name",
				Type:               mockutil.JSON("string"),
				AllowUnknownValues: true,
			},
			WantName:                 "name",
			WantType:                 cty.String,
			WantUnknownValuesAllowed: true,
		},
		"plain-text docstring": {
			Proto: &tfplugin5.Function_Parameter{
				Name:            "name",
				Type:            mockutil.JSON("string"),
				Description:     "The name of the thing.",
				DescriptionKind: tfplugin5.StringKind_PLAIN,
			},
			WantName:      "name",
			WantType:      cty.String,
			WantDoc:       "The name of the thing.",
			WantDocFormat: providerschema.DocStringPlain,
		},
		"markdown docstring": {
			Proto: &tfplugin5.Function_Parameter{
				Name:            "name",
				Type:            mockutil.JSON("string"),
				Description:     "The name of _The Thing_.",
				DescriptionKind: tfplugin5.StringKind_MARKDOWN,
			},
			WantName:      "name",
			WantType:      cty.String,
			WantDoc:       "The name of _The Thing_.",
			WantDocFormat: providerschema.DocStringMarkdown,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			param := providerschema.FunctionParameter(functionParameter{
				proto: test.Proto,
			})

			if got, want := param.Name(), test.WantName; got != want {
				t.Errorf("wrong name\ngot:  %s\nwant: %s", got, want)
			}

			gotType := param.Type()
			if gotTy, err := gotType.AsCtyType(); err == nil {
				if wantTy := test.WantType; !gotTy.Equals(test.WantType) {
					t.Errorf("wrong result type\ngot:  %#v\nwant: %#v", gotTy, wantTy)
				}
			} else {
				t.Errorf("invalid result type serialization: %s", err)
			}

			if got, want := param.NullValueAllowed(), test.WantNullValueAllowed; got != want {
				t.Errorf("wrong NullValueAllowed\ngot:  %t\nwant: %t", got, want)
			}

			if got, want := param.UnknownValuesAllowed(), test.WantUnknownValuesAllowed; got != want {
				t.Errorf("wrong UnknownValuesAllowed\ngot:  %t\nwant: %t", got, want)
			}

			doc, docFormat := param.DocDescription()
			if got, want := doc, test.WantDoc; got != want {
				t.Errorf("wrong docstring\ngot:  %s\nwant: %s", got, want)
			}
			if test.WantDoc != "" { // DocFormat is relevant only if doc is set
				if got, want := docFormat, test.WantDocFormat; got != want {
					t.Errorf("wrong docstring format\ngot:  %s\nwant: %s", got, want)
				}
			}
		})
	}
}
