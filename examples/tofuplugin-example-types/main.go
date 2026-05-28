package main

import (
	"context"
	"fmt"
	"maps"
	"os"

	"github.com/opentofu/provider-client/tofuprovider"
	"github.com/opentofu/provider-client/tofuprovider/providerops"
	"github.com/opentofu/provider-client/tofuprovider/providerschema"
)

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

	schema := resp.ProviderSchema()
	managedResourceTypes := maps.Collect(schema.ManagedResourceTypeSchemas())
	dataResourceTypes := maps.Collect(schema.DataResourceTypeSchemas())
	ephemeralResourceTypes := maps.Collect(schema.EphemeralResourceTypeSchemas())
	functions := maps.Collect(schema.FunctionSignatures())

	identityResp, err := provider.GetIdentitySchemas(ctx, &providerops.GetIdentitySchemasRequest{})
	var identitySchemas map[string]providerschema.IdentitySchema
	switch {
	case providerops.IsUnimplementedErr(err):
		// Older providers don't implement the resource identity RPC.
	case err != nil:
		fmt.Fprintf(os.Stderr, "Identity schemas request failed: %s\n", err)
		os.Exit(1)
	default:
		identitySchemas = maps.Collect(identityResp.IdentitySchemas())
	}

	printList("Managed Resource Types", managedResourceTypes)
	printList("Managed Resource Types with Identity", identitySchemas)
	printList("Data Resource Types", dataResourceTypes)
	printList("Ephemeral Resource Types", ephemeralResourceTypes)
	printList("Functions", functions)
	fmt.Printf("\n")
}

func printList[V any](title string, m map[string]V) {
	if len(m) == 0 {
		return
	}

	fmt.Printf("\n# %s\n\n", title)
	for name := range m {
		fmt.Printf("- `%s`\n", name)
	}
}
