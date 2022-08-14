package azurermagw

import (
	//"fmt"
	"fmt"
	//"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Probe_json struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState                   string `json:"provisioningState,omitempty"`
		Protocol                            string `json:"protocol,omitempty"`
		Path                                string `json:"path,omitempty"`
		Interval                            int    `json:"interval,omitempty"`
		Timeout                             int    `json:"timeout,omitempty"`
		UnhealthyThreshold                  int    `json:"unhealthyThreshold,omitempty"`
		PickHostNameFromBackendHTTPSettings bool   `json:"pickHostNameFromBackendHttpSettings,omitempty"`
		MinServers                          int    `json:"minServers,omitempty"`
		Match                               *struct {
			Body        string   `json:"body,omitempty"`
			StatusCodes []string `json:"statusCodes,omitempty"`
		} `json:"match"`
		Host                                string `json:"host,omitempty"`
		BackendHTTPSettings *[]struct {
			ID string `json:"id,omitempty"`
		} `json:"backendHttpSettings"`
	} `json:"properties"`
	Type string `json:"type,omitempty"`
} 

type Probe_tf struct {
	Name         								types.String	`tfsdk:"name"`	
	Id           								types.String	`tfsdk:"id"`
	Interval           							types.Int64		`tfsdk:"interval"`		
	Protocol                       				types.String	`tfsdk:"protocol"`								
	Path          								types.String	`tfsdk:"path"`		
	Timeout										types.Int64		`tfsdk:"timeout"`						
	Unhealthy_threshold							types.Int64		`tfsdk:"unhealthy_threshold"`
	Pick_host_name_from_backend_http_settings 	types.Bool		`tfsdk:"pick_host_name_from_backend_http_settings"`	
	Minimum_servers								types.Int64		`tfsdk:"minimum_servers"`
	//i choose to make Match param as required to avoid Value Conversion Error in terraform when 
	//this param is not provided
	Match										Match			`tfsdk:"match"`
	//the BackendHTTPSettings is added automatically when we provide probe name in BackendHTTPSettings params
}
type Match	struct{
	Body          								types.String	`tfsdk:"body"`
	Status_code									[]types.String	`tfsdk:"status_code"`
}

func createProbe(probe_plan Probe_tf,AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) Probe_json{	
	probe_json := Probe_json{
		Name:       probe_plan.Name.Value,
		//ID:         "",
		//Etag:       "",
		Properties: struct{
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			Protocol string "json:\"protocol,omitempty\""; 
			Path string "json:\"path,omitempty\""; 
			Interval int "json:\"interval,omitempty\""; 
			Timeout int "json:\"timeout,omitempty\""; 
			UnhealthyThreshold int "json:\"unhealthyThreshold,omitempty\""; 
			PickHostNameFromBackendHTTPSettings bool "json:\"pickHostNameFromBackendHttpSettings,omitempty\""; 
			MinServers int "json:\"minServers,omitempty\""; 
			Match *struct{
				Body string "json:\"body,omitempty\""; 
				StatusCodes []string "json:\"statusCodes,omitempty\""
			} "json:\"match\""; 
			Host string "json:\"host,omitempty\""; 
			BackendHTTPSettings *[]struct{
				ID string "json:\"id,omitempty\""
			} "json:\"backendHttpSettings\""
		}{
			Interval: int(probe_plan.Interval.Value),
			Protocol: probe_plan.Protocol.Value,
			Path: probe_plan.Path.Value,
			Timeout: int(probe_plan.Timeout.Value),
			UnhealthyThreshold: int(probe_plan.Unhealthy_threshold.Value),
			PickHostNameFromBackendHTTPSettings: bool(probe_plan.Pick_host_name_from_backend_http_settings.Value),
			MinServers: int(probe_plan.Minimum_servers.Value),
		},
		Type: "Microsoft.Network/applicationGateways/probes",
	}

	//it remains Match struct. we have to check if it is provided or not (actually, i don't have time)
	
	if &probe_plan.Match != nil {
		probe_json.Properties.Match = &struct{
			Body string "json:\"body,omitempty\""; 
			StatusCodes []string "json:\"statusCodes,omitempty\""
			}{
				Body: probe_plan.Match.Body.Value,
				//StatusCodes: ,
			}
		probe_json.Properties.Match.StatusCodes = make([]string, len(probe_plan.Match.Status_code))
		for i := 0; i < len(probe_plan.Match.Status_code); i++ {
			probe_json.Properties.Match.StatusCodes[i] = probe_plan.Match.Status_code[i].Value
		}
	}else{
		probe_json.Properties.Match = nil
	}		
	
	return probe_json
}
func generateProbeState(gw ApplicationGateway, ProbeName string) Probe_tf {
	//retrieve json element from gw
	index := getProbeElementKey(gw, ProbeName)
	probe_json := gw.Properties.Probes[index]
	
	// Map response body to resource schema attribute
	var probe_state Probe_tf
	probe_state = Probe_tf{
		Name         								: types.String {Value: probe_json.Name},
		Id           								: types.String {Value: probe_json.ID},
		Interval           							: types.Int64 {Value: int64(probe_json.Properties.Interval)},	
		Protocol                       				: types.String {Value: probe_json.Properties.Protocol},
		Path          								: types.String {Value: probe_json.Properties.Path},
		Timeout										: types.Int64	{Value: int64(probe_json.Properties.Timeout)},
		Unhealthy_threshold							: types.Int64	{Value: int64(probe_json.Properties.UnhealthyThreshold)},
		Pick_host_name_from_backend_http_settings 	: types.Bool {Value: bool(probe_json.Properties.PickHostNameFromBackendHTTPSettings)},
		Minimum_servers								: types.Int64	{Value: int64(probe_json.Properties.MinServers)},
		Match	: Match {
					Body		: types.String{Value: probe_json.Properties.Match.Body},
					Status_code	: []types.String{},
				},
	}

	if len(probe_json.Properties.Match.StatusCodes) != 0 {
		probe_state.Match.Status_code = make([]types.String,len(probe_json.Properties.Match.StatusCodes) )
	} else {
		probe_state.Match.Status_code = nil
	}
	for i := 0; i < len(probe_json.Properties.Match.StatusCodes); i++ {
		probe_state.Match.Status_code[i]=types.String{Value: probe_json.Properties.Match.StatusCodes[i]}
	}
	//verify if optional parameters are provided, otherwise, they have to set to null
		
	return probe_state
}
func getProbeElementKey(gw ApplicationGateway, ProbeName string) int {
	key := -1
	for i := len(gw.Properties.Probes) - 1; i >= 0; i-- {
		if gw.Properties.Probes[i].Name == ProbeName {
			key = i
		}
	}
	return key
}
func checkProbeElement(gw ApplicationGateway, ProbeName string) bool {
	exist := false
	for i := len(gw.Properties.Probes) - 1; i >= 0; i-- {
		if gw.Properties.Probes[i].Name == ProbeName {
			exist = true
		}
	}
	return exist
}
func removeProbeElement(gw *ApplicationGateway, ProbeName string) {
	for i := len(gw.Properties.Probes) - 1; i >= 0; i-- {
		if gw.Properties.Probes[i].Name == ProbeName {
			gw.Properties.Probes = append(gw.Properties.Probes[:i], gw.Properties.Probes[i+1:]...)
		}
	}
}