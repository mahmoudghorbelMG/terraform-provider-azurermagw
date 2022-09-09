package azurermagw

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func createRedirectConfiguration(redirectConfiguration_plan Redirect_configuration, AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) (RedirectConfiguration){	
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
	
	//var error_exclusivity string
	//var error_target string
	//there is a constraint for we have to check: Target_listener_name and target_url are mutually exclusive. 
	//only one of has to be set
	if redirectConfiguration_plan.Target_listener_name.Value != "" {
		redirectConfiguration_json.Properties.TargetListener = &struct{
				ID string "json:\"id,omitempty\""
			}{
				ID: target_listener_string + redirectConfiguration_plan.Target_listener_name.Value,
			}		
	}	
	if redirectConfiguration_plan.Target_url.Value != "" {
		//only Target_url is set.
		redirectConfiguration_json.Properties.TargetURL = redirectConfiguration_plan.Target_url.Value
	}
	return redirectConfiguration_json
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
func checkRedirectConfigurationCreate(plan BindingService, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) bool {
	//fatal-both-exist
	if plan.Redirect_configuration.Target_listener_name.Value != "" &&
		plan.Redirect_configuration.Target_url.Value != "" {
		resp.Diagnostics.AddError(
		"Unable to create binding. In the Redirect Configuration ("+plan.Redirect_configuration.Name.Value+"), 2 optional parameters mutually exclusive "+ 
		"are declared: Target_listener_name and Target_url. Only one has to be set. ",
		"Please, change configuration then retry.",)
		return true
	}
	//fatal both don't exist
	if plan.Redirect_configuration.Target_listener_name.Value == "" &&
	plan.Redirect_configuration.Target_url.Value == "" {
		resp.Diagnostics.AddError(
		"Unable to create binding. In the Redirect Configuration  ("+plan.Redirect_configuration.Name.Value+"), both optional parameters mutually exclusive "+ 
		"are missing: Target_listener_name and Target_url. At least and only one has to be set. ",
		"Please, change configuration then retry.",)
		return true
	}	
	// check if the given Target_listener_name exist in http_listeners map or in the gw
	if plan.Redirect_configuration.Target_listener_name.Value != "" &&
		!checkHTTPListenerNameInMap(plan.Redirect_configuration.Target_listener_name.Value, plan.Http_listeners) &&
		!checkHTTPListenerElement(gw, plan.Redirect_configuration.Target_listener_name.Value){
		resp.Diagnostics.AddError(
		"Unable to create binding. In the target HTTPS Listener ("+plan.Redirect_configuration.Target_listener_name.Value+") declared in Redirect Configuration : "+ 
		plan.Redirect_configuration.Name.Value+" doesn't match any existing (in the gw) nor declared (in the tf) HTTP Listener. ",
		"Please, change HTTP Listener name then retry.",
		)
		return true
	} 
	return false
}
func checkRedirectConfigurationUpdate(plan BindingService, gw ApplicationGateway, resp *tfsdk.UpdateResourceResponse) bool {
	//fatal-both-exist
	if plan.Redirect_configuration.Target_listener_name.Value != "" &&
		plan.Redirect_configuration.Target_url.Value != "" {
		resp.Diagnostics.AddError(
		"Unable to update binding. In the Redirect Configuration ("+plan.Redirect_configuration.Name.Value+"), 2 optional parameters mutually exclusive "+ 
		"are declared: Target_listener_name and Target_url. Only one has to be set. ",
		"Please, change configuration then retry.",)
		return true
	}
	//fatal both don't exist
	if plan.Redirect_configuration.Target_listener_name.Value == "" &&
	plan.Redirect_configuration.Target_url.Value == "" {
		resp.Diagnostics.AddError(
		"Unable to update binding. In the Redirect Configuration  ("+plan.Redirect_configuration.Name.Value+"), both optional parameters mutually exclusive "+ 
		"are missing: Target_listener_name and Target_url. At least and only one has to be set. ",
		"Please, change configuration then retry.",)
		return true
	}
	// check if the given Target_listener_name exist in http_listeners map or in the gw
	if plan.Redirect_configuration.Target_listener_name.Value != "" &&
		!checkHTTPListenerNameInMap(plan.Redirect_configuration.Target_listener_name.Value, plan.Http_listeners) &&
		!checkHTTPListenerElement(gw, plan.Redirect_configuration.Target_listener_name.Value){
		resp.Diagnostics.AddError(
		"Unable to create binding. In the target HTTPS Listener ("+plan.Redirect_configuration.Target_listener_name.Value+") declared in Redirect Configuration : "+ 
		plan.Redirect_configuration.Name.Value+" doesn't match any existing (in the gw) nor declared (in the tf) HTTP Listener. ",
		"Please, change HTTP Listener name then retry.",
		)
		return true
	} 
	return false
}