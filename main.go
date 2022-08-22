package main

import (
	"context"
	"terraform-provider-azurermagw/azurermagw"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {	
	tfsdk.Serve(context.Background(), azurermagw.New, tfsdk.ServeOpts{
		Name: "azurermagw",	})
}
