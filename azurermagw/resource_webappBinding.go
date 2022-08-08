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
						//to implment this, it requires additional effort. Actually, it is easier for me
						//to make it Required :)
						Required: true,
						//Default: false,
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
		},
	}, nil
}
func stringDefault(defaultValue string) stringDefaultModifier {
    return stringDefaultModifier{
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
	exist_element, exist := checkElementName(gw, plan)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. At least, those elements : "+ fmt.Sprint(exist_element),
			"Already exists in the app gateway. Please, modify their names.",
		)
		return
	}

	//create, map and add the new elements (json) object from the plan (plan) to the agw object
	gw.Properties.BackendAddressPools = append(
		gw.Properties.BackendAddressPools, createBackendAddressPool(
			plan.Backend_address_pool))
	gw.Properties.BackendHTTPSettingsCollection = append(
		gw.Properties.BackendHTTPSettingsCollection, createBackendHTTPSettings(
			plan.Backend_http_settings,
			r.p.AZURE_SUBSCRIPTION_ID,
			resourceGroupName,
			applicationGatewayName))

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
	
	//generate the States.
	nb_Fqdns 		:= len(plan.Backend_address_pool.Fqdns)
	nb_IpAddress	:= len(plan.Backend_address_pool.Ip_addresses)
	backendAddressPool_state 	:= generateBackendAddressPoolState(
		gw_response, 
		plan.Backend_address_pool.Name.Value,
		nb_Fqdns,
		nb_IpAddress)
	backendHTTPSettings_state 	:= generateBackendHTTPSettingsState(
		gw_response,
		plan.Backend_http_settings.Name.Value)

	// Generate resource state struct
	var result = WebappBinding{
		Name					: plan.Name,
		Agw_name				: types.String{Value: gw_response.Name},
		Agw_rg					: plan.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
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

	// *********** Processing backend address pool *********** //
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
	
	// *********** Processing backend http settings *********** //
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
	//fmt.Printf("\nHHHHHHHHHHHHHHHHH backendHTTPSettings_state =\n %+v ",backendHTTPSettings_state)
	

	// Generate resource state struct
	var result = WebappBinding{
		Name					: types.String{Value: webappBindingName},
		Agw_name				: state.Agw_name,
		Agw_rg					: state.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
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

	//preparing the new elements (json) from the plan
	//create and map the new Backend pool element (backendAddressPool_json) object from the plan (backendAddressPool_plan)
	backendAddressPool_plan := plan.Backend_address_pool
	backendAddressPool_json := createBackendAddressPool(backendAddressPool_plan)
	backendHTTPSettings_plan := plan.Backend_http_settings
	backendHTTPSettings_json := createBackendHTTPSettings(
		backendHTTPSettings_plan,r.p.AZURE_SUBSCRIPTION_ID,
		resourceGroupName,applicationGatewayName)

	//Verify if the agw already contains the elements to be updated.
	//the older ones should be removed before updating. 
	//we have also to prevent element name updating and manual deletion
	
	// *********** Processing backend address pool *********** //	
	//check if the backend name in the plan and state are different, that means that the
	if backendAddressPool_plan.Name.Value == state.Backend_address_pool.Name.Value {
		//it is about backend AddressPool update  with the same name
		//so we remove the old one before adding the new one.
		removeBackendAddressPoolElement(&gw, backendAddressPool_json.Name)
	}else{
		// it's about backend update with a new name
		// we have to check if the new backend name is already used
		if checkBackendAddressPoolElement(gw, backendAddressPool_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The Backend Adresse pool name : "+ backendAddressPool_json.Name+" already exists.",
				" Please, modify the name.",
			)
			return
		}
		//remove the old backend (old name) from the gateway
		removeBackendAddressPoolElement(&gw, state.Backend_address_pool.Name.Value)
	}

	// *********** Processing backend http settings *********** //	
	//check if the backend name in the plan and state are different, that means that the
	if backendHTTPSettings_plan.Name.Value == state.Backend_http_settings.Name.Value {
		//it is about backend http settings update  with the same name
		//so we remove the old one before adding the new one.
		removeBackendHTTPSettingsElement(&gw, backendHTTPSettings_json.Name)
	}else{
		// it's about backend http settings update with a new name
		// we have to check if the new backend http settings name is already used
		if checkBackendHTTPSettingsElement(gw, backendHTTPSettings_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The Backend Adresse pool name : "+ backendHTTPSettings_json.Name+" already exists.",
				" Please, modify the name.",
			)
			return
		}
		//remove the old backend http settings (old name) from the gateway
		removeBackendHTTPSettingsElement(&gw, state.Backend_http_settings.Name.Value)
	}

	//add the new elements
	gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, backendAddressPool_json)
	gw.Properties.BackendHTTPSettingsCollection = append(gw.Properties.BackendHTTPSettingsCollection, backendHTTPSettings_json)
	
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
	
	// in the Read method, the number of fqdns and Ip in a Backendpool should be calculated from the json object and not the plan or state,
	// because the purpose of the read is to see if there is a difference between the real element and the satate stored localy.
	/*********** Special for Backend Address Pool ********************/
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

	// Generate resource state struct
	var result = WebappBinding{
		Name					: state.Name,
		Agw_name				: types.String{Value: gw_response.Name},
		Agw_rg					: state.Agw_rg,
		Backend_address_pool	: backendAddressPool_state,
		Backend_http_settings	: backendHTTPSettings_state,
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
	backendAddressPool_name := state.Backend_address_pool.Name.Value
	backendHTTPSettings_name := state.Backend_http_settings.Name.Value
	
	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	
	//remove the backend from the gw
	removeBackendAddressPoolElement(&gw, backendAddressPool_name)
	removeBackendHTTPSettingsElement(&gw,backendHTTPSettings_name)

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
	exist := false
	var existing_element_list [] string
	//Create new var for all configurations
	backendAddressPool_plan := plan.Backend_address_pool 
	backendHTTPSettings_plan := plan.Backend_http_settings
	
	if checkBackendAddressPoolElement(gw, backendAddressPool_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- "+backendAddressPool_plan.Name.Value)
	}
	if checkBackendHTTPSettingsElement(gw, backendHTTPSettings_plan.Name.Value) {
		exist = true 
		existing_element_list = append(existing_element_list,"\n	- "+backendHTTPSettings_plan.Name.Value)
	}
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
