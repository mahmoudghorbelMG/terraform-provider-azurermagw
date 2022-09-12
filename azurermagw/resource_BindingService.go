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

type resourceBindingServiceType struct{}

type resourceBindingService struct {
	p provider
}

// Order Resource schema
func (r resourceBindingServiceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": { // Containe the name of the Binding resource
				Type:     types.StringType,
				Required: true,
				//Description: "This is a description message",
				MarkdownDescription: "The name of the binding service that bind an backend application (VM, web app, container web app, etc.) to the azure application gateway.",
			},
			"application_gateway_name": {
				Type:     types.StringType,
				Required: true,
				MarkdownDescription: "The name of the application gateway to which the backend application will be binded.",
			},
			"application_gateway_resource_group_name": {
				Type:     types.StringType,
				Required: true,
				MarkdownDescription: "The name of the resource group where the application gateway is deployed.",
			},
			"backend_address_pool": {
				Required: true,
				MarkdownDescription: "For this provider version, only one `backend_address_pool` block can be set as defined below.",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The name of the Backend Address Pool.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `backend_address_pool`.",
					},
					"fqdns": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
						MarkdownDescription: "A list of FQDN's which should be part of the Backend Address Pool.",
					},
					"ip_addresses": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
						MarkdownDescription: "A list of IP Addresses which should be part of the Backend Address Pool.",
					},
				}),
			},
			"backend_http_settings": {
				Required: true,
				MarkdownDescription: "For this provider version, only one `backend_http_settings` block can be set as defined below.",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The name of the Backend HTTP Settings Collection.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `backend_http_settings`.",
					},
					//the affinity has a default value if it's not provided: "ApplicationGatewayAffinity"
					"affinity_cookie_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The name of the affinity cookie. Required if `cookie_based_affinity` is `Enabled`",
						/*Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{stringDefault("ApplicationGatewayAffinity")},*/
					},
					"cookie_based_affinity": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "Is Cookie-Based Affinity enabled? Possible values are `Enabled` and `Disabled`.",
					},
					"pick_host_name_from_backend_address": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
						MarkdownDescription: "Whether host header should be picked from the host name of the backend server. Defaults to `false`.",
					},
					"port": {
						Type:     types.Int64Type,
						Required: true,
						MarkdownDescription: "The port which should be used for this Backend HTTP Settings Collection.",
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Protocol which should be used. Possible values are `Http` and `Https`.",
					},
					"request_timeout": {
						Type:     types.Int64Type,
						Required: true,
						MarkdownDescription: "The request timeout in seconds, which must be between 1 and 86400 seconds.",
					},
					"probe_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The name of an associated HTTP Probe.",
					},
				}),
			},
			"probe": {
				Required: true,
				MarkdownDescription: "For this provider version, only one `probe` block can be set as defined below",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: " The Name of the Probe.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `probe`.",
					},
					"interval": {
						Type:     types.Int64Type,
						Required: true,
						MarkdownDescription: "The Interval between two consecutive probes in seconds. Possible values range from 1 second to a maximum of 86,400 seconds.",
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Protocol used for this Probe. Possible values are `Http` and `Https`.",
					},
					"path": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: " The Path used for this Probe.",
					},
					"pick_host_name_from_backend_http_settings": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
						MarkdownDescription: "Whether the host header should be picked from the backend HTTP settings. Defaults to `false`.",
					},
					"timeout": {
						Type:     types.Int64Type,
						Required: true,
						MarkdownDescription: "The Timeout used for this Probe, which indicates when a probe becomes unhealthy. Possible values range from 1 second to a maximum of 86,400 seconds.",
					},
					"unhealthy_threshold": {
						Type:     types.Int64Type,
						Required: true,
						MarkdownDescription: "The Unhealthy Threshold for this Probe, which indicates the amount of retries which should be attempted before a node is deemed unhealthy. Possible values are from 1 to 20.",
					},	
					"minimum_servers": {
						Type:     types.Int64Type,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{intDefault(0)},
						MarkdownDescription: "The minimum number of servers that are always marked as healthy. Defaults to `0`.",
					},
					"match": {
						Required: true,
						MarkdownDescription: "A `match` block as defined above.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"body": {
								Type:     types.StringType,
								Required: true,
								MarkdownDescription: "A snippet from the Response Body which must be present in the Response.",
							},
							"status_code": {
								Type: types.ListType{
									ElemType: types.StringType,
								},
								Required: true,
								MarkdownDescription: "A list of allowed status codes for this Health Probe.",
							},
						}),
					},
				}),
			},			
			"ssl_certificate": {
				Required: true,
				MarkdownDescription: "For this provider version, only one `ssl_certificate` block can be set as defined below",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of the SSL certificate that is unique within this Application Gateway",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `ssl_certificate`.",
					},
					"key_vault_secret_id": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "Secret Id of (base-64 encoded unencrypted pfx) `Secret` or `Certificate` object stored in Azure KeyVault. "+
						"You need to enable soft delete for keyvault to use this feature. Required if `data` is not set. "+
						"TLS termination with Key Vault certificates is limited to the v2 SKUs.",
					},
					"data": {
						Type:     types.StringType,
						Optional: true,
						Sensitive: true,
						MarkdownDescription: "PFX certificate. Required if `key_vault_secret_id` is not set. Important: PFX certificate is not supported in the current provider version.",
					},
					"password": {
						Type:     types.StringType,
						Optional: true,
						Sensitive: true,
						MarkdownDescription: "Password for the pfx file specified in data. Required if `data` is set.",
					},/*
					"public_cert_data": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "",
					},*/
				}),
			},
			"redirect_configuration": {
				Required: true,
				MarkdownDescription: "For this provider version, only one `redirect_configuration` block can be set as defined below",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "Unique name of the redirect configuration block.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `redirect_configuration`.",
					},
					"redirect_type": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The type of redirect. Possible values are `Permanent`, `Temporary`, `Found` and `SeeOther`.",
					},
					"target_listener_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The name of the listener to redirect to. Cannot be set if `target_url` is set.",
					},
					"target_url": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Url to redirect the request to. Cannot be set if `target_listener_name` is set.",
					},
					"include_path": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
						MarkdownDescription: "Whether or not to include the path in the redirected Url. Defaults to `false`.",
					},
					"include_query_string": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
						MarkdownDescription: "Whether or not to include the query string in the redirected Url. Default to `false`.",
					},
				}),
			},
			"request_routing_rules": {
				Required: true,
				MarkdownDescription: "At least one block has to be defined. The request routing rules block has to be defiend as a map with a key name for each `request_routing_rule`. See Example usage for details.",
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of this Request Routing Rule.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `request_routing_rule`.",
					},
					"rule_type": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Type of Routing that should be used for this Rule. Possible values are `Basic` and `PathBasedRouting`.",
					},
					"priority": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "Rule evaluation order can be dictated by specifying an integer value from `1` to `20000` with `1` being the highest priority and `20000` being the lowest priority."+
						"For this version, the priority is computed by the provider (between 1 and 300) after getting the list of used values from the gateway.",
					},
					"http_listener_name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of the HTTP Listener which should be used for this Routing Rule. It has to match a Http Listener name declared in the binding service resource.",
					},
					"backend_address_pool_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Name of the Backend Address Pool which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set."+
						"It has to match a Backend Address Pool name declared in the binding service resource.",
					},
					"backend_http_settings_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Name of the Backend HTTP Settings Collection which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set."+
						"It has to match a Backend HTTP Settings name declared in the binding service resource.",
					},
					"redirect_configuration_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Name of the Redirect Configuration which should be used for this Routing Rule. "+
						"Cannot be set if either `backend_address_pool_name` or `backend_http_settings_name` is set."+
						"It has to match a Redirect Configuration name declared in the binding service resource.",
					},
					"rewrite_rule_set_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Name of the Rewrite Rule Set which should be used for this Routing Rule. Only valid for v2 SKUs. Not supported in this version",
					},
					"url_path_map_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Name of the URL Path Map which should be associated with this Routing Rule. Not supported in this version",
					},
				},tfsdk.MapNestedAttributesOptions{}),
			},
			"http_listeners": {
				Required: true,
				MarkdownDescription: "At least one block has to be defined. The http_listeners block has to be defiend as a mapwith a key name for each `http_listener`. See Example usage for details.",
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of the HTTP Listener.",
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
						MarkdownDescription: "The ID of the `http_listener`.",
					},
					"frontend_ip_configuration_name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of the Frontend IP Configuration used for this HTTP Listener.",
					},
					"frontend_port_name": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Name of the Frontend Port use for this HTTP Listener.",
					},
					"require_sni": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
						MarkdownDescription: "Should Server Name Indication be Required? Defaults to `false`.",
					},
					"protocol": {
						Type:     types.StringType,
						Required: true,
						MarkdownDescription: "The Protocol to use for this HTTP Listener. Possible values are `Http` and `Https`.",
					},
					"host_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The Hostname which should be used for this HTTP Listener. Setting this value changes Listener Type to 'Multi site', however, this option is not supported by the provider version.",
					},
					"host_names": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
						MarkdownDescription: "A list of Hostname(s) should be used for this HTTP Listener. It allows special wildcard characters."+
						"The `host_names` and `host_name` are mutually exclusive and cannot both be set.",
					},
					"ssl_certificate_name": {
						Type:     types.StringType,
						Optional: true,
						MarkdownDescription: "The name of the associated SSL Certificate which should be used for this HTTP Listener."+
						"It has to match a Ssl certificate name declared in the binding service resource.",
					},
				},tfsdk.MapNestedAttributesOptions{}),
			},
		},
	}, nil
}

// New resource instance
func (r resourceBindingServiceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceBindingService{
		p: *(p.(*provider)),
	}, nil
}

// Create a new resource
func (r resourceBindingService) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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
	var plan BindingService
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
	exist_element, exist := checkElementName(gw, plan)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. This (these) element(s) already exist(s) in the app gateway: \n"+ fmt.Sprint(exist_element),
			"Please, change its (their) name(s) then retry.",
		)
		return
	}/*
	exist_element, exist = checkPlanElementName(plan)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. This (these) element(s) have the same key in the configuration: \n"+ fmt.Sprint(exist_element),
			"Please, change its (their) key(s) then retry.",
		)
		return
	}*/
	
	//create, map and add the new elements (json) object from the plan to the agw object
	/************* generate and add BackendAddressPool **************/
	gw.Properties.BackendAddressPools = append(
		gw.Properties.BackendAddressPools, createBackendAddressPool(
			plan.Backend_address_pool))
	
	/************* generate and add request Routing Rule Map **************/
	for key, requestRoutingRule_plan := range plan.Request_routing_rules {
		if checkRequestRoutingRuleCreate(key, plan, gw, resp){
			return
		}
		priority := generatePriority(gw,"high")
		requestRoutingRule_json := createRequestRoutingRule(&requestRoutingRule_plan,priority,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules,requestRoutingRule_json)
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

	/************* generate and add Http listener Map **************/
	for _, httpListener_plan := range plan.Http_listeners { 
		if checkHTTPListenerCreate(httpListener_plan, plan, gw, resp) {
			return
		}
		httpListener_json := createHTTPListener(&httpListener_plan,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)	
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpListener_json)
	}	
	
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
	
	//printToFile(error_json,"updateGW_create.json")
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
	sslCertificate_state 			:= generateSslCertificateState(gw_response,plan.Ssl_certificate.Name.Value)
	redirectConfiguration_state 	:= generateRedirectConfigurationState(gw_response,plan.Redirect_configuration.Name.Value)
	
	httpListeners_state := make(map [string]Http_listener, len(plan.Http_listeners))
	for key, value := range plan.Http_listeners { 
		httpListeners_state[key] = generateHTTPListenerState(gw_response,value.Name.Value)
	}
	requestRoutingRules_state := make(map [string]Request_routing_rule, len(plan.Request_routing_rules))
	for key, value := range plan.Request_routing_rules { 
		requestRoutingRules_state[key] = generateRequestRoutingRuleState(gw_response,value.Name.Value)
	}
	
	var result BindingService
	result = BindingService{
		Name						: plan.Name,
		Agw_name					: types.String{Value: gw_response.Name},
		Agw_rg						: plan.Agw_rg,
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
		Http_listeners				: httpListeners_state,
		Request_routing_rules		: requestRoutingRules_state,
	}
	
	//store to the created object to the terraform state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceBindingService) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	fmt.Println("\n######################## Read Method ########################")
	
	// Get current state
	var state BindingService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	names_map := map[string]string{
		"bindingServiceName"			: state.Name.Value,
		"applicationGatewayName"		: state.Agw_name.Value,
		"resourceGroupName"				: state.Agw_rg.Value,
		"backendAddressPoolName"		: state.Backend_address_pool.Name.Value,
		"backendHTTPSettingsName"		: state.Backend_http_settings.Name.Value,
		"probeName"						: state.Probe.Name.Value,
		"sslCertificateName"			: state.Ssl_certificate.Name.Value,
		"redirectConfigurationName"		: state.Redirect_configuration.Name.Value,		
	}
	
	state = getBindingServiceState(r.p.AZURE_SUBSCRIPTION_ID, names_map, state.Http_listeners, state.Request_routing_rules, r.p.token.Access_token)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceBindingService) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	fmt.Println("\n######################## Update Method ########################")
	// Get plan values
	var plan BindingService
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state BindingService
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//check if there is an update of the gw name or resource group, otherwise, issue exit error. 
	//Propose to remove (destroy) the binding from the initial gw and create a new one in the new gw
	if (plan.Agw_rg.Value != state.Agw_rg.Value) || (plan.Agw_name.Value != state.Agw_name.Value) {
		//this is an error. issue an exit error.
		resp.Diagnostics.AddError(
			"You are willing to update the binding app gateway name or resource group name. This is not supported currently. "+ 
			"We recommand you to make a destroy of the existing binding before creating a new one with the suitable app gateway name and resource group.",
			"Please, change configuration then retry.",
		)
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
	
	// *********** Processing http Listener Map *********** //	
	//preparing the new elements (json) from the plan
	for key, httpListener_plan := range plan.Http_listeners {
		if checkHTTPListenerUpdate(httpListener_plan, plan, gw, resp) {
			return
		}
		// we have to remove the old http listener before creating the new one
		httpListener_state, exist := state.Http_listeners[key]
		// if the http_listener that exist in the plan exist also in the state
		if exist && (httpListener_plan.Name.Value == httpListener_state.Name.Value) {
			//so remove the old one before adding the new one.
			removeHTTPListenerElement(&gw, httpListener_plan.Name.Value)
		}else{
			// it's most likely about http Listener update:
			//	1) with a new name, 
			//	2) or with a new key 
			//	3) or it no longer exist
			
			//remove the old http Listener (old http listener name under the same key) from the gateway
			if exist {
				removeHTTPListenerElement(&gw, httpListener_state.Name.Value)
			}
			//check if the httpListener_plan name already exist in the old state but under different key, in order to remove it
			if checkHTTPListenerNameInMap(httpListener_plan.Name.Value, state.Http_listeners) {
				removeHTTPListenerElement(&gw, httpListener_plan.Name.Value)
			}
			// now check if the new http Listener name is already used in the gateway, no need to check it in the http listener map, 
			// because it will be done incrementally whenever a new http listener is added to the gw.
			if checkHTTPListenerElement(gw, httpListener_plan.Name.Value) {
				//this is an error. issue an exit error.
				resp.Diagnostics.AddError(
					"Unable to update the app gateway. The new http Listener name : "+ httpListener_plan.Name.Value+" already exists. "+
					"It can be due to the name of the http listener you are under declaring",
					" Please, change the name then retry.",				)
				return
			}
		}
		httpListener_json := createHTTPListener(&httpListener_plan,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)	
		//add the new one to the gw
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpListener_json)
	}
	//check if there are some http_listeners that exist in the state but no longer exist in the plan
	//they have to be removed from the gateway
	for _, httpListener_state := range state.Http_listeners {
		if !checkHTTPListenerNameInMap(httpListener_state.Name.Value, plan.Http_listeners) {
			removeHTTPListenerElement(&gw, httpListener_state.Name.Value)
		}
	}

	var priority int 	
	// *********** Processing request Routing Rule Map *********** //	
	//preparing the new elements (json) from the plan
	for key, requestRoutingRule_plan := range plan.Request_routing_rules { 
		if checkRequestRoutingRuleUpdate(key, plan, gw, resp) {
			return
		}
		//to compute priority, check if Request Routing Rule exist in the state, so we get the old priority
		// else, that means the old Request Routing Rule was removed manually, we have to generate a new priority
		requestRoutingRule_state, exist := state.Request_routing_rules[key]
		if exist {
			if requestRoutingRule_state.Priority.Value != "0" && requestRoutingRule_state.Priority.Value != "" {
				//the priority of new Request_routing_rule_http is already included in gw, so it's ok
				priority,_ = strconv.Atoi(requestRoutingRule_state.Priority.Value)
			}else{
				priority = generatePriority(gw,"high")
			}
		}else{
			priority = generatePriority(gw,"high")
		}
		
		//new request Routing Rule is ok. now we have to remove the old one
		//requestRoutingRule_state, exist := state.Request_routing_rules[key]
		// if the request Routing Rule that exist in the plan exist also in the state		  
		if exist && (requestRoutingRule_plan.Name.Value == requestRoutingRule_state.Name.Value) {
			//so we remove the old one before adding the new one.
			removeRequestRoutingRuleElement(&gw, requestRoutingRule_plan.Name.Value)
		}else{
			// it's most likely about request Routing Rule update:
			//	1) with a new name, 
			//	2) or with a new key 
			//	3) or it no longer exist

			//remove the old request Routing Rule (old name) from the gateway
			if exist {
				removeRequestRoutingRuleElement(&gw, requestRoutingRule_state.Name.Value)
			}
			//check if the requestRoutingRule_plan name already exist in the old state but under different key, in order to remove it
			if checkRequestRoutingRuleNameInMap(requestRoutingRule_plan.Name.Value, state.Request_routing_rules) {
				removeRequestRoutingRuleElement(&gw, requestRoutingRule_plan.Name.Value)
			}
			// we have to check if the new request Routing Rule name is already used in the gateway, no need to check it in the requestRoutingRule map, 
			// because it will be done incrementally whenever a new requestRoutingRule is added to the gw.
			if checkRequestRoutingRuleElement(gw, requestRoutingRule_plan.Name.Value) {
				//this is an error. issue an exit error.
				resp.Diagnostics.AddError(
					"Unable to update the app gateway. The new request Routing Rule name : "+ requestRoutingRule_plan.Name.Value+" already exists. "+
					"It can be due to the name of the request Routing Rule you are under declaring",
					" Please, change the name then retry.",
				)
				return
			}			
		}
		requestRoutingRule_json := createRequestRoutingRule(&requestRoutingRule_plan, priority, 
			r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName)		
		//add the new one to the gw
		gw.Properties.RequestRoutingRules = append(gw.Properties.RequestRoutingRules,requestRoutingRule_json)
	}
	//check if there are some request Routing Rules that exist in the state but no longer exist in the plan
	//they have to be removed from the gateway
	for _, requestRoutingRule_state := range state.Request_routing_rules {
		if !checkRequestRoutingRuleNameInMap(requestRoutingRule_state.Name.Value, plan.Request_routing_rules) {
			removeRequestRoutingRuleElement(&gw, requestRoutingRule_state.Name.Value)
		}
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
	gw.Properties.SslCertificates = append(gw.Properties.SslCertificates, sslCertificate_json)
	gw.Properties.RedirectConfigurations = append(gw.Properties.RedirectConfigurations, redirectConfiguration_json)
	
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
	sslCertificate_state 			:= generateSslCertificateState(gw_response,sslCertificate_json.Name)
	redirectConfiguration_state 	:= generateRedirectConfigurationState(gw_response,redirectConfiguration_json.Name)
	
	httpListeners_state := make(map [string]Http_listener, len(plan.Http_listeners))
	for key, value := range plan.Http_listeners { 
		httpListeners_state[key] = generateHTTPListenerState(gw_response,value.Name.Value)
	}
	requestRoutingRules_state := make(map [string]Request_routing_rule, len(plan.Request_routing_rules))
	for key, value := range plan.Request_routing_rules { 
		requestRoutingRules_state[key] = generateRequestRoutingRuleState(gw_response,value.Name.Value)
	}

	/*************** Special for Http listener **********************/
	// Generate resource state struct 
	//i moved "Generate resource state struct" with http listner block before it depends on the later.
	
	var result BindingService
	result = BindingService{
		Name						: plan.Name,
		Agw_name					: types.String{Value: gw_response.Name},
		Agw_rg						: plan.Agw_rg,
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
		Http_listeners				: httpListeners_state,
		Request_routing_rules		: requestRoutingRules_state,
	}
	
	//store to the created objecy to the terraform state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceBindingService) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	fmt.Println("\n######################## Delete Method ########################")
	// Get current state
	var state BindingService
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get elements names from state
	backendAddressPoolName 		:= state.Backend_address_pool.Name.Value
	backendHTTPSettingsName 	:= state.Backend_http_settings.Name.Value
	probeName 					:= state.Probe.Name.Value
	sslCertificateName 			:= state.Ssl_certificate.Name.Value
	redirectConfigurationName 	:= state.Redirect_configuration.Name.Value
	
	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//remove the elements from the gw
	removeBackendAddressPoolElement(&gw, backendAddressPoolName)
	removeBackendHTTPSettingsElement(&gw,backendHTTPSettingsName)
	removeProbeElement(&gw,probeName)
	removeSslCertificateElement(&gw,sslCertificateName)
	removeRedirectConfigurationElement(&gw,redirectConfigurationName)
	
	for _, httpListener_state := range state.Http_listeners { 
		removeHTTPListenerElement(&gw,httpListener_state.Name.Value)		
	}
	for _, requestRoutingRule_state := range state.Request_routing_rules { 
		removeRequestRoutingRuleElement(&gw,requestRoutingRule_state.Name.Value)		
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
func (r resourceBindingService) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	//the ID given in the import command should match exactly the following format:
	// <gw_name,gw-resourcegroup,backend_address_pool_name,backend_http_settings_name,probe_name,ssl_certificate,
	//https_listener_name,request_routing_rule_https_name,redirect_configuration_name,
	//http_listener_name(optional),request_routing_rule_http_name(optional, required if http_listener_name is set)>
	/*
	idParts := strings.Split(req.ID, ",")
	//check if the given ID contains the right number of params (9 or 11)
	if (len(idParts) != 11 && len(idParts) != 9) {
        resp.Diagnostics.AddError(
            "Unexpected Import Identifier. The identifier should be composed of 9 or 11 params matching exactly the following format: \n"+
			"<gw_name,gw-resourcegroup,backend_address_pool_name,backend_http_settings_name,probe_name,ssl_certificate,"+
			"https_listener_name,request_routing_rule_https_name,redirect_configuration_name,"+
			"http_listener_name(optional),request_routing_rule_http_name(optional, required if http_listener_name is set)>",
            "Please, check the import identifier then retry",
        )
        return
    }

	//check if there is an empty param
	for i := 0; i < len(idParts); i++ {
		if idParts[i] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier. A given param is empty",
				"Please, check the import identifier then retry",
			)
			return
		}
	}
	uniqueId := RandStringBytes(10)
	
	names_map := map[string]string{
		"bindingServiceName"			: "binding_"+uniqueId, //generate random unique name for the imported resource
		"applicationGatewayName"		: idParts[0],
		"resourceGroupName"				: idParts[1],
		"backendAddressPoolName"		: idParts[2],
		"backendHTTPSettingsName"		: idParts[3],
		"probeName"						: idParts[4],
		"sslCertificateName"			: idParts[5],
		"httpsListenerName"				: idParts[6],
		"requestRoutingRuleHttpsName"	: idParts[7],
		"redirectConfigurationName"		: idParts[8],		
	}
	if len(idParts) == 11{
		names_map["httpListenerName"] = idParts[9]
		names_map["requestRoutingRuleHttpName"] = idParts[10]
	}

	state := getBindingServiceState(r.p.AZURE_SUBSCRIPTION_ID,names_map,r.p.token.Access_token)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}*/
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

// specific processing for binding service
func getBindingServiceState(AZURE_SUBSCRIPTION_ID string, names_map map[string]string, http_listeners map[string]Http_listener, 
	request_routing_rules map[string]Request_routing_rule, Access_token string) BindingService {
	
	// Get gw from API and then update what is in state from what the API returns
	bindingServiceName := names_map["bindingServiceName"] 

	//Get the agw
	resourceGroupName := names_map["resourceGroupName"] 
	applicationGatewayName := names_map["applicationGatewayName"] 
	gw := getGW(AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, Access_token)
			
	// *********** Processing the backend address pool *********** //
	var backendAddressPool_state Backend_address_pool
	backendAddressPoolName := names_map["backendAddressPoolName"] 
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
	backendHTTPSettingsName := names_map["backendHTTPSettingsName"] 
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
	probeName := names_map["probeName"] 
	//check if the probe exists in the gateway, otherwise, it was removed manually
	if checkProbeElement(gw, probeName) {
		//generate probe state
		probe_state = generateProbeState(gw, probeName)
	}else{
		//generate an empty probe_state because it was removed manually
		probe_state = Probe_tf{}
	}
	
	// *********** Processing the SSL Certificate *********** //
	//check if the SSL Certificate  exists in  the gateway, otherwise, it was removed manually
	var sslCertificate_state Ssl_certificate
	sslCertificateName := names_map["sslCertificateName"] 
	if checkSslCertificateElement(gw, sslCertificateName) {
		sslCertificate_state = generateSslCertificateState(gw,sslCertificateName)
	}else{
		sslCertificate_state = Ssl_certificate{}
	}

	// *********** Processing the Redirect Configuration *********** //
	var redirectConfiguration_state Redirect_configuration
	redirectConfigurationName := names_map["redirectConfigurationName"] 
	//check if the Redirect Configuration exists in the gateway, otherwise, it was removed manually
	if checkRedirectConfigurationElement(gw, redirectConfigurationName) {
		//generate BackendState
		redirectConfiguration_state = generateRedirectConfigurationState(gw, redirectConfigurationName)
	}else{
		//generate an empty redirectConfiguration_state because it was removed manually
		redirectConfiguration_state = Redirect_configuration{}
	}	
	
	var result BindingService
	result = BindingService{
		Name						: types.String{Value: bindingServiceName},
		Agw_name					: types.String{Value: names_map["applicationGatewayName"]},
		Agw_rg						: types.String{Value: names_map["resourceGroupName"]},
		Backend_address_pool		: backendAddressPool_state,
		Backend_http_settings		: backendHTTPSettings_state,
		Probe						: probe_state,
		Ssl_certificate				: sslCertificate_state,
		Redirect_configuration		: redirectConfiguration_state,
	}
		
	// *********** Processing the http Listener Map *********** //
	//check if the Https listener  exists in  the gateway, otherwise, it was removed manually
	
	httpListeners_state := make(map [string]Http_listener, len(http_listeners))
	
	for key, value := range http_listeners { 
		var httpListener_state Http_listener
		if checkHTTPListenerElement(gw, value.Name.Value) {
			httpListener_state = generateHTTPListenerState(gw,value.Name.Value)
		}else{
			httpListener_state = Http_listener{}
		}
		httpListeners_state[key] = httpListener_state
	}
	result.Http_listeners = httpListeners_state
	
	// *********** Processing the request Routing Rule Map *********** //
	//check if the request Routing Rule exists in  the gateway, otherwise, it was removed manually
	requestRoutingRules_state := make(map [string]Request_routing_rule, len(request_routing_rules))
	
	for key, value := range request_routing_rules { 
		var requestRoutingRule_state Request_routing_rule
		if checkRequestRoutingRuleElement(gw, value.Name.Value) {
			requestRoutingRule_state = generateRequestRoutingRuleState(gw,value.Name.Value)
		}else{
			requestRoutingRule_state = Request_routing_rule{}
		}
		requestRoutingRules_state[key] = requestRoutingRule_state
	}
	result.Request_routing_rules = requestRoutingRules_state

	return result
}
func checkElementName(gw ApplicationGateway, plan BindingService) ([]string,bool){
	//This function allows to check if an element name in the required new configuration (plan BindingService) already exist in the gw.
	//if so, the provider has to stop executing and issue an exit error
	exist := false
	var existing_element_list [] string
	//Create new var for all configurations
	backendAddressPool_plan 	:= plan.Backend_address_pool 
	backendHTTPSettings_plan 	:= plan.Backend_http_settings
	probe_plan 					:= plan.Probe
	sslCertificate_plan			:= plan.Ssl_certificate
	redirectConfiguration_plan	:= plan.Redirect_configuration
	
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
	if checkSslCertificateElement(gw, sslCertificate_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- SSL Certificate: "+sslCertificate_plan.Name.Value)
	}
	if checkRedirectConfigurationElement(gw, redirectConfiguration_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- Redirect configuration: "+redirectConfiguration_plan.Name.Value)
	}
	for key, httpListener_plan := range plan.Http_listeners { 
		if checkHTTPListenerElement(gw, httpListener_plan.Name.Value) {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- HTTPListener ("+key+"): "+httpListener_plan.Name.Value)
		}
	}
	//check if the http_listener map contains a repetitive http_listener names
	for key, httpListener_plan := range plan.Http_listeners { 
		for key1, httpListener_plan1 := range plan.Http_listeners {
			if (httpListener_plan.Name.Value == httpListener_plan1.Name.Value) && (key != key1) {
				exist = true 
				existing_element_list = append(existing_element_list,"\n	- HTTPListener ("+key+" and "+key1+"): "+httpListener_plan.Name.Value)
			}
		}
	}
	for key, requestRoutingRule_plan := range plan.Request_routing_rules { 
		if checkRequestRoutingRuleElement(gw, requestRoutingRule_plan.Name.Value) {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- Request Routing Rule ("+key+"): "+requestRoutingRule_plan.Name.Value)
		}
	}
	//check if the requestRoutingRule map contains a repetitive requestRoutingRule names
	for key, requestRoutingRule_plan := range plan.Request_routing_rules { 
		for key1, requestRoutingRule_plan1 := range plan.Request_routing_rules {
			if (requestRoutingRule_plan.Name.Value == requestRoutingRule_plan1.Name.Value) && (key != key1) {
				exist = true 
				existing_element_list = append(existing_element_list,"\n	- Request Routing Rule ("+key+" and "+key1+"): "+requestRoutingRule_plan.Name.Value)
			}
		}
	}
	existing_element_list = append(existing_element_list,"\n")
	return existing_element_list,exist
}
func checkPlanElementName(plan BindingService) ([]string,bool){
	exist := false
	var existing_element_list [] string
	for key, httpListener_plan := range plan.Http_listeners { 
		for key1, httpListener_plan1 := range plan.Http_listeners {
			if (httpListener_plan.Name.Value != httpListener_plan1.Name.Value) && (key != key1) {
				exist = true 
				existing_element_list = append(existing_element_list,"\n	- HTTPListener ("+key+" and "+key1+"): "+httpListener_plan.Name.Value)
			}
		}
	}
	for key, requestRoutingRule_plan := range plan.Request_routing_rules { 
		for key1, requestRoutingRule_plan1 := range plan.Request_routing_rules {
			if (requestRoutingRule_plan.Name.Value == requestRoutingRule_plan1.Name.Value) && (key != key1) {
				exist = true 
				existing_element_list = append(existing_element_list,"\n	- Request Routing Rule ("+key+" and "+key1+"): "+requestRoutingRule_plan.Name.Value)
			}
		}
	}
	existing_element_list = append(existing_element_list,"\n")
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
func RandStringBytes(n int) string {
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}