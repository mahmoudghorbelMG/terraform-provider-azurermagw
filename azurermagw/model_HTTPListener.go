package azurermagw

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HTTPListener struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		FirewallPolicy *struct {
			ID string `json:"id,omitempty"`
		} `json:"firewallPolicy"`
		ProvisioningState       string `json:"provisioningState,omitempty"`
		FrontendIPConfiguration *struct {ID string `json:"id,omitempty"`} `json:"frontendIPConfiguration"`
		FrontendPort *struct {ID string `json:"id,omitempty"`} `json:"frontendPort,omitempty"`
		Protocol                    string   `json:"protocol,omitempty"`
		HostName                    string   `json:"hostName,omitempty"`
		HostNames                   []string `json:"hostNames,omitempty"`
		RequireServerNameIndication bool     `json:"requireServerNameIndication,omitempty"`
		SslCertificate *struct {ID string `json:"id,omitempty"`} `json:"sslCertificate"`
		SslProfile *struct {ID string `json:"id,omitempty"`	} `json:"sslProfile"`
		CustomErrorConfigurations *[]struct {
			CustomErrorPageURL string `json:"customErrorPageUrl,omitempty"`
			StatusCode         string `json:"statusCode,omitempty"`
		} `json:"customErrorConfigurations"`
		RequestRoutingRules *[]struct {	ID string `json:"id,omitempty"`} `json:"requestRoutingRules,omitempty"`
	} `json:"properties"`
	Type string `json:"type,omitempty"`
} 

type Http_listener struct {
	//required
	Name         						types.String	`tfsdk:"name"`	
	Id           						types.String	`tfsdk:"id"`
	Frontend_ip_configuration_name		types.String	`tfsdk:"frontend_ip_configuration_name"`					
	Frontend_port_name           		types.String	`tfsdk:"frontend_port_name"`					
	Protocol                       		types.String	`tfsdk:"protocol"`								
	//optional but mutually exclusive
	Host_name  							types.String	`tfsdk:"host_name"`	
	Host_names  						[]types.String	`tfsdk:"host_names"`	
	//optional
	Require_sni 						types.Bool		`tfsdk:"require_sni"`	//default to false
	Ssl_certificate_name 				types.String	`tfsdk:"ssl_certificate_name"`							
	//Ssl_profile_name 
	//Firewall_policy_id 
	//Custom_error_configuration 

}

func createHTTPListener(httpListener_plan *Http_listener, AZURE_SUBSCRIPTION_ID string, 
								rg_name string, agw_name string) (HTTPListener){
	httpListener_json := HTTPListener{
		Name:       httpListener_plan.Name.Value,
		//ID:         AZURE_SUBSCRIPTION_ID,
		//Etag:       "",
		Properties: struct{
			FirewallPolicy *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"firewallPolicy\""; 
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			FrontendIPConfiguration *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"frontendIPConfiguration\""; 
			FrontendPort *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"frontendPort,omitempty\""; 
			Protocol string "json:\"protocol,omitempty\""; 
			HostName string "json:\"hostName,omitempty\""; 
			HostNames []string "json:\"hostNames,omitempty\""; 
			RequireServerNameIndication bool "json:\"requireServerNameIndication,omitempty\""; 
			SslCertificate *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"sslCertificate\""; 
			SslProfile *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"sslProfile\""; 
			CustomErrorConfigurations *[]struct{
				CustomErrorPageURL string "json:\"customErrorPageUrl,omitempty\""; 
				StatusCode string "json:\"statusCode,omitempty\""
			} "json:\"customErrorConfigurations\""; 
			RequestRoutingRules *[]struct{
				ID string "json:\"id,omitempty\""
			} "json:\"requestRoutingRules,omitempty\""
		}{
			Protocol: httpListener_plan.Protocol.Value,			
			RequireServerNameIndication: bool(httpListener_plan.Require_sni.Value),
		},
		Type: "Microsoft.Network/applicationGateways/httpListeners",
	}
	
	//frontendIPConfiguration is required, so no test to do
	frontendIPConfigurationID :="/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+
				"/providers/Microsoft.Network/applicationGateways/"+agw_name+"/frontendIPConfigurations/"+
				httpListener_plan.Frontend_ip_configuration_name.Value
	httpListener_json.Properties.FrontendIPConfiguration = &struct{ID string "json:\"id,omitempty\""}{ID: frontendIPConfigurationID,}

	//frontendPort is required, so no test to do
	frontendPortID :="/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+
	"/providers/Microsoft.Network/applicationGateways/"+agw_name+"/frontendPorts/"+httpListener_plan.Frontend_port_name.Value
	httpListener_json.Properties.FrontendPort = &struct{ID string "json:\"id,omitempty\""}{ID: frontendPortID,}

	//ssl certificate id is optional, but when provided, it has to be conform with the certificate name in the binding
	sslCertificateID := "/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+
	"/providers/Microsoft.Network/applicationGateways/"+agw_name+"/sslCertificates/"
	// if there is a Ssl_certificate_name, then put it, else, nil
	//var error_SslCertificateName string
	if httpListener_plan.Ssl_certificate_name.Value != "" {
		//we have to check here if the probe name matches probe name in terraform conf in plan.
		httpListener_json.Properties.SslCertificate = &struct{
			ID string "json:\"id,omitempty\""
		}{
			ID: sslCertificateID + httpListener_plan.Ssl_certificate_name.Value,
		}		
	}
	
	//verify the mutual exclusivity of the optional attributes hostname and hostnames
	//var error_Hostname string
	if httpListener_plan.Host_name.Value != "" {
			//hostname is provided but not hostnames
			httpListener_json.Properties.HostName = httpListener_plan.Host_name.Value	
	}
	if len(httpListener_plan.Host_names)!=0 {
		//hostnames is provided but not hostname
		httpListener_json.Properties.HostNames = make([]string,len(httpListener_plan.Host_names))
		for i := 0; i < len(httpListener_plan.Host_names); i++ {
			httpListener_json.Properties.HostNames[i] = httpListener_plan.Host_names[i].Value
		}
	}
	return httpListener_json
}
func generateHTTPListenerState(gw ApplicationGateway, HTTPListenerName string) Http_listener {
	//retrieve json element from gw
	index := getHTTPListenerElementKey_gw(gw, HTTPListenerName)
	httpListener_json := gw.Properties.HTTPListeners[index]
	
	// Map response body to resource schema attribute
	var httpListener_state Http_listener
	httpListener_state = Http_listener{
		Name:                       	types.String	{Value: httpListener_json.Name},
		Id:                         	types.String	{Value: httpListener_json.ID},
		Protocol:                   	types.String	{Value: httpListener_json.Properties.Protocol},
		Require_sni: 					types.Bool		{Value: httpListener_json.Properties.RequireServerNameIndication},
		Ssl_certificate_name:       	types.String	{},
		Frontend_ip_configuration_name:	types.String	{},
		Frontend_port_name:       		types.String	{},
		//Host_name:       				types.String	{},
		Host_names:						[]types.String	{},			
	}
	if httpListener_json.Properties.HostName != "" {
		httpListener_state.Host_name = types.String{Value: httpListener_json.Properties.HostName}
	}else{
		httpListener_state.Host_name.Null = true
	}
	//map host_names. check if it is an empty array.
	if len(httpListener_json.Properties.HostNames)!=0 {
		httpListener_state.Host_names = make([]types.String,len(httpListener_json.Properties.HostNames))
		for i := 0; i < len(httpListener_json.Properties.HostNames); i++ {
			httpListener_state.Host_names[i] = types.String{Value: httpListener_json.Properties.HostNames[i]}
		}
	}else{
		//we have to verify if it has to be nil or empty array
		httpListener_state.Host_names = nil
	}

	//map Frontend_ip_configuration_name
	//split the Frontend_ip_configuration_name ID using the separator "/". the Frontend_ip_configuration_name name is the last one
	if httpListener_json.Properties.FrontendIPConfiguration != nil {
		splitted_list := strings.Split(httpListener_json.Properties.FrontendIPConfiguration.ID,"/")
		httpListener_state.Frontend_ip_configuration_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		httpListener_state.Frontend_ip_configuration_name = types.String{Null: true}
	}

	//map Frontend_port_name
	//split the Frontend_port_name ID using the separator "/". the Frontend_port_name name is the last one
	if httpListener_json.Properties.FrontendPort != nil {
		splitted_list := strings.Split(httpListener_json.Properties.FrontendPort.ID,"/")
		httpListener_state.Frontend_port_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		httpListener_state.Frontend_port_name = types.String{Null: true}
	}

	//map Ssl_certificate_name
	//split the Ssl_certificate_name ID using the separator "/". the Ssl_certificate_name name is the last one
	if httpListener_json.Properties.SslCertificate != nil {
		splitted_list := strings.Split(httpListener_json.Properties.SslCertificate.ID,"/")
		httpListener_state.Ssl_certificate_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		httpListener_state.Ssl_certificate_name = types.String{Null: true}
	}

	return httpListener_state
}
func getHTTPListenerElementKey_gw(gw ApplicationGateway, HTTPListenerName string) int {
	key := -1
	for i := len(gw.Properties.HTTPListeners) - 1; i >= 0; i-- {
		if gw.Properties.HTTPListeners[i].Name == HTTPListenerName {
			key = i
		}
	}
	return key
}
func getHTTPListenerElementKey_state(httpListeners_state []Http_listener, httpListener_plan Http_listener) int {
	key := -1
	for i := len(httpListeners_state) - 1; i >= 0; i-- {
		//3 conditions has to be satisfied: 1)same Frontend Port, 2) same FrontendIpConfiguration and 3)same HostName or HostNames.
		//condition 1
		condition_1 := false
		if httpListeners_state[i].Frontend_port_name.Value == httpListener_plan.Frontend_port_name.Value {
			condition_1 = true
		}
		condition_2 := false
		if httpListeners_state[i].Frontend_ip_configuration_name.Value == httpListener_plan.Frontend_ip_configuration_name.Value {
			condition_2 = true
		}
		condition_3 := false		
		if httpListeners_state[i].Host_name.Value == httpListener_plan.Host_name.Value && httpListener_plan.Host_name.Value != ""{
			condition_3 = true
		}else{
			if httpListener_plan.Host_name.Value == "" {			
				if httpListeners_state[i].Host_name.Value == "" {
					for i := 0; i < len(httpListeners_state[i].Host_names); i++ {
						if check(httpListener_plan.Host_names,httpListeners_state[i].Host_names[i].Value){
							condition_3 = true
						}
					}
				}else{
					condition_3 = check(httpListener_plan.Host_names,httpListeners_state[i].Host_name.Value)					
				}
			}else{
				if httpListeners_state[i].Host_name.Value=="" {
					condition_3 = check(httpListeners_state[i].Host_names,httpListener_plan.Host_name.Value)
				}
			}
		}
		if condition_1&&condition_2&&condition_3 {
			key = i
			return key
		}		
	}
	return key
}
func check(hostnames []types.String, hostname string) bool{
	for i := 0; i < len(hostnames); i++ {
		if hostnames[i].Value == hostname{
			return true
		}
	}
	return false
}
func checkHTTPListenerElement(gw ApplicationGateway, HTTPListenerName string) bool {
	exist := false
	for i := len(gw.Properties.HTTPListeners) - 1; i >= 0; i-- {
		if gw.Properties.HTTPListeners[i].Name == HTTPListenerName {
			exist = true
		}
	}
	return exist
}
func removeHTTPListenerElement(gw *ApplicationGateway, HTTPListenerName string) {
	for i := len(gw.Properties.HTTPListeners) - 1; i >= 0; i-- {
		if gw.Properties.HTTPListeners[i].Name == HTTPListenerName {
			gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners[:i], gw.Properties.HTTPListeners[i+1:]...)
		}
	}
}
/*
func checkHTTPListenerCreate(plan BindingService, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) bool {
	if plan.Http_listener.Ssl_certificate_name.Value != "" {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to create binding. A SslCertificate name ("+plan.Http_listener.Ssl_certificate_name.Value+") is declared in Http_listener: "+ 
		plan.Http_listener.Name.Value+". There is no need for SslCertificate name in this case. ",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	if plan.Http_listener.Host_name.Value != "" && len(plan.Http_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTPS Listener "+ plan.Http_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	if plan.Http_listener.Host_name.Value == "" && len(plan.Http_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTP Listener "+ plan.Http_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	return false
}
func checkHTTPListenerUpdate(plan BindingService, gw ApplicationGateway, resp *tfsdk.UpdateResourceResponse) bool {
	if plan.Http_listener.Ssl_certificate_name.Value != "" {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to update binding. A SslCertificate name ("+plan.Http_listener.Ssl_certificate_name.Value+") is declared in Http_listener: "+ 
		plan.Http_listener.Name.Value+". There is no need for SslCertificate name in this case. ",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	if plan.Http_listener.Host_name.Value != "" && len(plan.Http_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTPS Listener "+ plan.Http_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	if plan.Http_listener.Host_name.Value == "" && len(plan.Http_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTP Listener "+ plan.Http_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	return false
}
func checkHTTPSListenerCreate(plan BindingService, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) bool {
	if plan.Https_listener.Ssl_certificate_name.Value != "" &&
		plan.Https_listener.Ssl_certificate_name.Value != plan.Ssl_certificate.Name.Value{
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to create binding. The SslCertificate name ("+plan.Https_listener.Ssl_certificate_name.Value+") declared in Https_listener: "+ 
		plan.Https_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Ssl_certificate.Name.Value,
		"Please, change Ssl Certificate name then retry.",)
		return true
	}
	if plan.Https_listener.Host_name.Value != "" && len(plan.Https_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTPS Listener "+ plan.Https_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTPS Listener configuration then retry.",)
		return true
	}
	if plan.Https_listener.Host_name.Value == "" && len(plan.Https_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTP Listener "+ plan.Https_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTPS Listener configuration then retry.",)
		return true
	}
	return false
}
func checkHTTPSListenerUpdate(plan BindingService, gw ApplicationGateway, resp *tfsdk.UpdateResourceResponse) bool {
	if plan.Https_listener.Ssl_certificate_name.Value != "" &&
		plan.Https_listener.Ssl_certificate_name.Value != plan.Ssl_certificate.Name.Value{
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to update binding. The SslCertificate name ("+plan.Https_listener.Ssl_certificate_name.Value+") declared in Https_listener: "+ 
		plan.Https_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Ssl_certificate.Name.Value,
		"Please, change Ssl Certificate name then retry.",)
		return true
	}
	if plan.Https_listener.Host_name.Value != "" && len(plan.Https_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTPS Listener "+ plan.Https_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTPS Listener configuration then retry.",)
		return true
	}
	if plan.Https_listener.Host_name.Value == "" && len(plan.Https_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTP Listener "+ plan.Https_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTPS Listener configuration then retry.",)
		return true
	}
	return false
}*/
func checkHTTPListenerCreate(http_listener Http_listener, plan BindingService, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) bool {
	if http_listener.Ssl_certificate_name.Value != "" && strings.EqualFold(http_listener.Protocol.Value,"http") {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to create binding. A SslCertificate name ("+http_listener.Ssl_certificate_name.Value+") is declared in Http_listener: "+ 
		http_listener.Name.Value+". There is no need for SslCertificate name in this case. ",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	if http_listener.Ssl_certificate_name.Value == "" && strings.EqualFold(http_listener.Protocol.Value,"https") {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to create binding. A SslCertificate name is missing for the Http_listener: "+ 
		http_listener.Name.Value+".",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	//if it's about https, check if the certificate name match the one declared un the binding service, or (coming soon) in the gw
	if http_listener.Ssl_certificate_name.Value != "" &&
		http_listener.Ssl_certificate_name.Value != plan.Ssl_certificate.Name.Value{
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to create binding. The SslCertificate name ("+http_listener.Ssl_certificate_name.Value+") declared in Http_listener: "+ 
		http_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Ssl_certificate.Name.Value,
		"Please, change Ssl Certificate name then retry.",)
		return true
	}
	if http_listener.Host_name.Value != "" && len(http_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTPS Listener "+ http_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	if http_listener.Host_name.Value == "" && len(http_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTP Listener "+ http_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	return false
}
func checkHTTPListenerUpdate(http_listener Http_listener, plan BindingService, gw ApplicationGateway, resp *tfsdk.UpdateResourceResponse) bool {
	if http_listener.Ssl_certificate_name.Value != "" && strings.EqualFold(http_listener.Protocol.Value,"http") {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to update binding. A SslCertificate name ("+http_listener.Ssl_certificate_name.Value+") is declared in Http_listener: "+ 
		http_listener.Name.Value+". There is no need for SslCertificate name in this case. ",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	if http_listener.Ssl_certificate_name.Value == "" && strings.EqualFold(http_listener.Protocol.Value,"https") {
		//no need for SslCertificate
		resp.Diagnostics.AddError(
		"Unable to update binding. A SslCertificate name is missing for the Http_listener: "+ 
		http_listener.Name.Value+".",
		"Please, change Http listener configuration then retry.",)
		return true
	}
	//if it's about https, check if the certificate name match the one declared un the binding service, or in the gw
	if http_listener.Ssl_certificate_name.Value != "" &&
		http_listener.Ssl_certificate_name.Value != plan.Ssl_certificate.Name.Value{
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to update binding. The SslCertificate name ("+http_listener.Ssl_certificate_name.Value+") declared in Http_listener: "+ 
		http_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Ssl_certificate.Name.Value,
		"Please, change Ssl Certificate name then retry.",)
		return true
	}
	if http_listener.Host_name.Value != "" && len(http_listener.Host_names) != 0 {
		//hostname and hostnames are mutually exclusive. only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTPS Listener "+ http_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	if http_listener.Host_name.Value == "" && len(http_listener.Host_names) == 0 {
		//hostname and hostnames are mutually exclusive. at least and only one should be set
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTP Listener "+ http_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be set",
			"Please, change HTTP Listener configuration then retry.",)
		return true
	}
	return false
}
func checkHTTPListenerNameInMap(HTTPListenerName string, http_listeners map[string]Http_listener) bool{
	for _, value := range http_listeners {
		if HTTPListenerName == value.Name.Value {
			return true
		}
	}
	return false
}