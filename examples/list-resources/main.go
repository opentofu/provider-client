package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/provider-client/tofuprovider"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

const resourceTypeName = "aws_s3_bucket"

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <provider-executable> [provider-args...]\n", args[0])
		os.Exit(1)
	}
	args = args[1:]

	ctx := context.Background()
	provider, err := tofuprovider.StartGRPCPlugin(ctx, args[0], args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer provider.Close()

	resp, err := provider.GetProviderSchema(ctx, &providerops.GetProviderSchemaRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Provider schema request failed: %s\n", err)
		os.Exit(1)
	}

	var listSchema providerschema.Schema
	for name, v := range resp.ProviderSchema().ManagedResourceTypeListSchemas() {
		if name == resourceTypeName {
			listSchema = v
			break
		}
	}
	if listSchema == nil {
		fmt.Fprintf(os.Stderr, "Error: provider has no %s list resource schema\n", resourceTypeName)
		os.Exit(1)
	}

	var resourceSchema providerschema.Schema
	for name, v := range resp.ProviderSchema().ManagedResourceTypeSchemas() {
		if name == resourceTypeName {
			resourceSchema = v
			break
		}
	}

	resourceType, err := impliedType(resourceSchema)
	if err != nil {
		log.Fatal(err)
	}

	providerConfigType, err := impliedType(resp.ProviderSchema().ProviderConfigSchema())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot derive provider config type: %s\n", err)
		os.Exit(1)
	}
	providerConfig := providerschema.NewDynamicValue(cty.NullVal(providerConfigType), providerConfigType)

	configureResp, err := provider.ConfigureProvider(ctx, &providerops.ConfigureProviderRequest{
		TerraformVersion: "1.15.5",
		Config:           providerConfig,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if configureResp.Diagnostics().HasErrors() {
		for d := range configureResp.Diagnostics().All() {
			fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", d.Summary(), d.Detail())
		}
		os.Exit(1)
	}

	fmt.Println("Provider configuration done")

	listConfigType, err := impliedType(listSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot derive list config type: %s\n", err)
		os.Exit(1)
	}
	listConfigValue := nullObjectWith(listConfigType, map[string]cty.Value{
		"region": cty.StringVal("eu-central-1"),
	})

	fmt.Printf("Listing %s in region %s:\n", resourceTypeName, listConfigValue.GetAttr("region").AsString())

	list, err := provider.ListManagedResources(ctx, &providerops.ListManagedResourcesRequest{
		TypeName: resourceTypeName,
		Config:   providerschema.NewDynamicValue(listConfigValue, listConfigType),
		Limit:    10000,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	defer list.Close(ctx)

	for {
		res, err := list.ReadResult(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		if res.Diagnostics().HasErrors() {
			for d := range res.Diagnostics().All() {
				fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", d.Summary(), d.Detail())
			}
			os.Exit(1)
		}

		fmt.Printf("- display name: %s\n", res.DisplayName())

		if r := res.Resource(); r != nil {
			resVal, err := res.Resource().AsCtyValue(resourceType)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(" %+v\n", resVal)
		}
	}
}

func impliedType(block providerschema.BlockType) (cty.Type, error) {
	attrTypes := map[string]cty.Type{}

	for name, attr := range block.Attributes() {
		ty, err := attributeType(attr)
		if err != nil {
			return cty.NilType, fmt.Errorf("attribute %q: %w", name, err)
		}
		attrTypes[name] = ty
	}

	for name, nested := range block.NestedBlockTypes() {
		nestedTy, err := impliedType(nested)
		if err != nil {
			return cty.NilType, fmt.Errorf("block %q: %w", name, err)
		}
		switch nested.Nesting() {
		case providerschema.NestingSingle, providerschema.NestingGroup:
			attrTypes[name] = nestedTy
		case providerschema.NestingList:
			attrTypes[name] = cty.List(nestedTy)
		case providerschema.NestingSet:
			attrTypes[name] = cty.Set(nestedTy)
		case providerschema.NestingMap:
			attrTypes[name] = cty.Map(nestedTy)
		default:
			return cty.NilType, fmt.Errorf("block %q: unsupported nesting mode %s", name, nested.Nesting())
		}
	}

	return cty.Object(attrTypes), nil
}

func attributeType(attr providerschema.Attribute) (cty.Type, error) {
	if nested := attr.NestedType(); nested != nil {
		attrTypes := map[string]cty.Type{}
		for name, a := range nested.Attributes() {
			ty, err := attributeType(a)
			if err != nil {
				return cty.NilType, fmt.Errorf("attribute %q: %w", name, err)
			}
			attrTypes[name] = ty
		}
		obj := cty.Object(attrTypes)
		switch nested.Nesting() {
		case providerschema.NestingSingle, providerschema.NestingGroup:
			return obj, nil
		case providerschema.NestingList:
			return cty.List(obj), nil
		case providerschema.NestingSet:
			return cty.Set(obj), nil
		case providerschema.NestingMap:
			return cty.Map(obj), nil
		default:
			return cty.NilType, fmt.Errorf("unsupported nesting mode %s", nested.Nesting())
		}
	}

	return attr.Type().AsCtyType()
}

func nullObjectWith(ty cty.Type, overrides map[string]cty.Value) cty.Value {
	vals := map[string]cty.Value{}
	for name, attrTy := range ty.AttributeTypes() {
		if v, ok := overrides[name]; ok {
			vals[name] = v
		} else {
			vals[name] = cty.NullVal(attrTy)
		}
	}
	return cty.ObjectVal(vals)
}
