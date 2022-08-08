package azurermagw

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BackendHTTPSettings struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		//I will focus on the params that we use
		AffinityCookieName             string      `json:"affinityCookieName,omitempty"`
		CookieBasedAffinity            string      `json:"cookieBasedAffinity,omitempty"`
		PickHostNameFromBackendAddress bool        `json:"pickHostNameFromBackendAddress,omitempty"`
		Port           				   int    	   `json:"port,omitempty"`
		Probe          *struct {
			ID string `json:"id,omitempty"`
		} `json:"probe"`
		Protocol                       string      `json:"protocol,omitempty"`
		RequestTimeout int    `json:"requestTimeout,omitempty"`

		//we don't use these ones. They will not appear in this current implementation
		//I there is someone interested in these params, (s)he should implement them in accordance with Backend_http_settings struct
		AuthenticationCertificates     *[]struct { ////ajouté
			ID string `json:"id"`
		} `json:"authenticationCertificates"`
		ConnectionDraining *struct { ////ajouté
			DrainTimeoutInSec int  `json:"drainTimeoutInSec,omitempty"`
			Enabled           bool `json:"enabled,omitempty"`
		} `json:"connectionDraining"`
		HostName       string `json:"hostName,omitempty"` ////ajouté
		Path           string `json:"path,omitempty"`
		ProbeEnabled        bool `json:"probeEnabled,omitempty"` ////ajouté
		ProvisioningState              string      `json:"provisioningState,omitempty"`
		RequestRoutingRules *[]struct {
			ID string `json:"id,omitempty"`
		} `json:"requestRoutingRules"`
		TrustedRootCertificates *[]struct { ////ajouté
			ID string `json:"id,omitempty"`
		} `json:"trustedRootCertificates"`
	} `json:"properties"`
	Type string `json:"type,omitempty"`
} 

type Backend_http_settings struct {
	Name         						types.String	`tfsdk:"name"`	
	Id           						types.String	`tfsdk:"id"`
	Affinity_cookie_name           		types.String	`tfsdk:"affinity_cookie_name"`					
	Cookie_based_affinity          		types.String	`tfsdk:"cookie_based_affinity"`					
	Pick_host_name_from_backend_address types.Bool		`tfsdk:"pick_host_name_from_backend_address"`	
	Port								types.Int64		`tfsdk:"port"`									
	Protocol                       		types.String	`tfsdk:"protocol"`								
	Request_timeout						types.Int64		`tfsdk:"request_timeout"`						
	Probe_name							types.String	`tfsdk:"probe_name"`							
}

func createBackendHTTPSettings(backend_plan Backend_http_settings,AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) BackendHTTPSettings{
	fmt.Printf("\nIIIIIIIIIIIIIIIIIIII  backend_plan =\n %+v ",backend_plan)
	backend_json := BackendHTTPSettings{
		Name:       backend_plan.Name.Value,
		//ID:         "",
		//Etag:       "",
		Properties: struct{
			AffinityCookieName 	string "json:\"affinityCookieName,omitempty\""; 
			CookieBasedAffinity string "json:\"cookieBasedAffinity,omitempty\""; 
			PickHostNameFromBackendAddress bool "json:\"pickHostNameFromBackendAddress,omitempty\""; 
			Port int "json:\"port,omitempty\""; 
			Probe *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"probe\""; 
			Protocol string "json:\"protocol,omitempty\""; 
			RequestTimeout int "json:\"requestTimeout,omitempty\""; 

			AuthenticationCertificates *[]struct{
				ID string "json:\"id\""
			} "json:\"authenticationCertificates\""; 
			ConnectionDraining *struct{
				DrainTimeoutInSec int "json:\"drainTimeoutInSec,omitempty\""; 
				Enabled bool "json:\"enabled,omitempty\""
			} "json:\"connectionDraining\""; 
			HostName string "json:\"hostName,omitempty\""; 
			Path string "json:\"path,omitempty\""; 
			ProbeEnabled bool "json:\"probeEnabled,omitempty\""; 
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			RequestRoutingRules *[]struct{
				ID string "json:\"id,omitempty\""
			} "json:\"requestRoutingRules\""; 
			TrustedRootCertificates *[]struct{
				ID string "json:\"id,omitempty\""
			} "json:\"trustedRootCertificates\""
			}{	//initialisation of the Properties Struct
				CookieBasedAffinity:			backend_plan.Cookie_based_affinity.Value,
				//if PickHostNameFromBackendAddress attribute become optional with false default value, this line should be replaced
				PickHostNameFromBackendAddress: bool(backend_plan.Pick_host_name_from_backend_address.Value),
				Port: 							int(backend_plan.Port.Value),
				Protocol: 						backend_plan.Protocol.Value,
				RequestTimeout: 				int(backend_plan.Request_timeout.Value),
			},
		Type: "Microsoft.Network/applicationGateways/backendHttpSettingsCollection",
	}

	
	//the probe name should treated specifically to construct the ID
	probe_string := "/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+"/providers/Microsoft.Network/applicationGateways/"+agw_name+"/probes/"
	// if there is à probe, then copy it, else, nil
	if backend_plan.Probe_name.Value != "" {
		backend_json.Properties.Probe = &struct{
			ID string "json:\"id,omitempty\""
		}{
			ID: probe_string + backend_plan.Probe_name.Value,
		}
	}else{
		//backend_json.Properties.Probe = 
	}

	//AffinityCookieName: 			backend_plan.Affinity_cookie_name.Value,
	if !backend_plan.Affinity_cookie_name.Null {
		backend_json.Properties.AffinityCookieName = backend_plan.Affinity_cookie_name.Value
	}else{
		backend_json.Properties.AffinityCookieName = ""
	}		
	fmt.Printf("\nHHHHHHHHHHHHHH  backend_json =\n %+v ",backend_json)
	
	// add the backend to the agw and update the agw
	return backend_json
}
func generateBackendHTTPSettingsState(gw ApplicationGateway, BackendHTTPSettingsName string) Backend_http_settings {
	// we have to give the nb_Fqdns and nb_IpAddress in order to make this function reusable in create, read and update method
	index := getBackendHTTPSettingsElementKey(gw, BackendHTTPSettingsName)
	backend_json := gw.Properties.BackendHTTPSettingsCollection[index]
	
	// Map response body to resource schema attribute
	//split the probe ID using the separator "/". the probe name is the last one
	
	var backend_state Backend_http_settings
	backend_state = Backend_http_settings{
		Name:                                types.String	{Value: backend_json.Name},
		Id:                                  types.String	{Value: backend_json.ID},
		Affinity_cookie_name:                types.String	{},
		Cookie_based_affinity:               types.String	{Value: backend_json.Properties.CookieBasedAffinity},
		Pick_host_name_from_backend_address: types.Bool		{Value: backend_json.Properties.PickHostNameFromBackendAddress},
		Port:                                types.Int64	{Value: int64(backend_json.Properties.Port)},
		Protocol:                            types.String	{Value: backend_json.Properties.Protocol},
		Request_timeout:                     types.Int64	{Value: int64(backend_json.Properties.RequestTimeout)},
		Probe_name:                          types.String	{},
	}
	//verify if optional parameters are provided, otherwise, they have to set to null
	if backend_json.Properties.Probe != nil {
		splitted_list := strings.Split(backend_json.Properties.Probe.ID,"/")
		backend_state.Probe_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		backend_state.Probe_name = types.String{Null: true}
	}	
	if (backend_json.Properties.AffinityCookieName != "")&&
		(&backend_json.Properties.AffinityCookieName != nil) {
		backend_state.Affinity_cookie_name = types.String {Value: backend_json.Properties.AffinityCookieName}
	}else{
		backend_state.Affinity_cookie_name = types.String{Value: ""}
	}

	fmt.Printf("\n--------------------- BackendHTTPSettings state createBackendHTTPSettings() =\n %+v ",backend_state)
	
	return backend_state
}
func getBackendHTTPSettingsElementKey(gw ApplicationGateway, BackendHTTPSettingsName string) int {
	key := -1
	for i := len(gw.Properties.BackendHTTPSettingsCollection) - 1; i >= 0; i-- {
		if gw.Properties.BackendHTTPSettingsCollection[i].Name == BackendHTTPSettingsName {
			key = i
		}
	}
	return key
}
func checkBackendHTTPSettingsElement(gw ApplicationGateway, BackendHTTPSettingsName string) bool {
	exist := false
	for i := len(gw.Properties.BackendHTTPSettingsCollection) - 1; i >= 0; i-- {
		if gw.Properties.BackendHTTPSettingsCollection[i].Name == BackendHTTPSettingsName {
			exist = true
		}
	}
	return exist
}
func removeBackendHTTPSettingsElement(gw *ApplicationGateway, BackendHTTPSettingsName string) {
	for i := len(gw.Properties.BackendHTTPSettingsCollection) - 1; i >= 0; i-- {
		if gw.Properties.BackendHTTPSettingsCollection[i].Name == BackendHTTPSettingsName {
			gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection[:i], gw.Properties.BackendHTTPSettingsCollection[i+1:]...)
		}
	}
}