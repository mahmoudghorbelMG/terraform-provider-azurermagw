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
		},
	}, nil
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
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
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

	//Check if the agw already contains an element that has the same name
	exist_element, exist := checkElementName(gw, plan)
	if exist {
		resp.Diagnostics.AddError(
			"Unable to create binding. At least, those elements : "+ exist_element,
			"Already exists in the app gateway. Please, modify the element name.",
		)
		return
	}

	//create and map the new Backend pool element (backend_json) object from the plan (backend_plan)
	backend_json := createBackendAddressPool(plan.Backend_address_pool)
	gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, backend_json)

	gw_response, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)
	
	//verify if the backend address pool is added to the gateway, otherwise exit error
	if code != 200 {
		// Error  - backend address pool wasn't added to the app gateway
		resp.Diagnostics.AddError(
			"Unable to add Backend Address pool ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Backend Address pool Name doesn't exist in the response of the app gateway",
		)
		return
	}
	
	//generate BackendState.
	nb_Fqdns 		:= len(plan.Backend_address_pool.Fqdns)
	nb_IpAddress	:= len(plan.Backend_address_pool.Ip_addresses)
	backend_state := generateBackendAddressPoolState(gw_response, plan.Backend_address_pool.Name.Value,nb_Fqdns,nb_IpAddress)

	// Generate resource state struct
	var result = WebappBinding{
		Name:                 plan.Name,
		Agw_name:             types.String{Value: gw_response.Name},
		Agw_rg:               plan.Agw_rg,
		Backend_address_pool: backend_state,
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

	var backend_state Backend_address_pool
	backendAddressPoolName := state.Backend_address_pool.Name.Value

	fmt.Printf("\n--------------------- state.Backend_address_pool Content before Read :\n %+v ",state.Backend_address_pool)
	//check if the backend address pool exist in the gateway, otherwise, it was removed manually
	if checkBackendAddressPoolElement(gw, backendAddressPoolName) {
		// in the Read method, the number of fqdns and Ip in a Backendpool should be calculated from the json object and not the plan or state,
		// because the purpose of the read is to see if there is a difference between the real element and the satate stored localy.
		index := getBackendAddressPoolElementKey(gw, backendAddressPoolName)
		backend_json := gw.Properties.BackendAddressPools[index]
		fmt.Printf("\n******************** backend_json Content from gw :\n %+v ",backend_json)
		nb_BackendAddresses := len(backend_json.Properties.BackendAddresses)
		nb_Fqdns := 0
		for i := 0; i < nb_BackendAddresses; i++ {
			if (backend_json.Properties.BackendAddresses[i].Fqdn != "") && (&backend_json.Properties.BackendAddresses[i].Fqdn != nil) {
				nb_Fqdns++
			}
		}
		nb_IpAddress := nb_BackendAddresses - nb_Fqdns

		//generate BackendState
		backend_state = generateBackendAddressPoolState(gw, backendAddressPoolName,nb_Fqdns,nb_IpAddress)
	}else{
		//generate an empty Backend_State because it was removed manually
		backend_state = Backend_address_pool{}
	}
	fmt.Printf("\n+++++++++++++++++++++ state.Backend_address_pool Content after Read :\n %+v ",backend_state)
	
	// Generate resource state struct
	var result = WebappBinding{
		Name:                 types.String{Value: webappBindingName},
		Agw_name:             state.Agw_name,
		Agw_rg:               state.Agw_rg,
		Backend_address_pool: backend_state,
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


	//Verify if the agw already contains the wanted element
	backend_plan := plan.Backend_address_pool

	//create and map the new Backend pool element (backend_json) object from the plan (backend_plan)
	backend_json := createBackendAddressPool(backend_plan)

	//check if the backend name in the plan and state are different
	if backend_plan.Name.Value == state.Backend_address_pool.Name.Value {
		//it is about backend update with the same name
		//so we remove the old one 
		removeBackendAddressPoolElement(&gw, backend_json.Name)
	}else{
		// it's about backend update with a new name
		// we have to check if the new backend name is already used
		if checkBackendAddressPoolElement(gw, backend_json.Name) {
			//this is an error. issue an exit error.
			resp.Diagnostics.AddError(
				"Unable to update the app gateway. The Backend Adresse pool name : "+ backend_json.Name+" already exists.",
				" Please, modify the name.",
			)
			return
		}
		//remove the old backend (old name) from the gateway
		removeBackendAddressPoolElement(&gw, state.Backend_address_pool.Name.Value)
	}

	//add the new one
	gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, backend_json)
	//and update the gateway
	gw_response, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)

	//verify if the backend address pool is added to the gateway, otherwise exit error
	if code != 200 {
		// Error  - backend address pool wasn't added to the app gateway
		resp.Diagnostics.AddError(
			"Unable to update Backend Address pool ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Backend Address pool Name doesn't exist in the response of the app gateway",
		)
		return
	}

	var backend_state Backend_address_pool	
	// in the Read method, the number of fqdns and Ip in a Backendpool should be calculated from the json object and not the plan or state,
	// because the purpose of the read is to see if there is a difference between the real element and the satate stored localy.
		
	index := getBackendAddressPoolElementKey(gw, backend_json.Name)
	backend_json2 := gw.Properties.BackendAddressPools[index]
	nb_BackendAddresses := len(backend_json2.Properties.BackendAddresses)
	nb_Fqdns := 0
	for i := 0; i < nb_BackendAddresses; i++ {
		if (backend_json2.Properties.BackendAddresses[i].Fqdn != "") && (&backend_json2.Properties.BackendAddresses[i].Fqdn != nil) {
			nb_Fqdns++
		} 
	}
	nb_IpAddress := nb_BackendAddresses - nb_Fqdns

	//generate BackendState
	backend_state = generateBackendAddressPoolState(gw_response, backend_json.Name,nb_Fqdns,nb_IpAddress)
	
	// Generate resource state struct
	var result = WebappBinding{
		Name:                 state.Name,
		Agw_name:             types.String{Value: gw_response.Name},
		Agw_rg:               state.Agw_rg,
		Backend_address_pool: backend_state,
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
	// Get backend address pool name from state
	backend_name := state.Backend_address_pool.Name.Value
	
	//Get the agw
	resourceGroupName := state.Agw_rg.Value
	applicationGatewayName := state.Agw_name.Value
	gw := getGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, r.p.token.Access_token)
	/*
	//test if the backend address pool doen't exist in the gateway, then it is an error
	if !checkBackendAddressPoolElement(gw, backend_name) {
		// Error  - the non existance of backend_plan address pool name must stop execution
		resp.Diagnostics.AddError(
			"Unable to delete Backend Address pool",
			"Backend Address pool Name doesn't exist in the app gateway. ### Definitely, it was removed manually###",
		)
		return
	}*/

	//remove the backend from the gw
	removeBackendAddressPoolElement(&gw, backend_name)

	//and update the gateway
	_, error_json, code := updateGW(r.p.AZURE_SUBSCRIPTION_ID, resourceGroupName, applicationGatewayName, gw, r.p.token.Access_token)
	//verify if the backend address pool is added to the gateway, otherwise exit error
	if code != 200 {
		// Error  - backend address pool wasn't added to the app gateway
		resp.Diagnostics.AddError(
			"Unable to delete Backend Address pool ######## API response = "+fmt.Sprint(code)+"\n"+error_json,
			"Backend Address pool Name doesn't exist in the response of the app gateway",
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


func checkElementName(gw ApplicationGateway, plan WebappBinding) (string,bool){
	//This function allows to check if an element name in the required new configuration (plan WebappBinding) already exist in the gw.
	//if so, the provider has to stop executing and issue an exit error
	exist := false

	//Create new var for all configurations
	backend_plan := plan.Backend_address_pool
	if checkBackendAddressPoolElement(gw, backend_plan.Name.Value) {
		exist = true 
	}
	return backend_plan.Name.Value,exist
}

//Backend pool operations
func createBackendAddressPool(backend_plan Backend_address_pool) BackendAddressPool{
	backend_json := BackendAddressPool{
		Name: backend_plan.Name.Value,
		Properties: struct {
			ProvisioningState string "json:\"provisioningState,omitempty\""
			BackendAddresses  []struct {
				Fqdn      string "json:\"fqdn,omitempty\""
				IPAddress string "json:\"ipAddress,omitempty\""
			} "json:\"backendAddresses\""
			RequestRoutingRules []struct {
				ID string "json:\"id,omitempty\""
			} "json:\"requestRoutingRules,omitempty\""
		}{},
		Type: "Microsoft.Network/applicationGateways/backendAddressPools",
	}
	length := len(backend_plan.Fqdns) + len(backend_plan.Ip_addresses)

	//If there is no fqdn nor IPaddress for the backend pool, initialize the BackendAddresses to nil to avoid a terraform provider error when making the state
	if length == 0 {
		backend_json.Properties.BackendAddresses = nil
	} else {
		backend_json.Properties.BackendAddresses = make([]struct {
			Fqdn      string "json:\"fqdn,omitempty\""
			IPAddress string "json:\"ipAddress,omitempty\""
		}, length)
	}
	for i := 0; i < len(backend_plan.Fqdns); i++ {
		backend_json.Properties.BackendAddresses[i].Fqdn = backend_plan.Fqdns[i].Value
	}
	for i := 0; i < len(backend_plan.Ip_addresses); i++ {
		backend_json.Properties.BackendAddresses[i+len(backend_plan.Fqdns)].IPAddress = backend_plan.Ip_addresses[i].Value
	}
	// add the backend to the agw and update the agw
	return backend_json
}
func generateBackendAddressPoolState(gw ApplicationGateway, backendAddressPoolName string,nb_Fqdns int,nb_IpAddress int) Backend_address_pool {
	// we have to give the nb_Fqdns and nb_IpAddress in order to make this function reusable in create, read and update method
	index := getBackendAddressPoolElementKey(gw, backendAddressPoolName)
	backend_json := gw.Properties.BackendAddressPools[index]
	// log the added backend address pool
	//tflog.Trace(ctx, "created BackendAddressPool", "BackendAddressPool ID", backend_json.ID)

	// Map response body to resource schema attribute
	backend_state := Backend_address_pool{
		Name:         types.String{Value: backend_json.Name},
		Id:           types.String{Value: backend_json.ID},
		Fqdns:        []types.String{},
		Ip_addresses: []types.String{},
	}
	
	//fmt.Println("------------------ The number nb_Fqdns is:", nb_Fqdns)
	if nb_Fqdns != 0 {
		backend_state.Fqdns = make([]types.String, nb_Fqdns)
	} else {
		backend_state.Fqdns = nil
	}
	//fmt.Println("------------------ The number nb_IpAddress is:", nb_IpAddress)

	if nb_IpAddress != 0 {
		backend_state.Ip_addresses = make([]types.String, nb_IpAddress)
	} else {
		backend_state.Ip_addresses = nil
	}

	for j := 0; j < nb_Fqdns; j++ {
		backend_state.Fqdns[j] = types.String{Value: backend_json.Properties.BackendAddresses[j].Fqdn}
	}
	for j := 0; j < nb_IpAddress; j++ {
		backend_state.Ip_addresses[j] = types.String{Value: backend_json.Properties.BackendAddresses[j+nb_Fqdns].IPAddress}
	}

	return backend_state
}
func checkBackendAddressPoolElement(gw ApplicationGateway, backendAddressPoolName string) bool {
	exist := false
	//fmt.Println("OOOOOOOO  looking for =", backendAddressPoolName)
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		//fmt.Println("wwwwwwwwww  =", gw.Properties.BackendAddressPools[i].Name)
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			//gw.Properties.BackendAddressPools =append(gw.Properties.BackendAddressPools[:i], gw.Properties.BackendAddressPools[i+1:]...)
			exist = true
		}
	}
	//fmt.Println("ww         Exist =",exist)
	return exist
}
func removeBackendAddressPoolElement(gw *ApplicationGateway, backendAddressPoolName string) {
	//removed := false
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools[:i], gw.Properties.BackendAddressPools[i+1:]...)
			//removed = true
		}
	}
	//fmt.Println("#############################removed =", removed)
}
func getBackendAddressPoolElementKey(gw ApplicationGateway, backendAddressPoolName string) int {
	key := -1
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- {
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName {
			key = i
		}
	}
	return key
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
