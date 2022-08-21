package main

import (
	"context"
	"terraform-provider-azurermagw/azurermagw"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)


func main() {	
	tfsdk.Serve(context.Background(), azurermagw.New, tfsdk.ServeOpts{
		Name: "azurermagw",	})
}
