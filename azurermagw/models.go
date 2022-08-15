package azurermagw

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebappBinding
type WebappBinding struct {
	Name                 	types.String         	`tfsdk:"name"`
	Agw_name             	types.String         	`tfsdk:"agw_name"`
	Agw_rg               	types.String         	`tfsdk:"agw_rg"`
	Backend_address_pool	Backend_address_pool 	`tfsdk:"backend_address_pool"`
	Backend_http_settings   Backend_http_settings	`tfsdk:"backend_http_settings"`
	Probe					Probe_tf				`tfsdk:"probe"`
	Http_listener			*Http_listener			`tfsdk:"http_listener"`
	Https_listener			Http_listener			`tfsdk:"https_listener"`
}
