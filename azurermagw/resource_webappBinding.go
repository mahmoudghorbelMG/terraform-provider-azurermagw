package azurermagw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
						//this params should be optional but whith default value (false)
						//to implment this, it requires additional effort. Actually, it's easier for me
						//to make it Required :)
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
		},
	}, nil
}
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
	fmt.Println("\n######################## Create Method1 ########################")
	
	// Retrieve values from plan
	var plan WebappBinding
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	fmt.Println("\n######################## Create Method 2 ########################")
	
	//Get the agw (app gateway) from Azure with its Rest API
	resourceGroupName := plan.Agw_rg.Value
	applicationGatewayName := plan.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//Check if the agw already contains an existing element that has the same name of a new element to add
	exist_element, exist := checkElementName(gw, plan)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. At least, these elements already exists in the app gateway: "+ fmt.Sprint(exist_element),
			"Please, change their names.",
		)
		return
	}
	fmt.Println("\n######################## Create Method 3 ########################")
	
	//create, map and add the new elements (json) object from the plan (plan) to the agw object
	gw.Properties.BackendAddressPools = append(
		gw.Properties.BackendAddressPools, createBackendAddressPool(
			plan.Backend_address_pool))

	backendHTTPSettings_json, error_probeName := createBackendHTTPSettings(plan.Backend_http_settings,plan.Probe.Name.Value,
												r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	if error_probeName== "fatal" {
		resp.Diagnostics.AddError(
			"Unable to create binding. The probe name ("+plan.Backend_http_settings.Probe_name.Value+") declared in Backend_http_settings: "+ 
			plan.Backend_http_settings.Name.Value+" doesn't match the probe name conf : "+plan.Probe.Name.Value,
			"Please, change probe name then retry.",
		)
		return
	}
	gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection,backendHTTPSettings_json)
	gw.Properties.Probes = append(gw.Properties.Probes,
		createProbe(plan.Probe,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName))

	/************* Processing Http listener **************/
	// no ssl certificate to provider, so no need to check error_SslCertificateName
	/*var h *Http_listener	
	h = &plan.Http_listener*/
	fmt.Println("\n######################## Create Method 4 ########################")
	
	//if hasField(plan,"Http_listener"){
	
	if &plan.Http_listener != nil {
		fmt.Println("\n######################## Create Method 5 ########################")
		SslCertificateName:=""
		httpListener_json, _,error_Hostname := createHTTPListener(plan.Http_listener,SslCertificateName,
				r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		if error_Hostname == "fatal-exclusivity" {
			//hostname and hostnames are mutually exclusive. only one should be provided
			resp.Diagnostics.AddError(
				"Unable to create binding. In HTTP Listener "+ plan.Http_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
				"Only one should be provided",
				"Please, change HTTPListener configuration then retry.",)
			return
		}
		if error_Hostname == "fatal-missing" {
			//hostname and hostnames are mutually exclusive. only one should be provided
			resp.Diagnostics.AddError(
				"Unable to create binding. In HTTP Listener "+ plan.Http_listener.Name.Value+", both Hostname and Hostnames are missing. "+
				"At least and only one should be provided",
				"Please, change HTTPListener configuration then retry.",)
			return
		}
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpListener_json)
	}
	/************* Processing Https listener **************/
	//SslCertificateName := plan.SslCertificate.Name.Value // (not yet implemented till now)
	SslCertificateName :="default-citeo-adelphe-cert"
	httpsListener_json, error_SslCertificateName,error_Hostname := createHTTPListener(plan.Https_listener,SslCertificateName,
		r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	if error_SslCertificateName == "fatal" {
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to create binding. The SslCertificate name ("+SslCertificateName+") declared in Http_listener: "+ 
		plan.Https_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Https_listener.Ssl_certificate_name.Value,
		"Please, change probe name then retry.",)
		return
	}
	if error_Hostname == "fatal-exclusivity" {
		//hostname and hostnames are mutually exclusive. only one should be provided
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTP Listener "+ plan.Https_listener.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be provided",
			"Please, change HTTPListener configuration then retry.",)
		return
	}
	if error_Hostname == "fatal-missing" {
		//hostname and hostnames are mutually exclusive. only one should be provided
		resp.Diagnostics.AddError(
			"Unable to create binding. In HTTP Listener "+ plan.Https_listener.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be provided",
			"Please, change HTTPListener configuration then retry.",)
		return
	}		
	gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners,httpsListener_json)	
	

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
	backendAddressPool_state 	:= generateBackendAddressPoolState(
		gw_response,plan.Backend_address_pool.Name.Value,nb_Fqdns,nb_IpAddress)
	backendHTTPSettings_state 	:= generateBackendHTTPSettingsState(
		gw_response,plan.Backend_http_settings.Name.Value)
	probe_state := generateProbeState(gw_response,plan.Probe.Name.Value)
	
	var httpListener_state Http_listener
	//if hasField(plan,"Http_listener"){
	if &plan.Http_listener != nil {
		httpListener_state 	= generateHTTPListenerState(gw_response,plan.Http_listener.Name.Value)
	}else{
		httpListener_state = Http_listener{}
	}
	tflog.Info(ctx,"[CREATE] httpListener_state :",  map[string]interface{}{"httpListener_state ": httpListener_state,})
	
	httpsListener_state 	:= generateHTTPListenerState(gw_response,plan.Https_listener.Name.Value)
	

	// Generate resource state struct
	var result = WebappBinding{
		Name					: plan.Name,
		Agw_name				: types.String{Value: gw_response.Name},
		Agw_rg					: plan.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
		Probe					: probe_state,
		Http_listener			: &httpListener_state,
		Https_listener			: &httpsListener_state,
	}
	//store to the created objecy to the terraform state
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

	// *********** Processing the http Listener *********** //
	//check if the Http listener  exists in the old state (because it's optional param) 
	//in order to check if it's in the gateway, otherwise, it was removed manually
	var httpListener_state Http_listener
	//if hasField(state,"Http_listener"){
	if &state.Http_listener != nil{
		httpListenerName := state.Http_listener.Name.Value
		if checkHTTPListenerElement(gw, httpListenerName) {
			httpListener_state 	= generateHTTPListenerState(gw,httpListenerName)
		}else{
			httpListener_state = Http_listener{}
		}
	}
	
	// *********** Processing the https Listener *********** //
	//check if the Https listener  exists in in the gateway, otherwise, it was removed manually
	var httpsListener_state Http_listener
	httpsListenerName := state.Https_listener.Name.Value
	if checkHTTPListenerElement(gw, httpsListenerName) {
		httpsListener_state = generateHTTPListenerState(gw,httpsListenerName)
	}else{
		httpsListener_state = Http_listener{}
	}
		
	// Generate resource state struct
	var result = WebappBinding{
		Name					: types.String{Value: webappBindingName},
		Agw_name				: state.Agw_name,
		Agw_rg					: state.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
		Probe					: probe_state,
		Http_listener			: &httpListener_state,
		Https_listener			: &httpsListener_state,
	}

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
	tflog.Info(ctx,"\n######################## Update Method ########################")
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
	tflog.Info(ctx,"\n*********** Processing backend address pool ***********")
	
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
				" Please, change the name.",
			)
			return
		}
		//remove the old backend (old name) from the gateway
		removeBackendAddressPoolElement(&gw, state.Backend_address_pool.Name.Value)
	}

	// *********** Processing backend http settings *********** //	
	tflog.Info(ctx,"\n*********** Processing backend http settings ***********")
	//preparing the new elements (json) from the plan
	backendHTTPSettings_plan := plan.Backend_http_settings
	backendHTTPSettings_json, error_probeName := createBackendHTTPSettings(backendHTTPSettings_plan,plan.Probe.Name.Value,r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	
	//check the provided probe name 
	if error_probeName== "fatal" {
		resp.Diagnostics.AddError(
			"Unable to update binding. The probe name ("+backendHTTPSettings_plan.Probe_name.Value+") declared in Backend_http_settings: "+ 
			backendHTTPSettings_plan.Name.Value+" doesn't match the probe name conf : "+plan.Probe.Name.Value,
			"Please, change probe name then retry.",
		)
		return
	}
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
	tflog.Info(ctx,"\n*********** Processing the probe ***********")
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
				" Please, change the name.",
			)
			return
		}
		//remove the old backend http settings (old name) from the gateway
		removeProbeElement(&gw, state.Probe.Name.Value)
	}

	// *********** Processing http Listener *********** //	
	tflog.Info(ctx,"\n*********** Processing http Listener ***********")
	//preparing the new elements (json) from the plan
	tflog.Info(ctx,"plan.Http_listener before if :",  map[string]interface{}{"plan.Http_listener ": plan.Http_listener,})
	fmt.Printf("\nIIIIIIIIIIIIIIIIIIII  httpListener_plan =\n %+v ",plan.Http_listener)
	//if hasField(plan,"Http_listener"){
	if &plan.Http_listener != nil {
		tflog.Info(ctx,"in if :")
		SslCertificateName:=""
		httpListener_plan := plan.Http_listener
		httpListener_json, _,error_Hostname := createHTTPListener(httpListener_plan,SslCertificateName,
				r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
		if error_Hostname == "fatal-exclusivity" {
			//hostname and hostnames are mutually exclusive. only one should be provided
			resp.Diagnostics.AddError(
				"Unable to update binding. In HTTP Listener "+ httpListener_plan.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
				"Only one should be provided",
				"Please, change HTTPListener configuration then retry.",)
			return
		}
		if error_Hostname == "fatal-missing" {
			//hostname and hostnames are mutually exclusive. only one should be provided
			resp.Diagnostics.AddError(
				"Unable to update binding. In HTTP Listener "+ httpListener_plan.Name.Value+", both Hostname and Hostnames are missing. "+
				"At least and only one should be provided",
				"Please, change HTTPListener configuration then retry.",)
			return
		}
		//new http listener is ok. now we have to remove the old one
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
					" Please, change the name.",
				)
				return
			}
			//remove the old http Listener (old name) from the gateway
			removeHTTPListenerElement(&gw, state.Http_listener.Name.Value)
		}
		//we have to add the http listener here because it's optional
		gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners, httpListener_json)			
	}
			
	// *********** Processing https Listener *********** //	
	//preparing the new elements (json) from the plan
	//SslCertificateName := plan.SslCertificate.Name.Value // (not yet implemented till now)
	SslCertificateName:="default-citeo-adelphe-cert"
	httpsListener_plan := plan.Https_listener
	httpsListener_json, error_SslCertificateName,error_Hostname := createHTTPListener(httpsListener_plan,SslCertificateName,
			r.p.AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName)
	if error_SslCertificateName == "fatal" {
		//wrong SslCertificate Name
		resp.Diagnostics.AddError(
		"Unable to update binding. The SslCertificate name ("+SslCertificateName+") declared in Https_listener: "+ 
		plan.Https_listener.Name.Value+" doesn't match the SslCertificate name conf : "+plan.Https_listener.Ssl_certificate_name.Value,
		"Please, change probe name then retry.",)
		return
	}
	if error_Hostname == "fatal-exclusivity" {
		//hostname and hostnames are mutually exclusive. only one should be provided
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTPS Listener "+ httpsListener_plan.Name.Value+", Hostname and Hostnames are mutually exclusive. "+
			"Only one should be provided",
			"Please, change HTTPS Listener configuration then retry.",)
		return
	}
	if error_Hostname == "fatal-missing" {
		//hostname and hostnames are mutually exclusive. only one should be provided
		resp.Diagnostics.AddError(
			"Unable to update binding. In HTTPS Listener "+ httpsListener_plan.Name.Value+", both Hostname and Hostnames are missing. "+
			"At least and only one should be provided",
			"Please, change HTTPS Listener configuration then retry.",)
		return
	}
	//new http listener is ok. now we have to remove the old one
	if httpsListener_plan.Name.Value == state.Https_listener.Name.Value {
		//so we remove the old one before adding the new one.
		removeHTTPListenerElement(&gw, httpsListener_json.Name)
	}else{
		// it's most likely about http Listener update with a new name
		// we have to check if the new http Listener name is already used
		if checkHTTPListenerElement(gw, httpsListener_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The new http Listener name : "+ httpsListener_json.Name+" already exists.",
				" Please, change the name.",
			)
			return
		}
		//remove the old http Listener (old name) from the gateway
		removeHTTPListenerElement(&gw, state.Https_listener.Name.Value)
	}
	
	//add the new elements (http Listener elements are already added). 
	gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, backendAddressPool_json)
	gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection, backendHTTPSettings_json)
	gw.Properties.Probes = append(gw.Properties.Probes, probe_json)
	gw.Properties.HTTPListeners = append(gw.Properties.HTTPListeners, httpsListener_json)

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
	
	backendAddressPool_state	:= generateBackendAddressPoolState(gw_response, backendAddressPool_json.Name,nb_Fqdns,nb_IpAddress)
	backendHTTPSettings_state	:= generateBackendHTTPSettingsState(gw_response,backendHTTPSettings_json.Name)
	probe_state	:= generateProbeState(gw_response,probe_json.Name)

	/*************** Special for Http listener **********************/
	var httpListener_state Http_listener
	//if hasField(plan,"Http_listener"){
	if &plan.Http_listener != nil {
		httpListener_state 	= generateHTTPListenerState(gw_response,plan.Http_listener.Name.Value)
	}else{
		httpListener_state = Http_listener{}
	}
	/***************************************************************/
	httpsListener_state 	:= generateHTTPListenerState(gw_response,plan.Https_listener.Name.Value)
		
	// Generate resource state struct
	var result = WebappBinding{
		Name					: state.Name,
		Agw_name				: types.String{Value: gw_response.Name},
		Agw_rg					: state.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
		Probe					: probe_state,
		Http_listener			: &httpListener_state,
		Https_listener			: &httpsListener_state,
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
	backendAddressPoolName := state.Backend_address_pool.Name.Value
	backendHTTPSettingsName := state.Backend_http_settings.Name.Value
	probeName := state.Probe.Name.Value
	HTTPSListenerName := state.Https_listener.Name.Value
	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//remove the elements from the gw
	removeBackendAddressPoolElement(&gw, backendAddressPoolName)
	removeBackendHTTPSettingsElement(&gw,backendHTTPSettingsName)
	removeProbeElement(&gw,probeName)
	removeHTTPListenerElement(&gw,HTTPSListenerName)
	/*************** Special for Http listener **********************/
	//if hasField(state,"Http_listener"){
	if &state.Http_listener != nil {
		removeHTTPListenerElement(&gw,state.Http_listener.Name.Value)
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


func checkElementName(gw ApplicationGateway, plan WebappBinding) ([]string,bool){
	//This function allows to check if an element name in the required new configuration (plan WebappBinding) already exist in the gw.
	//if so, the provider has to stop executing and issue an exit error
	fmt.Println("\n######################## Create Method inside check ########################")
	
	exist := false
	var existing_element_list [] string
	//Create new var for all configurations
	backendAddressPool_plan 	:= plan.Backend_address_pool 
	backendHTTPSettings_plan 	:= plan.Backend_http_settings
	probe_plan 					:= plan.Probe
	httpsListener_plan 			:= plan.Https_listener
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
	fmt.Println("\n######################## Create Method inside check before if ########################")
	
	//if hasField(plan,"Http_listener"){
	if &plan.Http_listener != nil {
		httpListener_plan 			:= plan.Http_listener
		if checkHTTPListenerElement(gw, httpListener_plan.Name.Value) {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- HTTPListener: "+httpListener_plan.Name.Value)
		}
		if httpListener_plan.Name.Value == httpsListener_plan.Name.Value {
			exist = true 
			existing_element_list = append(existing_element_list,"\n	- HTTP and HTTPS Listener (new): "+httpListener_plan.Name.Value)
		}
	}
	fmt.Println("\n######################## Create Method inside check after if ########################")
	
	return existing_element_list,exist
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