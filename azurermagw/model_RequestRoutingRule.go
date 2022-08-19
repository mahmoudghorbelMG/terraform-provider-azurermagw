package azurermagw

import (
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RequestRoutingRule struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState string `json:"provisioningState,omitempty"`
		RuleType          string `json:"ruleType,omitempty"`
		Priority          int    `json:"priority,omitempty"`
		HTTPListener      *struct {
			ID string `json:"id,omitempty"`
		} `json:"httpListener"`
		BackendAddressPool *struct {
			ID string `json:"id,omitempty"`
		} `json:"backendAddressPool"`
		BackendHTTPSettings *struct {
			ID string `json:"id,omitempty"`
		} `json:"backendHttpSettings"`
		LoadDistributionPolicy *struct {
			ID string `json:"id,omitempty"`
		} `json:"loadDistributionPolicy"`
		RedirectConfiguration *struct {
			ID string `json:"id,omitempty"`
		} `json:"redirectConfiguration"`
		RewriteRuleSet *struct {
			ID string `json:"id,omitempty"`
		} `json:"rewriteRuleSet"`
		URLPathMap *struct {
			ID string `json:"id,omitempty"`
		} `json:"urlPathMap"`
	} `json:"properties"`
	Type string `json:"type"`
}

type Request_routing_rule struct {
	//required
	Name         						types.String	`tfsdk:"name"`	
	Id           						types.String	`tfsdk:"id"`
	Rule_type							types.String	`tfsdk:"rule_type"`					
	Http_listener_name           		types.String	`tfsdk:"http_listener_name"`
	Priority 							types.Int64		`tfsdk:"priority"`
	//Cannot be set if redirect_configuration_name is not set
	Backend_address_pool_name			types.String	`tfsdk:"backend_address_pool_name"`
	Backend_http_settings_name			types.String	`tfsdk:"backend_http_settings_name"`								
	//Cannot be set if both backend_address_pool_name and backend_http_settings_name are not set
	Redirect_configuration_name			types.String	`tfsdk:"redirect_configuration_name"`	
	//Only valid for v2 SKUs.
	Rewrite_rule_set_name  				types.String	`tfsdk:"rewrite_rule_set_name"`	
	//optional
	Url_path_map_name					types.String		`tfsdk:"url_path_map_name"`
}

func createRequestRoutingRule(requestRoutingRule_plan *Request_routing_rule, priority int, AZURE_SUBSCRIPTION_ID string, 
								rg_name string, agw_name string) (RequestRoutingRule){
	requestRoutingRule_json := RequestRoutingRule{
		Name:       requestRoutingRule_plan.Name.Value,
		ID:         "",
		Etag:       "",
		Properties: struct{
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			RuleType string "json:\"ruleType,omitempty\""; 
			Priority int "json:\"priority,omitempty\""; 
			HTTPListener *struct{ID string "json:\"id,omitempty\""} "json:\"httpListener\""; 
			BackendAddressPool *struct{ID string "json:\"id,omitempty\""} "json:\"backendAddressPool\""; 
			BackendHTTPSettings *struct{ID string "json:\"id,omitempty\""} "json:\"backendHttpSettings\""; 
			LoadDistributionPolicy *struct{ID string "json:\"id,omitempty\""} "json:\"loadDistributionPolicy\""; 
			RedirectConfiguration *struct{ID string "json:\"id,omitempty\""} "json:\"redirectConfiguration\""; 
			RewriteRuleSet *struct{ID string "json:\"id,omitempty\""} "json:\"rewriteRuleSet\""; 
			URLPathMap *struct{ID string "json:\"id,omitempty\""} "json:\"urlPathMap\""
			}{
				RuleType		: requestRoutingRule_plan.Rule_type.Value,
				Priority		: priority,				
			},
		Type: "Microsoft.Network/applicationGateways/requestRoutingRules",
	}
	ID:="/subscriptions/"+AZURE_SUBSCRIPTION_ID+"/resourceGroups/"+rg_name+"/providers/Microsoft.Network/applicationGateways/"+agw_name
	
	HTTPListenerID :=ID+"/httpListeners/"+requestRoutingRule_plan.Http_listener_name.Value
	requestRoutingRule_json.Properties.HTTPListener = &struct{ID string "json:\"id,omitempty\""}{ID: HTTPListenerID,}
	
	if requestRoutingRule_plan.Redirect_configuration_name.Value != "" {
		//redirect_configuration_name is set
		redirectConfigurationID := ID+"/redirectConfigurations/"+requestRoutingRule_plan.Redirect_configuration_name.Value
		requestRoutingRule_json.Properties.RedirectConfiguration = &struct{ID string "json:\"id,omitempty\""}{ID: redirectConfigurationID,}
	}else{
		//backend_address_pool_name and backend_http_settings_name are set:
		backendAddressPoolID := ID+"/backendAddressPools/"+requestRoutingRule_plan.Backend_address_pool_name.Value
		requestRoutingRule_json.Properties.BackendAddressPool = &struct{ID string "json:\"id,omitempty\""}{ID: backendAddressPoolID,}
	
		backendHttpSettingsID := ID+"/backendHttpSettingsCollection/"+requestRoutingRule_plan.Backend_http_settings_name.Value
		requestRoutingRule_json.Properties.BackendHTTPSettings = &struct{ID string "json:\"id,omitempty\""}{ID: backendHttpSettingsID,}
	}
	if requestRoutingRule_plan.Rewrite_rule_set_name.Value != ""{
		//rewrite_rule_set_name is set
		rewriteRuleSetID := ID+"/rewriteRuleSets/"+requestRoutingRule_plan.Rewrite_rule_set_name.Value
		requestRoutingRule_json.Properties.RewriteRuleSet = &struct{ID string "json:\"id,omitempty\""}{ID: rewriteRuleSetID,}
	}
	if requestRoutingRule_plan.Url_path_map_name.Value != "" {
		//url_path_map_name is set
		URLPathMapID:= ID+"/urlPathMaps/"+requestRoutingRule_plan.Url_path_map_name.Value
		requestRoutingRule_json.Properties.URLPathMap = &struct{ID string "json:\"id,omitempty\""}{ID: URLPathMapID,}
	}
	
	return requestRoutingRule_json
}
func generateRequestRoutingRuleState(gw ApplicationGateway, RequestRoutingRuleName string) Request_routing_rule {
	//retrieve json element from gw
	index := getRequestRoutingRuleElementKey_gw(gw, RequestRoutingRuleName)
	requestRoutingRule_json := gw.Properties.RequestRoutingRules[index]
	
	// Map response body to resource schema attribute
	var requestRoutingRule_state Request_routing_rule
	requestRoutingRule_state = Request_routing_rule{
		Name:                        types.String{Value: requestRoutingRule_json.Name},
		Id:                          types.String{Value: requestRoutingRule_json.ID},
		Rule_type:                   types.String{Value: requestRoutingRule_json.Properties.RuleType},
		Http_listener_name:          types.String{},
		Priority:                    types.Int64{Value: int64(requestRoutingRule_json.Properties.Priority)},
		Backend_address_pool_name:   types.String{},
		Backend_http_settings_name:  types.String{},
		Redirect_configuration_name: types.String{},
		Rewrite_rule_set_name:       types.String{},
		Url_path_map_name:           types.String{},
	}
	splitted_list := strings.Split(requestRoutingRule_json.Properties.HTTPListener.ID,"/")
	requestRoutingRule_state.Http_listener_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	
	if requestRoutingRule_json.Properties.RedirectConfiguration != nil {
		splitted_list := strings.Split(requestRoutingRule_json.Properties.RedirectConfiguration.ID,"/")
		requestRoutingRule_state.Redirect_configuration_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		requestRoutingRule_state.Redirect_configuration_name = types.String{Null: true}
	}
	if requestRoutingRule_json.Properties.BackendAddressPool != nil {
		splitted_list := strings.Split(requestRoutingRule_json.Properties.BackendAddressPool.ID,"/")
		requestRoutingRule_state.Backend_address_pool_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		requestRoutingRule_state.Backend_address_pool_name = types.String{Null: true}
	}
	if requestRoutingRule_json.Properties.BackendHTTPSettings != nil {
		splitted_list := strings.Split(requestRoutingRule_json.Properties.BackendHTTPSettings.ID,"/")
		requestRoutingRule_state.Backend_http_settings_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		requestRoutingRule_state.Backend_http_settings_name = types.String{Null: true}
	}
	if requestRoutingRule_json.Properties.RewriteRuleSet != nil {
		splitted_list := strings.Split(requestRoutingRule_json.Properties.RewriteRuleSet.ID,"/")
		requestRoutingRule_state.Rewrite_rule_set_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		requestRoutingRule_state.Rewrite_rule_set_name = types.String{Null: true}
	}
	if requestRoutingRule_json.Properties.URLPathMap != nil {
		splitted_list := strings.Split(requestRoutingRule_json.Properties.URLPathMap.ID,"/")
		requestRoutingRule_state.Url_path_map_name = types.String{Value: splitted_list[len(splitted_list)-1]}
	}else{
		requestRoutingRule_state.Url_path_map_name = types.String{Null: true}
	}

	return requestRoutingRule_state
}
func getRequestRoutingRuleElementKey_gw(gw ApplicationGateway, RequestRoutingRuleName string) int {
	key := -1
	for i := len(gw.Properties.RequestRoutingRules) - 1; i >= 0; i-- {
		if gw.Properties.RequestRoutingRules[i].Name == RequestRoutingRuleName {
			key = i
		}
	}
	return key
}
func checkRequestRoutingRuleElement(gw ApplicationGateway, RequestRoutingRuleName string) bool {
	exist := false
	for i := len(gw.Properties.RequestRoutingRules) - 1; i >= 0; i-- {
		if gw.Properties.RequestRoutingRules[i].Name == RequestRoutingRuleName {
			exist = true
		}
	}
	return exist
}
func removeRequestRoutingRuleElement(gw *ApplicationGateway, RequestRoutingRuleName string) {
	for i := len(gw.Properties.RequestRoutingRules) - 1; i >= 0; i-- {
		if gw.Properties.RequestRoutingRules[i].Name == RequestRoutingRuleName {
			gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules[:i], gw.Properties.RequestRoutingRules[i+1:]...)
		}
	}
}