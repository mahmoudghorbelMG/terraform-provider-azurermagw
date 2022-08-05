package azurermagw

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BackendAddressPool struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState string `json:"provisioningState,omitempty"`
		BackendAddresses  []struct {
			Fqdn      string `json:"fqdn,omitempty"`
			IPAddress string `json:"ipAddress,omitempty"`
		} `json:"backendAddresses"`
		RequestRoutingRules []struct {
			ID string `json:"id,omitempty"`
		} `json:"requestRoutingRules,omitempty"`
	} `json:"properties"`
	Type string `json:"type,omitempty"`
}
type Backend_address_pool struct {
	Name         types.String   `tfsdk:"name"`
	Id           types.String   `tfsdk:"id"`
	Fqdns        []types.String `tfsdk:"fqdns"`
	Ip_addresses []types.String `tfsdk:"ip_addresses"`
}
