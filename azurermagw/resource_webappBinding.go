package azurermagw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceWebappBindingType struct{}

type resourceWebappBinding struct {
	p provider
}

// Order Resource schema
func (r resourceWebappBindingType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": { // Containe the name of the Binding resource
				Type:     types.StringType,
				Required: true,
			},
			"agw_name": {
				Type:     types.StringType,
				Required: true,
			},
			"agw_rg": {
				Type:     types.StringType,
				Required: true,
			},
			"backend_address_pool": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"fqdns": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"ip_addresses": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
				}),
			},
			"backend_http_settings": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					//the affinity has a default value if it's not provided: "ApplicationGatewayAffinity"
					"affinity_cookie_name": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{stringDefault("ApplicationGatewayAffinity")},
					},
					"cookie_based_affinity": {
						Type:     types.StringType,
						Required: true,
					},
					"pick_host_name_from_backend_address": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
					"port": {
						Type:     types.Int64Type,
						Required: true,
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
					},
					"request_timeout": {
						Type:     types.Int64Type,
						Required: true,
					},
					"probe_name": {
						Type:     types.StringType,
						Optional: true,
					},
				}),
			},
			"probe": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"interval": {
						Type:     types.Int64Type,
						Required: true,
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
					},
					"path": {
						Type:     types.StringType,
						Required: true,
					},
					"pick_host_name_from_backend_http_settings": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
					"timeout": {
						Type:     types.Int64Type,
						Required: true,
					},
					"unhealthy_threshold": {
						Type:     types.Int64Type,
						Required: true,
					},	
					"minimum_servers": {
						Type:     types.Int64Type,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{intDefault(0)},
					},
					"match": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"body": {
								Type:     types.StringType,
								Required: true,
							},
							"status_code": {
								Type: types.ListType{
									ElemType: types.StringType,
								},
								Required: true,
							},
						}),
					},
				}),
			},
			"http_listener": {
				Optional: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"frontend_ip_configuration_name": {
						Type:     types.StringType,
						Required: true,
					},
					"frontend_port_name": {
						Type:     types.StringType,
						Required: true,
					},
					"require_sni": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
					},
					"host_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"host_names": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"ssl_certificate_name": {
						Type:     types.StringType,
						Optional: true,
					},
				},),
			},
			"https_listener": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"frontend_ip_configuration_name": {
						Type:     types.StringType,
						Required: true,
					},
					"frontend_port_name": {
						Type:     types.StringType,
						Required: true,
					},
					"require_sni": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
					},
					"host_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"host_names": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"ssl_certificate_name": {
						Type:     types.StringType,
						Optional: true,
					},
				},),
			},
			"ssl_certificate": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"key_vault_secret_id": {
						Type:     types.StringType,
						Optional: true,
					},
					"data": {
						Type:     types.StringType,
						Optional: true,
					},
					"password": {
						Type:     types.StringType,
						Optional: true,
						Sensitive: true,
					},
				}),
			},
			"redirect_configuration": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"redirect_type": {
						Type:     types.StringType,
						Required: true,
					},
					"target_listener_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"target_url": {
						Type:     types.StringType,
						Optional: true,
					},
					"include_path": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
					"include_query_string": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
					},
				}),
			},
			"request_routing_rule_http": {
				Optional: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"rule_type": {
						Type:     types.StringType,
						Required: true,
					},
					"priority": {
						Type:     types.StringType,
						Computed: true,
					},
					"http_listener_name": {
						Type:     types.StringType,
						Required: true,
					},
					"backend_address_pool_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"backend_http_settings_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"redirect_configuration_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"rewrite_rule_set_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"url_path_map_name": {
						Type:     types.StringType,
						Optional: true,
					},
				}),
			},
			"request_routing_rule_https": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"rule_type": {
						Type:     types.StringType,
						Required: true,
					},
					"priority": {
						Type:     types.StringType,
						Computed: true,
					},
					"http_listener_name": {
						Type:     types.StringType,
						Required: true,
					},
					"backend_address_pool_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"backend_http_settings_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"redirect_configuration_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"rewrite_rule_set_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"url_path_map_name": {
						Type:     types.StringType,
						Optional: true,
					},
				}),
			},
		},
	}, nil
}
// set default values
func intDefault(defaultValue int64) intDefaultModifier{
	return intDefaultModifier{
        Default: defaultValue,
    }
}
func stringDefault(defaultValue string) stringDefaultModifier {
	return stringDefaultModifier{
        Default: defaultValue,
    }
}
func boolDefault(defaultValue bool) boolDefaultModifier {
	return boolDefaultModifier{
        Default: defaultValue,
    }
}

// New resource instance
func (r resourceWebappBindingType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceWebappBinding{
		p: *(p.(*provider)),
	}, nil
}

// Create a new resource
func (r resourceWebappBinding) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	fmt.Println("\n######################## Create Method ########################")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource."+
			"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}
	
	// Retrieve values from plan
	var plan WebappBinding
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	//Get the agw (app gateway) from Azure with its Rest API
	resourceGroupName := plan.Agw_rg.Value
	applicationGatewayName := plan.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//Check if the agw already contains an existing element that has the same name of a new element to add
	exist_element, exist := checkElementName(gw, plan,plan.Http_listener,plan.Request_routing_rule_http)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. This (these) element(s) already exist(s) in the app gateway: \n"+ fmt.Sprint(exist_element),
			"Please, change its (their) name(s) then retry.",
		)
		return
	}
	
	//create, map and add the new elements (json) object from the plan to the agw object
	/************* generate and add BackendAddressPool **************/
	gw.Properties.BackendAddressPools = append(
		gw.Properties.BackendAddressPools, createBackendAddressPool(
			plan.Backend_address_pool))
	
	/************* generate and add request Routing Rule for HTTPS **************/
	if checkRequestRoutingRuleHttpsCreate(plan,gw,resp){
		return
	}
	priority := generatePriority(gw,"high")
	requestRoutingRuleHttps_json := createRequestRoutingRule(plan.Request_routing_rule_https,priority,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules,requestRoutingRuleHttps_json)
	
	/************* generate and add request Routing Rule for HTTP **************/
	if plan.Request_routing_rule_http != nil {
		if checkRequestRoutingRuleHttpCreate(plan,gw,resp){
			return
		}
		priority = generatePriority(gw,"high")
		requestRoutingRuleHttp_json := createRequestRoutingRule(plan.Request_routing_rule_http,priority,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules,requestRoutingRuleHttp_json)
	}

	/************* generate and add Backend HTTP Settings **************/
	if checkBackendHTTPSettingsCreate(plan,gw,resp){
		return
	}
	backendHTTPSettings_json := createBackendHTTPSettings(plan.Backend_http_settings,r.p.AZURE_SUBSCRIPTION_ID,
					resourceGroupName,applicationGatewayName)
	gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection,backendHTTPSettings_json)
	
	
	/************* generate and add probe **************/
	gw.Properties.Probes = append(gw.Properties.Probes,
		createProbe(plan.Probe,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName))

	/************* generate and add Http listener **************/
	// no ssl certificate to provide, so no need to check error_SslCertificateName
	if plan.Http_listener != nil {
		if checkHTTPListenerCreate(plan, gw, resp) {
			return
		}
		httpListener_json := createHTTPListener(plan.Http_listener,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpListener_json)
	}
	
	/************* generate and add Https listener **************/
	if checkHTTPSListenerCreate(plan, gw, resp) {
		return
	}
	httpsListener_json := createHTTPListener(plan.Https_listener,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)	
	gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpsListener_json)	
	
	/************* generate and add ssl Certificate **************/
	if checkSslCertificateCreate(plan, gw, resp) {
		return
	}
	sslCertificate_json := createSslCertificate(plan.Ssl_certificate,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	gw.Properties.SslCertificates = append(gw.Properties.SslCertificates,sslCertificate_json)

	/************* generate and add Redirect Configuration **************/
	if checkRedirectConfigurationCreate(plan, gw, resp) {
		return
	}
	redirectConfiguration_json:= createRedirectConfiguration(plan.Redirect_configuration,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	gw.Properties.RedirectConfigurations = append(gw.Properties.RedirectConfigurations,redirectConfiguration_json)


	//call the API to update the gw
	gw_response, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)
	
	//verify if the API response is 200 (that means, normaly, elements were added to the gateway), otherwise exit error
	if code != 200 {
		// Error  - backend address pool wasn't added to the app gateway
		resp.Diagnostics.AddError(
			"Unable to create the resource. ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Check the API response",
		)
		return
	}
	
	//generate the States based on gw_response from API.
	nb_Fqdns 		:= len(plan.Backend_address_pool.Fqdns)
	nb_IpAddress	:= len(plan.Backend_address_pool.Ip_addresses)
	backendAddressPool_state 		:= generateBackendAddressPoolState(gw_response,plan.Backend_address_pool.Name.Value,nb_Fqdns,nb_IpAddress)
	backendHTTPSettings_state 		:= generateBackendHTTPSettingsState(gw_response,plan.Backend_http_settings.Name.Value)
	probe_state 					:= generateProbeState(gw_response,plan.Probe.Name.Value)
	httpsListener_state 			:= generateHTTPListenerState(gw_response,plan.Https_listener.Name.Value)
	sslCertificate_state 			:= generateSslCertificateState(gw_response,plan.Ssl_certificate.Name.Value)
	redirectConfiguration_state 	:= generateRedirectConfigurationState(gw_response,plan.Redirect_configuration.Name.Value)
	requestRoutingRuleHttps_state	:= generateRequestRoutingRuleState(gw_response,plan.Request_routing_rule_https.Name.Value)

	//i moved "Generate resource state struct" with http listner block before it depends on the later.
	var result WebappBinding
	result = WebappBinding{
		Name						: plan.Name,
		Agw_name					: types.String{Value: gw_response.Name},
		Agw_rg						: plan.Agw_rg,
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Https_listener				: &httpsListener_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
		Request_routing_rule_https	: &requestRoutingRuleHttps_state,
	}
	//add Http_listener and Request_routing_rule_http if they are not nil
	if plan.Http_listener != nil {
		httpListener_state 	:= generateHTTPListenerState(gw_response,plan.Http_listener.Name.Value)
		result.Http_listener = &httpListener_state
	}else{
		result.Http_listener = nil
	}
	if plan.Request_routing_rule_http != nil {
		request_routing_rule_http_state := generateRequestRoutingRuleState(gw_response, plan.Request_routing_rule_http.Name.Value)
		result.Request_routing_rule_http = &request_routing_rule_http_state
	}else{
		result.Request_routing_rule_http = nil
	}

	//store to the created object to the terraform state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceWebappBinding) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	fmt.Println("\n######################## Read Method ########################")
	
	// Get current state
	var state WebappBinding
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get gw from API and then update what is in state from what the API returns
	webappBindingName := state.Name.Value

	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	// *********** Processing the request Routing Rule for HTTPS *********** //
	//check if the request Routing Rule exists in  the gateway, otherwise, it was removed manually
	var requestRoutingRuleHttps_state Request_routing_rule
	requestRoutingRuleHttpsName := state.Request_routing_rule_https.Name.Value
	//check if the backend http settings exists in the gateway, otherwise, it was removed manually
	if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttpsName) {
		//generate BackendState
		requestRoutingRuleHttps_state = generateRequestRoutingRuleState(gw, requestRoutingRuleHttpsName)
	}else{
		//generate an empty requestRoutingRuleHttps_state because it was removed manually
		requestRoutingRuleHttps_state = Request_routing_rule{}
	}
	
	// *********** Processing the backend address pool *********** //
	var backendAddressPool_state Backend_address_pool
	backendAddressPoolName := state.Backend_address_pool.Name.Value
	//check if the backend address pool exist in the gateway, otherwise, it was removed manually
	if checkBackendAddressPoolElement(gw, backendAddressPoolName) {
		// in the Read method, the number of fqdns and Ip in a Backendpool should be calculated from the json object and not the plan or state,
		// because the purpose of the read is to see if there is a difference between the real element and the satate stored localy.
		index := getBackendAddressPoolElementKey(gw, backendAddressPoolName)
		backendAddressPool_json := gw.Properties.BackendAddressPools[index]
		nb_BackendAddresses := len(backendAddressPool_json.Properties.BackendAddresses)
		nb_Fqdns := 0
		for i := 0; i < nb_BackendAddresses; i++ {
			if 	(backendAddressPool_json.Properties.BackendAddresses[i].Fqdn != "") && 
				(&backendAddressPool_json.Properties.BackendAddresses[i].Fqdn != nil) {
				nb_Fqdns++
			}
		}
		nb_IpAddress := nb_BackendAddresses - nb_Fqdns
		//generate BackendState
		backendAddressPool_state = generateBackendAddressPoolState(gw, backendAddressPoolName,nb_Fqdns,nb_IpAddress)
	}else{
		//generate an empty backendAddressPool_state because it was removed manually
		backendAddressPool_state = Backend_address_pool{}
	}
	
	// *********** Processing the backend http settings *********** //
	var backendHTTPSettings_state Backend_http_settings
	backendHTTPSettingsName := state.Backend_http_settings.Name.Value
	//check if the backend http settings exists in the gateway, otherwise, it was removed manually
	if checkBackendHTTPSettingsElement(gw, backendHTTPSettingsName) {
		//generate BackendState
		backendHTTPSettings_state = generateBackendHTTPSettingsState(gw, backendHTTPSettingsName)
	}else{
		//generate an empty backendHTTPSettings_state because it was removed manually
		backendHTTPSettings_state = Backend_http_settings{}
	}
	
	// *********** Processing the probe *********** //
	var probe_state Probe_tf
	probeName := state.Probe.Name.Value
	//check if the probe exists in the gateway, otherwise, it was removed manually
	if checkProbeElement(gw, probeName) {
		//generate probe state
		probe_state = generateProbeState(gw, probeName)
	}else{
		//generate an empty probe_state because it was removed manually
		probe_state = Probe_tf{}
	}

	// *********** Processing the https Listener *********** //
	//check if the Https listener  exists in  the gateway, otherwise, it was removed manually
	var httpsListener_state Http_listener
	httpsListenerName := state.Https_listener.Name.Value
	if checkHTTPListenerElement(gw, httpsListenerName) {
		httpsListener_state = generateHTTPListenerState(gw,httpsListenerName)
	}else{
		httpsListener_state = Http_listener{}
	}

	// *********** Processing the SSL Certificate *********** //
	//check if the SSL Certificate  exists in  the gateway, otherwise, it was removed manually
	var sslCertificate_state Ssl_certificate
	sslCertificateName := state.Ssl_certificate.Name.Value
	if checkSslCertificateElement(gw, sslCertificateName) {
		sslCertificate_state = generateSslCertificateState(gw,sslCertificateName)
	}else{
		sslCertificate_state = Ssl_certificate{}
	}

	// *********** Processing the Redirect Configuration *********** //
	var redirectConfiguration_state Redirect_configuration
	redirectConfigurationName := state.Redirect_configuration.Name.Value
	//check if the Redirect Configuration exists in the gateway, otherwise, it was removed manually
	if checkRedirectConfigurationElement(gw, redirectConfigurationName) {
		//generate BackendState
		redirectConfiguration_state = generateRedirectConfigurationState(gw, redirectConfigurationName)
	}else{
		//generate an empty redirectConfiguration_state because it was removed manually
		redirectConfiguration_state = Redirect_configuration{}
	}	
	
	var result WebappBinding
	result = WebappBinding{
		Name						: types.String{Value: webappBindingName},
		Agw_name					: state.Agw_name,
		Agw_rg						: state.Agw_rg,
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Https_listener				: &httpsListener_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
		Request_routing_rule_https	: &requestRoutingRuleHttps_state,
	}
	
	// *********** Processing the http Listener *********** //
	//check if the Http listener and request Routing Rule for HTTP exist in the old state (because they are optional param) 
	//in order to check if it's in the gateway, otherwise, it was removed manually
	if state.Http_listener != nil{
		httpListenerName := state.Http_listener.Name.Value
		if checkHTTPListenerElement(gw, httpListenerName) {
			httpListener_state 	:= generateHTTPListenerState(gw,httpListenerName)			
			result.Http_listener = &httpListener_state
		}else{
			result.Http_listener = nil
		}
	}else{
		result.Http_listener = nil
	}
	
	// *********** Processing the request Routing Rule for HTTP *********** //
	//check if the request Routing Rule exists in  the gateway, otherwise, it was removed manually
	if state.Request_routing_rule_http != nil {
		requestRoutingRuleHttpName := state.Request_routing_rule_http.Name.Value
		if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttpName) {
			//generate State
			requestRoutingRuleHttp_state := generateRequestRoutingRuleState(gw, requestRoutingRuleHttpName)
			result.Request_routing_rule_http = &requestRoutingRuleHttp_state
		}else{
			//generate an empty requestRoutingRuleHttp_state because it was removed manually
			result.Request_routing_rule_http = nil
		}
	}else{
		result.Request_routing_rule_http = nil // &Request_routing_rule{}//nil
	}
	fmt.Println("\n######################## after Routing Rule for HTTP ########################\n", result)
	
	state = result
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
// Update resource
func (r resourceWebappBinding) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	fmt.Println("\n######################## Update Method ########################")
	// Get plan values
	var plan WebappBinding
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state WebappBinding
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	//Get the agw in order to update it with new values from plan
	resourceGroupName := plan.Agw_rg.Value
	applicationGatewayName := plan.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)

	//for all elements (attributes), prepare the new elements (json) from the plan
	//Verify if the agw already contains the elements to be updated beacause:
	//		- the older ones has be removed before updating. 
	//		- we have also to prevent element name updating and manual deletion

	// *********** Processing backend address pool *********** //	
	//preparing the new elements (json) from the plan
	backendAddressPool_plan := plan.Backend_address_pool
	backendAddressPool_json := createBackendAddressPool(backendAddressPool_plan)
	
	//check if the backend name in the plan and state are different, that means that
	//it's about backend AddressPool update  with the same name
	if backendAddressPool_plan.Name.Value == state.Backend_address_pool.Name.Value {
		//so we remove the old one before adding the new one.
		removeBackendAddressPoolElement(&gw, backendAddressPool_json.Name)
	}else{
		// it's most likely about backend update with a new name
		// we have to check if the new backend name is already used
		if checkBackendAddressPoolElement(gw, backendAddressPool_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new Backend Adresse pool name : "+ backendAddressPool_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old backend (old name) from the gateway
		removeBackendAddressPoolElement(&gw, state.Backend_address_pool.Name.Value)
	}
	
	var priority int
	// *********** Processing request Routing Rule (HTTP) *********** //	
	if plan.Request_routing_rule_http != nil {	
		if checkRequestRoutingRuleHttpUpdate(plan,gw,resp) {
			return
		}
		//to compute priority, check if Request Routing Rule exist in the state, so we get the old priority
		// else, that means the old Request Routing Rule was removed manually, we have to generate a new priority
		//var priority int
		if state.Request_routing_rule_http != nil {
			if state.Request_routing_rule_http.Priority.Value != "0" {
				priority,_ = strconv.Atoi(state.Request_routing_rule_http.Priority.Value)
			}else{
				priority = generatePriority(gw,"high")
			}
		}else{
			priority = generatePriority(gw,"high")
		}
		requestRoutingRuleHttp_json := createRequestRoutingRule(plan.Request_routing_rule_http,priority,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		
		//new Request routing rule is ok. now we have to remove the old one
		if state.Request_routing_rule_http != nil {		
			requestRoutingRuleHttp_plan := plan.Request_routing_rule_http
			if requestRoutingRuleHttp_plan.Name.Value == state.Request_routing_rule_http.Name.Value {
				//so we remove the old one before adding the new one.
				removeRequestRoutingRuleElement(&gw,requestRoutingRuleHttp_json.Name)
			}else{
				// it's most likely about Request routing rule update with a new name
				// we have to check if the new Request routing rule name is already used
				if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttp_json.Name) {
					//this is an error. issue an exit error.
					resp.Diagnostics.AddError(
						"Unable to update the app gateway. The new Request routing rule (HTTP) name : "+ requestRoutingRuleHttp_json.Name+" already exists.",
						" Please, change the name then retry.",)
					return
				}
				//remove the old  Request routing rule (old name) from the gateway
				removeRequestRoutingRuleElement(&gw, state.Request_routing_rule_http.Name.Value)
			}
		}
		//we have to add the Request Routing Rule (HTTP) here because it's optional
		gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules, requestRoutingRuleHttp_json)
	}else{//plan.Request_routing_rule_http = nil, no Request_routing_rule_http in the plan (ex: removed)
		//check if there is an old one in the state in order to remove it
		if state.Request_routing_rule_http != nil {
			removeRequestRoutingRuleElement(&gw, state.Request_routing_rule_http.Name.Value)
		}
	}

	// *********** Processing request Routing Rule (HTTPS) *********** //	
	if checkRequestRoutingRuleHttpsUpdate(plan,gw,resp) {
		return
	}
	//to compute priority, check if Request Routing Rule exist in the state, so we get the old priority
	// else, that means the old Request Routing Rule was removed manually, we have to generate a new priority
	
	if state.Request_routing_rule_https.Priority.Value != "0" {
		//the priority of new Request_routing_rule_http is already included in gw, so it's ok
		priority,_ = strconv.Atoi(state.Request_routing_rule_https.Priority.Value)
	}else{
		priority = generatePriority(gw,"high")
	}
	requestRoutingRuleHttps_json := createRequestRoutingRule(plan.Request_routing_rule_https,priority,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	
	//new Request routing rule is ok. now we have to remove the old one
	requestRoutingRuleHttps_plan := plan.Request_routing_rule_https
	if requestRoutingRuleHttps_plan.Name.Value == state.Request_routing_rule_https.Name.Value {
		//so we remove the old one before adding the new one.
		removeRequestRoutingRuleElement(&gw,requestRoutingRuleHttps_json.Name)
	}else{
		// it's most likely about Request routing rule update with a new name
		// we have to check if the new Request routing rule name is already used
		if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttps_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new Request routing rule (HTTPS) name : "+ requestRoutingRuleHttps_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old Request routing rule (old name) from the gateway
		removeRequestRoutingRuleElement(&gw, state.Request_routing_rule_https.Name.Value)
	}	
	
	// *********** Processing backend http settings *********** //	
	if checkBackendHTTPSettingsUpdate(plan,gw,resp){
		return
	}
	//preparing the new elements (json) from the plan
	backendHTTPSettings_plan := plan.Backend_http_settings
	backendHTTPSettings_json:= createBackendHTTPSettings(backendHTTPSettings_plan,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	
	//check if the backend HTTPSettings name in the plan and state are different, that means that
	//it's about backend HTTPSettings update  with the same name
	if backendHTTPSettings_plan.Name.Value == state.Backend_http_settings.Name.Value {
		//it's about backend http settings update  with the same name
		//so we remove the old one before adding the new one.
		removeBackendHTTPSettingsElement(&gw, backendHTTPSettings_json.Name)
	}else{
		// it's about backend http settings update with a new name
		// we have to check if the new backend http settings name is already used
		if checkBackendHTTPSettingsElement(gw, backendHTTPSettings_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new Backend HTTP settings name : "+ backendHTTPSettings_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old backend http settings (old name) from the gateway
		removeBackendHTTPSettingsElement(&gw, state.Backend_http_settings.Name.Value)
	}

	// *********** Processing the probe *********** //	
	//preparing the new elements (json) from the plan
	probe_plan := plan.Probe	
	probe_json := createProbe(probe_plan,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)

	//check if the probe name in the plan and state are different,that means that
	//it's about probe update  with the same name
	if probe_plan.Name.Value == state.Probe.Name.Value {
		//so we remove the old one before adding the new one.
		removeProbeElement(&gw, probe_json.Name)
	}else{
		// it's about probe update with a new name
		// we have to check if the new probe name is already used
		if checkProbeElement(gw, probe_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new probe name : "+ probe_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old backend http settings (old name) from the gateway
		removeProbeElement(&gw, state.Probe.Name.Value)
	}

	// *********** Processing http Listener *********** //	
	//preparing the new elements (json) from the plan
	if plan.Http_listener != nil {
		if checkHTTPListenerUpdate(plan, gw, resp) {
			return
		}
		httpListener_plan := plan.Http_listener
		httpListener_json := createHTTPListener(httpListener_plan,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		//new http listener is ok. now we have to remove the old one if there is already one, else, nothing to remove
		if state.Http_listener != nil {
			if httpListener_plan.Name.Value == state.Http_listener.Name.Value {
				//so we remove the old one before adding the new one.
				removeHTTPListenerElement(&gw, httpListener_json.Name)
			}else{
				// it's most likely about http Listener update with a new name
				// we have to check if the new http Listener name is already used
				if checkHTTPListenerElement(gw, httpListener_json.Name) {
					//this is an error. issue an exit error.
					resp.Diagnostics.AddError(
						"Unable to update the app gateway. The new http Listener name : "+ httpListener_json.Name+" already exists.",
						" Please, change the name then retry.",
					)
					return
				}
				//remove the old http Listener (old name) from the gateway
				removeHTTPListenerElement(&gw, state.Http_listener.Name.Value)
			}
		}
		//we have to add the http listener here because it's optional
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners, httpListener_json)			
	}else{//plan.Http_listener = nil, no http listener in the plan (ex: removed)
		//check if there is an old one in the state in order to remove it
		if state.Http_listener != nil {
			removeHTTPListenerElement(&gw, state.Http_listener.Name.Value)
		}
	}
			
	// *********** Processing https Listener *********** //	
	//preparing the new elements (json) from the plan
	if checkHTTPSListenerUpdate(plan, gw, resp) {
		return
	}
	httpsListener_plan := plan.Https_listener
	httpsListener_json := createHTTPListener(httpsListener_plan,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	//new https listener is ok. now we have to remove the old one
	if httpsListener_plan.Name.Value == state.Https_listener.Name.Value {
		//so we remove the old one before adding the new one.
		removeHTTPListenerElement(&gw, httpsListener_json.Name)
	}else{
		// it's most likely about https Listener update with a new name
		// we have to check if the new https Listener name is already used
		if checkHTTPListenerElement(gw, httpsListener_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new https Listener name : "+ httpsListener_json.Name+" already exists."+
				"Could be due to the name of the http listener you are under declaring",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old http Listener (old name) from the gateway
		removeHTTPListenerElement(&gw, state.Https_listener.Name.Value)
	}

	// *********** Processing SSL Certificate *********** //	
	//preparing the new elements (json) from the plan
	if checkSslCertificateUpdate(plan, gw, resp){
		return
	}
	sslCertificate_plan := plan.Ssl_certificate
	sslCertificate_json := createSslCertificate(plan.Ssl_certificate,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	
	//check if the SSL Certificate name in the plan and state are different, that means that
	//it's about SSL Certificate update  with the same name
	if sslCertificate_plan.Name.Value == state.Ssl_certificate.Name.Value {
		//it's about SSL Certificate update  with the same name
		//so we remove the old one before adding the new one.
		removeSslCertificateElement(&gw, sslCertificate_json.Name)
	}else{
		// it's about SSL Certificate update with a new name
		// we have to check if the new SSL Certificate name is already used
		if checkSslCertificateElement(gw, sslCertificate_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new SSL Certificate name : "+ sslCertificate_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old SSL Certificate (old name) from the gateway
		removeSslCertificateElement(&gw, state.Ssl_certificate.Name.Value)
	}

	// *********** Processing Redirect Configuration *********** //	
	//preparing the new element (json) from the plan
	if checkRedirectConfigurationUpdate(plan,gw,resp) {
		return
	}
	redirectConfiguration_plan := plan.Redirect_configuration
	redirectConfiguration_json := createRedirectConfiguration(redirectConfiguration_plan,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	
	//check if the Redirect Configuration name in the plan and state are different, that means that
	//it's about Redirect Configuration update  with the same name
	if redirectConfiguration_plan.Name.Value == state.Redirect_configuration.Name.Value {
		//it's about Redirect Configuration update  with the same name
		//so we remove the old one before adding the new one.
		removeRedirectConfigurationElement(&gw, redirectConfiguration_json.Name)
	}else{
		// it's about Redirect Configuration update with a new name
		// we have to check if the new Redirect Configuration name is already used
		if checkRedirectConfigurationElement(gw, redirectConfiguration_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new Redirect Configuration name : "+ redirectConfiguration_json.Name+" already exists.",
				" Please, change the name then retry.",
			)
			return
		}
		//remove the old Redirect Configuration (old name) from the gateway
		removeRedirectConfigurationElement(&gw, state.Redirect_configuration.Name.Value)
	}

	//add the new elements (http Listener and Request Routing Rule (HTTP) elements are already added because they are optionals). 
	gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, backendAddressPool_json)
	gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection, backendHTTPSettings_json)
	gw.Properties.Probes = append(gw.Properties.Probes, probe_json)
	gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners, httpsListener_json)
	gw.Properties.SslCertificates = append(gw.Properties.SslCertificates, sslCertificate_json)
	gw.Properties.RedirectConfigurations = append(gw.Properties.RedirectConfigurations, redirectConfiguration_json)
	gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules,requestRoutingRuleHttps_json)


	//and update the gateway
	gw_response, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)

	//verify if the API response is 200 (that means, normaly, elements were added to the gateway), otherwise exit error
	if code != 200 {
		// Error  - when adding new elements to the app gateway
		resp.Diagnostics.AddError(
			"Unable to update the resource. ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Check the API response",
		)
		return
	}

	// Generate new states 
	/*********** Special for Backend Address Pool ********************/
	// in the Read method, the number of fqdns and Ip in a Backendpool should be calculated from the json object and not the plan or state,
	// because the purpose of the read is to see if there is a difference between the real element and the satate stored localy.
	index := getBackendAddressPoolElementKey(gw, backendAddressPool_json.Name)
	backendAddressPool_json2 := gw.Properties.BackendAddressPools[index]
	nb_BackendAddresses := len(backendAddressPool_json2.Properties.BackendAddresses)
	nb_Fqdns := 0
	for i := 0; i < nb_BackendAddresses; i++ {
		if (backendAddressPool_json2.Properties.BackendAddresses[i].Fqdn != "") && (&backendAddressPool_json2.Properties.BackendAddresses[i].Fqdn != nil) {
			nb_Fqdns++
		} 
	}
	nb_IpAddress := nb_BackendAddresses - nb_Fqdns
	/*****************************************************************/
	
	backendAddressPool_state		:= generateBackendAddressPoolState(gw_response, backendAddressPool_json.Name,nb_Fqdns,nb_IpAddress)
	backendHTTPSettings_state		:= generateBackendHTTPSettingsState(gw_response,backendHTTPSettings_json.Name)
	probe_state						:= generateProbeState(gw_response,probe_json.Name)
	httpsListener_state 			:= generateHTTPListenerState(gw_response,httpsListener_json.Name)
	sslCertificate_state 			:= generateSslCertificateState(gw_response,sslCertificate_json.Name)
	redirectConfiguration_state 	:= generateRedirectConfigurationState(gw_response,redirectConfiguration_json.Name)
	requestRoutingRuleHttps_state 	:= generateRequestRoutingRuleState(gw_response,requestRoutingRuleHttps_json.Name)

	/*************** Special for Http listener **********************/
	// Generate resource state struct 
	//i moved "Generate resource state struct" with http listner block before it depends on the later.
	
	var result WebappBinding
	result = WebappBinding{
		Name						: state.Name,
		Agw_name					: types.String{Value: gw_response.Name},
		Agw_rg						: state.Agw_rg,
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Https_listener				: &httpsListener_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
		Request_routing_rule_https	: &requestRoutingRuleHttps_state,
	}
	if plan.Http_listener != nil {
		httpListener_state 	:= generateHTTPListenerState(gw_response,plan.Http_listener.Name.Value)
		result.Http_listener = &httpListener_state
	}else{
		result.Http_listener = nil
	}
	if plan.Request_routing_rule_http != nil {
		requestRoutingRuleHttp_state 	:= generateRequestRoutingRuleState(gw_response,plan.Request_routing_rule_http.Name.Value)
		result.Request_routing_rule_http = &requestRoutingRuleHttp_state
	}else{
		result.Request_routing_rule_http = nil
	}
	
	//store to the created objecy to the terraform state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceWebappBinding) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	fmt.Println("\n######################## Delete Method ########################")
	// Get current state
	var state WebappBinding
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get elements names from state
	backendAddressPoolName 		:= state.Backend_address_pool.Name.Value
	backendHTTPSettingsName 	:= state.Backend_http_settings.Name.Value
	probeName 					:= state.Probe.Name.Value
	httpsListenerName 			:= state.Https_listener.Name.Value
	sslCertificateName 			:= state.Ssl_certificate.Name.Value
	redirectConfigurationName 	:= state.Redirect_configuration.Name.Value
	requestRoutingRuleHttpsName := state.Request_routing_rule_https.Name.Value

	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//remove the elements from the gw
	removeBackendAddressPoolElement(&gw, backendAddressPoolName)
	removeBackendHTTPSettingsElement(&gw,backendHTTPSettingsName)
	removeProbeElement(&gw,probeName)
	removeHTTPListenerElement(&gw,httpsListenerName)
	removeSslCertificateElement(&gw,sslCertificateName)
	removeRedirectConfigurationElement(&gw,redirectConfigurationName)
	removeRequestRoutingRuleElement(&gw,requestRoutingRuleHttpsName)

	/*************** Special for Http listener **********************/
	if state.Http_listener != nil {
		removeHTTPListenerElement(&gw,state.Http_listener.Name.Value)
	}		
	/*************** Special for Http listener **********************/
	if state.Request_routing_rule_http != nil {
		removeRequestRoutingRuleElement(&gw,state.Request_routing_rule_http.Name.Value)
	}
	
	//and update the gateway
	_, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)
	//verify if the API response is 200 (that means, normaly, elements were deleted to the gateway), otherwise exit error
	if code != 200 {
		// Error  - when deleting new elements to the app gateway
		resp.Diagnostics.AddError(
			"Unable to delete the resource. ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Check the API response",
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceWebappBinding) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	//tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
func checkElementName(gw ApplicationGateway, plan WebappBinding,httpListener_plan *Http_listener,requestRoutingRuleHttp_plan *Request_routing_rule) ([]string,bool){
	//This function allows to check if an element name in the required new configuration (plan WebappBinding) already exist in the gw.
	//if so, the provider has to stop executing and issue an exit error
	exist := false
	var existing_element_list [] string
	//Create new var for all configurations
	backendAddressPool_plan 	:= plan.Backend_address_pool 
	backendHTTPSettings_plan 	:= plan.Backend_http_settings
	probe_plan 					:= plan.Probe
	httpsListener_plan 			:= plan.Https_listener
	sslCertificate_plan			:= plan.Ssl_certificate
	redirectConfiguration_plan	:= plan.Redirect_configuration
	requestRoutingRuleHttps_plan:= plan.Request_routing_rule_https

	if checkBackendAddressPoolElement(gw, backendAddressPool_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- BackendAddressPool: "+backendAddressPool_plan.Name.Value)
	}
	if checkBackendHTTPSettingsElement(gw, backendHTTPSettings_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- BackendHTTPSettings: "+backendHTTPSettings_plan.Name.Value)
	}
	if checkProbeElement(gw, probe_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- Probe: "+probe_plan.Name.Value)
	}
	if checkHTTPListenerElement(gw, httpsListener_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- HTTPListener: "+httpsListener_plan.Name.Value)
	}
	if httpListener_plan != nil {
		//httpListener_plan 			:= httpListener_plan
		if checkHTTPListenerElement(gw, httpListener_plan.Name.Value) {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- HTTPListener: "+httpListener_plan.Name.Value)
		}
		if httpListener_plan.Name.Value == httpsListener_plan.Name.Value {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- HTTP and HTTPS Listener (new): "+httpListener_plan.Name.Value)
		}
	}
	if checkSslCertificateElement(gw, sslCertificate_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- SSL Certificate: "+sslCertificate_plan.Name.Value)
	}
	if checkRedirectConfigurationElement(gw, redirectConfiguration_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- Redirect configuration: "+redirectConfiguration_plan.Name.Value)
	}
	if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttps_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- Request Routing Rule for HTTPS: "+requestRoutingRuleHttps_plan.Name.Value)
	}
	if requestRoutingRuleHttp_plan != nil {
		if checkRequestRoutingRuleElement(gw, requestRoutingRuleHttp_plan.Name.Value) {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- Request Routing Rule for HTTP: "+requestRoutingRuleHttp_plan.Name.Value)
		}
		if requestRoutingRuleHttp_plan.Name.Value == requestRoutingRuleHttps_plan.Name.Value {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- Request Routing Rule for HTTP and HTTPS (new): "+requestRoutingRuleHttp_plan.Name.Value)
		}
	}
	return existing_element_list,exist
}
func generatePriority(gw ApplicationGateway, level string) int {
	priority := 0
	rand.Seed(time.Now().UnixNano())
	var priorities = make([]int,len(gw.Properties.RequestRoutingRules))
	for i := 0; i < len(gw.Properties.RequestRoutingRules); i++ {
		priorities[i] = gw.Properties.RequestRoutingRules[i].Properties.Priority
	}
	good_priority :=false
	for good_priority == false{
		priority = rand.Intn(300-1) + 1
		for i := 0; i < len(priorities); i++ {
			if priority == priorities[i] {
				good_priority = true
			}
		}
		good_priority = !good_priority
	}
	return priority
}
//Client operations
func getGW(subscriptionId string, resourceGroupName string, applicationGatewayName string, token string) ApplicationGateway {
	requestURI := "https://management.azure.com/subscriptions/" + subscriptionId + "/resourceGroups/" +
		resourceGroupName + "/providers/Microsoft.Network/applicationGateways/" + applicationGatewayName + "?api-version=2021-08-01"
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Call failure: %+v", err)
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var agw ApplicationGateway
	err = json.Unmarshal(responseData, &agw)

	if err != nil {
		fmt.Println(err)
	}
	return agw
}
func updateGW(subscriptionId string, resourceGroupName string, applicationGatewayName string, gw ApplicationGateway, token string) (ApplicationGateway, string, int) {
	requestURI := "https://management.azure.com/subscriptions/" + subscriptionId + "/resourceGroups/" +
		resourceGroupName + "/providers/Microsoft.Network/applicationGateways/" + applicationGatewayName + "?api-version=2021-08-01"
	payloadBytes, err := json.Marshal(gw)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", requestURI, body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Call failure: %+v", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//if code != 200, the responseData contain a json that describe the error
	rs := string(responseData)
	error_json, err := PrettyString(rs)
	if err != nil {
		log.Fatal(err)
	}
	var agw ApplicationGateway
	err = json.Unmarshal(responseData, &agw)
	if err != nil {
		fmt.Println(err)
	}
	return agw, error_json, code
}

//some debugging tools
func PrettyStringGW(gw ApplicationGateway) string {
	payloadBytes, err := json.Marshal(gw)
	if err != nil {
		// handle err
	}
	str := string(payloadBytes)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "error"
	}
	return prettyJSON.String()
}
func PrettyStringFromByte(str []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
func printToFile(str string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	fmt.Fprintln(mw, str)
}
func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
	  rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
	  return false
	}
	return rv.FieldByName(name).IsValid()
}