package azurermagw

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
func checkRequestRoutingRuleCreate(plan WebappBinding, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) {
	//check http_listener_name (https)
	if plan.Request_routing_rule.Http_listener_name.Value != plan.Https_listener.Name.Value {
		// http_listener_name don't match with Https_listener.Name, issue exit error
		resp.Diagnostics.AddError(
			"Unable to create binding. The Https listener name ("+plan.Request_routing_rule.Http_listener_name.Value+
			") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't match the Https listener name conf : "+
			plan.Https_listener.Name.Value,"Please, change Https listener name then retry.",
		)
		return
	}
	//check mutual exclusivity
	if plan.Request_routing_rule.Redirect_configuration_name.Value != "" {
		//check if one or both are provided, then issue exit error
		if plan.Request_routing_rule.Backend_address_pool_name.Value != "" ||
		 	plan.Request_routing_rule.Backend_http_settings_name.Value != ""{
			// mutual exclusivity error betwenn => exit
			resp.Diagnostics.AddError(
				"Unable to create binding. In the Request Routing Rule  ("+plan.Request_routing_rule.Name.Value+") configuration, "+
				"redirect_configuration_name cannot be set if both backend_address_pool_name or backend_http_settings_name are set ",
				"Please, change configuration then retry.",
				)
			return
		}
		//check redirect_configuration name
		if plan.Request_routing_rule.Redirect_configuration_name.Value != plan.Redirect_configuration.Name.Value {
			// redirect_configuration_name don't match with Redirect_configuration.Name, issue exit error
			resp.Diagnostics.AddError(
				"Unable to create binding. The redirect configuration name ("+plan.Request_routing_rule.Redirect_configuration_name.Value+
				") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't match the redirect configuration name conf : "+
				plan.Redirect_configuration.Name.Value,"Please, change redirect configuration name then retry.",
			)
			return
		}
	}else{
		//check if one or both are missing, then issue exit error
		if plan.Request_routing_rule.Backend_address_pool_name.Value == "" ||
			plan.Request_routing_rule.Backend_http_settings_name.Value == "" {
			// mutual exclusivity error betwenn => exit			
			resp.Diagnostics.AddError(
				"Unable to create binding. In the Request Routing Rule  ("+plan.Request_routing_rule.Name.Value+") configuration, "+
				"a paramameter is missing: [redirect_configuration_name] or [backend_address_pool_name and backend_http_settings_name]",
				"Please, change configuration then retry.",
				)
			return
		}
		//it's ok, check next constraints
		//check backend_address_pool_name 
		if plan.Request_routing_rule.Backend_address_pool_name.Value != plan.Backend_address_pool.Name.Value {
			resp.Diagnostics.AddError(
				"Unable to create binding. The backend address pool name ("+plan.Request_routing_rule.Backend_address_pool_name.Value+
				") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't match the Backend_address_pool name conf : "+
				plan.Backend_address_pool.Name.Value,"Please, change backend address pool name then retry.",
			)
			return
		}
		//check backend_http_settings_name 
		if plan.Request_routing_rule.Backend_http_settings_name.Value != plan.Backend_http_settings.Name.Value {
			resp.Diagnostics.AddError(
				"Unable to create binding. The Backend http settings name ("+plan.Request_routing_rule.Backend_http_settings_name.Value+
				") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't match the Backend http settings name conf : "+
				plan.Backend_http_settings.Name.Value,"Please, change Backend http settings name then retry.",
			)
			return
		}
	}
	//check rewrite_rule_set_name
	if plan.Request_routing_rule.Rewrite_rule_set_name.Value != ""{
		if !checkRewriteRuleSetElement(gw,plan.Request_routing_rule.Rewrite_rule_set_name.Value){
			resp.Diagnostics.AddError(
				"Unable to create binding. The rewrite_rule_set name ("+plan.Request_routing_rule.Rewrite_rule_set_name.Value+
				") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't exist in the gateway.",
				"Please, remove or change rewrite_rule_set name then retry.",
			)
			return
		}
	}
	//check url_path_map_name
	if plan.Request_routing_rule.Url_path_map_name.Value != ""{
		if !checkURLPathMapElement(gw,plan.Request_routing_rule.Url_path_map_name.Value){
			resp.Diagnostics.AddError(
				"Unable to create binding. The url_path_map name ("+plan.Request_routing_rule.Url_path_map_name.Value+
				") declared in Request_routing_rule: "+ plan.Request_routing_rule.Name.Value+" doesn't exist in the gateway.",
				"Please, remove or change url_path_map name then retry.",
			)
			return
		}
	}
}