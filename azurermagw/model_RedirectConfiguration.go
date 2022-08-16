package azurermagw

import (
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RedirectConfiguration struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState string `json:"provisioningState,omitempty"`
		RedirectType      string `json:"redirectType,omitempty"`
		TargetListener    *struct {
			ID string `json:"id,omitempty"`
		} `json:"targetListener"`
		TargetURL           string `json:"targetUrl,omitempty"`
		IncludePath         bool   `json:"includePath,omitempty"`
		IncludeQueryString  bool   `json:"includeQueryString,omitempty"`
		RequestRoutingRules *[]struct {
			ID string `json:"id,omitempty"`
		} `json:"requestRoutingRules"`
		URLPathMaps *[]struct {
			ID string `json:"id,omitempty"`
		} `json:"urlPathMaps"`
	} `json:"properties"`
	Type string `json:"type"`
}

type Redirect_configuration struct {
	Name         								types.String	`tfsdk:"name"`	
	Id           								types.String	`tfsdk:"id"`
	Redirect_type 								types.String	`tfsdk:"redirect_type"`
	//optional
	Target_listener_name						types.String	`tfsdk:"target_listener_name"`								
	Target_url									types.String	`tfsdk:"target_url"`
	Include_path							 	types.Bool		`tfsdk:"include_path"`	
	Include_query_string					 	types.Bool		`tfsdk:"include_query_string"`	
	
}

func createRedirectConfiguration(redirectConfiguration_plan Redirect_configuration,HTTPSListenerName string, AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) (RedirectConfiguration, string, string){	
	redirectConfiguration_json := RedirectConfiguration{
		Name:       redirectConfiguration_plan.Name.Value,
		//ID:         "",
		//Etag:       "",
		Properties: struct{
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			RedirectType string "json:\"redirectType,omitempty\""; 
			TargetListener *struct{
				ID string "json:\"id,omitempty\""
			} "json:\"targetListener\""; 
			TargetURL string "json:\"targetUrl,omitempty\""; 
			IncludePath bool "json:\"includePath,omitempty\""; 
			IncludeQueryString bool "json:\"includeQueryString,omitempty\""; 
			RequestRoutingRules *[]struct{
				ID string "json:\"id,omitempty\""
			} "json:\"requestRoutingRules\""; 
			URLPathMaps *[]struct{
				ID string "json:\"id,omitempty\""
				} "json:\"urlPathMaps\""
			}{
				RedirectType: redirectConfiguration_plan.Redirect_type.Value,
				IncludePath: bool(redirectConfiguration_plan.Include_path.Value),
				IncludeQueryString: bool(redirectConfiguration_plan.Include_query_string.Value),			
			},
		Type:       "Microsoft.Network/applicationGateways/redirectConfigurations",
	}
	target_listener_string := "/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+"/providers/Microsoft.Network/applicationGateways/"+agw_name+"/httpListeners/"
	
	var error_exclusivity string
	var error_target string
	//there is a constraint for we have to check: Target_listener_name and target_url are mutually exclusive. 
	//only one of has to be set
	if redirectConfiguration_plan.Target_listener_name.Value != "" {
		if redirectConfiguration_plan.Target_url.Value != "" {
			//both are set
			error_exclusivity = "fatal-both-exist"
		}else{
			//only Key_vault_secret_id is set
			//check if its name match the HTTPS listener of the current config, else issue a warning
			if redirectConfiguration_plan.Target_listener_name.Value == HTTPSListenerName {
				redirectConfiguration_json.Properties.TargetListener = &struct{
					ID string "json:\"id,omitempty\""
				}{
					ID: target_listener_string + redirectConfiguration_plan.Target_listener_name.Value,
				}
			}else{
				//Error exit
				error_target = "fatal"
			}			
		}
	}else{
		if redirectConfiguration_plan.Target_url.Value != "" {
			//only Target_url is set.
			redirectConfiguration_json.Properties.TargetURL = redirectConfiguration_plan.Target_url.Value
		}else{
			//both are empty
			error_exclusivity = "fatal-both-miss"
		}
	}
	return redirectConfiguration_json,error_exclusivity,error_target
}
func generateRedirectConfigurationState(gw ApplicationGateway, RedirectConfigurationName string) Redirect_configuration {
	//retrieve json element from gw
	index := getRedirectConfigurationElementKey(gw, RedirectConfigurationName)
	redirectConfiguration_json := gw.Properties.RedirectConfigurations[index]
	
	// Map response body to resource schema attribute
	var redirectConfiguration_state Redirect_configuration
	redirectConfiguration_state = Redirect_configuration{
		Name:                 types.String{Value: redirectConfiguration_json.Name},
		Id:                   types.String{Value: redirectConfiguration_json.ID},
		Redirect_type:        types.String{Value: redirectConfiguration_json.Properties.RedirectType},
		Target_listener_name: types.String{},
		Target_url:           types.String{},
		Include_path:         types.Bool{Value: redirectConfiguration_json.Properties.IncludePath},
		Include_query_string: types.Bool{Value: redirectConfiguration_json.Properties.IncludeQueryString},
	}
	if redirectConfiguration_json.Properties.TargetListener != nil {
		//split the Target Listener ID using the separator "/". the Target Listener name is the last one
		splitted_list := strings.Split(redirectConfiguration_json.Properties.TargetListener.ID,"/")
		redirectConfiguration_state.Target_listener_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		redirectConfiguration_state.Target_listener_name = types.String{Null: true}
	}

	if redirectConfiguration_json.Properties.TargetURL != "" {
		redirectConfiguration_state.Target_url = types.String{Value: redirectConfiguration_json.Properties.TargetURL}
	}else{
		redirectConfiguration_state.Target_url.Null = true
	}
	
	return redirectConfiguration_state
}
func getRedirectConfigurationElementKey(gw ApplicationGateway, RedirectConfigurationName string) int {
	key := -1
	for i := len(gw.Properties.RedirectConfigurations) - 1; i >= 0; i-- {
		if gw.Properties.RedirectConfigurations[i].Name == RedirectConfigurationName {
			key = i
		}
	}
	return key
}
func checkRedirectConfigurationElement(gw ApplicationGateway, RedirectConfigurationName string) bool {
	exist := false
	for i := len(gw.Properties.RedirectConfigurations) - 1; i >= 0; i-- {
		if gw.Properties.RedirectConfigurations[i].Name == RedirectConfigurationName {
			exist = true
		}
	}
	return exist
}
func removeRedirectConfigurationElement(gw *ApplicationGateway, RedirectConfigurationName string) {
	for i := len(gw.Properties.RedirectConfigurations) - 1; i >= 0; i-- {
		if gw.Properties.RedirectConfigurations[i].Name == RedirectConfigurationName {
			gw.Properties.RedirectConfigurations = append(gw.Properties.RedirectConfigurations[:i], gw.Properties.RedirectConfigurations[i+1:]...)
		}
	}
}