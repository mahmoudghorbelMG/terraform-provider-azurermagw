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

func createBackendAddressPool(backend_plan Backend_address_pool) BackendAddressPool{
	backend_json := BackendAddressPool{
		Name: backend_plan.Name.Value,
		Properties: struct {
			ProvisioningState string "json:\"provisioningState,omitempty\""
			BackendAddresses  []struct {
				Fqdn      string "json:\"fqdn,omitempty\""
				IPAddress string "json:\"ipAddress,omitempty\""
			} "json:\"backendAddresses\""
			RequestRoutingRules []struct {
				ID string "json:\"id,omitempty\""
			} "json:\"requestRoutingRules,omitempty\""
		}{},
		Type: "Microsoft.Network/applicationGateways/backendAddressPools",
	}
	length := len(backend_plan.Fqdns) + len(backend_plan.Ip_addresses)

	//If there is no fqdn nor IPaddress for the backend pool, initialize the BackendAddresses to nil to avoid a terraform provider error when making the state
	if length == 0 {
		backend_json.Properties.BackendAddresses = nil
	} else {
		backend_json.Properties.BackendAddresses = make([]struct {
			Fqdn      string "json:\"fqdn,omitempty\""
			IPAddress string "json:\"ipAddress,omitempty\""
		}, length)
	}
	for i := 0; i < len(backend_plan.Fqdns); i++ {
		backend_json.Properties.BackendAddresses[i].Fqdn = backend_plan.Fqdns[i].Value
	}
	for i := 0; i < len(backend_plan.Ip_addresses); i++ {
		backend_json.Properties.BackendAddresses[i+len(backend_plan.Fqdns)].IPAddress = backend_plan.Ip_addresses[i].Value
	}
	
	return backend_json
}
// we have to give the nb_Fqdns and nb_IpAddress in order to make this function reusable in create, read and update method
func generateBackendAddressPoolState(gw ApplicationGateway, backendAddressPoolName string,nb_Fqdns int,nb_IpAddress int) Backend_address_pool {
	//retrieve json element from gw	
	index := getBackendAddressPoolElementKey(gw, backendAddressPoolName)
	backend_json := gw.Properties.BackendAddressPools[index]

	
	// Map response body to resource schema attribute
	backend_state := Backend_address_pool{
		Name:         types.String{Value: backend_json.Name},
		Id:           types.String{Value: backend_json.ID},
		Fqdns:        []types.String{},
		Ip_addresses: []types.String{},
	}
	
	
	if nb_Fqdns != 0 {
		backend_state.Fqdns = make([]types.String, nb_Fqdns)
	} else {
		backend_state.Fqdns = nil
	}
	
	if nb_IpAddress != 0 {
		backend_state.Ip_addresses = make([]types.String, nb_IpAddress)
	} else {
		backend_state.Ip_addresses = nil
	}
	
	index_nb_Fqdns		:=0
	index_nb_IpAddress 	:= 0
	for i := 0; i < nb_Fqdns + nb_IpAddress; i++ {
		if backend_json.Properties.BackendAddresses[i].Fqdn != "" {
			backend_state.Fqdns[index_nb_Fqdns] = types.String{Value: backend_json.Properties.BackendAddresses[i].Fqdn}
			index_nb_Fqdns++
		}
		if backend_json.Properties.BackendAddresses[i].IPAddress != "" {
			backend_state.Ip_addresses[index_nb_IpAddress] = types.String{Value: backend_json.Properties.BackendAddresses[i].IPAddress}
			index_nb_IpAddress++
		}
	}

	return backend_state
}
func getBackendAddressPoolElementKey(gw ApplicationGateway, backendAddressPoolName string) int {
	key := -1
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			key = i
		}
	}
	return key
}
func checkBackendAddressPoolElement(gw ApplicationGateway, backendAddressPoolName string) bool {
	exist := false
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			exist = true
		}
	}
	//fmt.Println("ww         Exist =",exist)
	return exist
}
func removeBackendAddressPoolElement(gw *ApplicationGateway, backendAddressPoolName string) {
	//removed := false
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools[:i], gw.Properties.BackendAddressPools[i+1:]...)
			//removed = true
		}
	}	
}