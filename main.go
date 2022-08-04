package main

import (
	"azurerm_agw/azurermagw"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)


func main() {
	
	/*for i := 0; i < 3; i++ {
		fmt.Print("##################################")
	}*/
	tfsdk.Serve(context.Background(), azurermagw.New, tfsdk.ServeOpts{
		Name: "azurermagw",
	})
/*
	AZURE_TENANT_ID:="*******************************"
	AZURE_CLIENT_ID :="*****************************"
	AZURE_CLIENT_SECRET :="******************************"
	AZURE_SUBSCRIPTION_ID :="*********************************"
	
	resourceGroupName:= "shared-app-gateway"
	applicationGatewayName := "default-app-gateway-mahmoud"
	token := getToken(AZURE_CLIENT_ID,AZURE_CLIENT_SECRET,AZURE_TENANT_ID)
	var gw azurermagw.ApplicationGateway
	gw = getGW(AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName,token.Access_token)
	fmt.Print("##################################",
				"\nGateway Name =",gw.Name,				
				"\nGateway Properties GatewayIPConfigurations[0].Name=",gw.Properties.GatewayIPConfigurations[0].Name,
				"\nGateway Name =",gw.Properties.SslCertificates[0].Name,
				"\nGateway Name =",gw.Properties.BackendHTTPSettingsCollection[0].Name,
				"\nGateway Name =",gw.Properties.BackendAddressPools[2].Name,"\n")


	
	// create a new backendAddressPool 
	mahmoud := azurermagw.BackendAddressPools{
		Name: "mahmoud-backendAddressPool-name",
		//ID:   "qlsdjflqsdjfqsd",
		//Etag: "Etag-string",
		Properties: struct {
			ProvisioningState string "json:\"provisioningState,omitempty\""
			BackendAddresses  []struct {
				Fqdn      string "json:\"fqdn,omitempty\""
				IPAddress string "json:\"ipAddress,omitempty\""
			} "json:\"backendAddresses\""
			RequestRoutingRules []struct {
				ID string "json:\"id\""
			} "json:\"requestRoutingRules,omitempty\""
		}{},
		Type: "Microsoft.Network/applicationGateways/backendAddressPools",
	}
	
	//fmt.Printf("\nGateway Identity = %+v\n",mahmoud)
	//mahmoud.Properties.BackendAddresses = make([]struct{Fqdn string "json:\"fqdn,omitempty\""; IPAddress string "json:\"ipAddress,omitempty\""}, 2)
	mahmoud.Properties.BackendAddresses = make([]struct{Fqdn string "json:\"fqdn,omitempty\""; IPAddress string "json:\"ipAddress,omitempty\""}, 2)
	
	mahmoud.Properties.BackendAddresses[0].Fqdn ="fqdn.mahmoud"
	mahmoud.Properties.BackendAddresses[1].IPAddress = "10.2.3.3"
	
	//add the newly created backendAddressPool to the gateway
	//gw.Properties.BackendAddressPools = append(gw.Properties.BackendAddressPools, mahmoud)
	//updateGW(AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName,gw,token.Access_token)	
	
	//"name": "mahmoud-backendAddressPool-name"
	removeBackendAddressPoolElement(&gw,"mahmoud-backendAddressPool-name")
	//printGWtoFile(gw,"gw-removed.json")
	updateGW(AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName,gw,token.Access_token)
*/
}
func removeBackendAddressPoolElement(gw *azurermagw.ApplicationGateway,backendAddressPoolName string ){
	removed := false
	for i := len(gw.Properties.BackendAddressPools) - 1; i >= 0; i-- { 
		if gw.Properties.BackendAddressPools[i].Name == backendAddressPoolName { 
			gw.Properties.BackendAddressPools =append(gw.Properties.BackendAddressPools[:i], gw.Properties.BackendAddressPools[i+1:]...) 
			removed=true
		}
	}	
	fmt.Println("#############################removed =",removed)
}

func updateGW(subscriptionId string,resourceGroupName string,applicationGatewayName string,gw azurermagw.ApplicationGateway, token string){
	requestURI := "https://management.azure.com/subscriptions/"+subscriptionId+"/resourceGroups/"+ 
	resourceGroupName+"/providers/Microsoft.Network/applicationGateways/"+applicationGatewayName+"?api-version=2021-08-01"
	payloadBytes, err := json.Marshal(gw)
	if err != nil {
		// handle err
	}

////////print json gw
	//rs := string(payloadBytes)
	//ress, err := PrettyString(rs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("agw_before.json||||||||||||||||||||||||||||||||||||||||||||||||||||||||")
	//printToFile(ress,"agw_before.json")
	fmt.Println("agw_before.json||||||||||||||||||||||||||||||||||||||||||||||||||||||||")
	
////////////

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
	//responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
	
	//responseString := string(responseData)
    //res, err := PrettyString(responseString)
    if err != nil {
        log.Fatal(err)
    }
    //printToFile(res,"agw_after.json")
	return 
}
/*
func monMain(){
	//AZURE_TENANT_ID :="*****************************"
	//AZURE_CLIENT_ID :="**********************"
	//AZURE_CLIENT_SECRET :="**************************************"
	AZURE_SUBSCRIPTION_ID :="*********************************"
	//AZURE_SUBSCRIPTION_ID :="***********************************"
	resourceGroupName:= "shared-app-gateway"
	applicationGatewayName := "default-app-gateway-mahmoud"
	//applicationGatewayName := "dev-app-gateway"
	//token := getToken()
	//fmt.Printf("###################################\n%s\n",token.Access_token)
	var gw azurermagw.ApplicationGateway
	gw = getGW(AZURE_SUBSCRIPTION_ID,resourceGroupName,applicationGatewayName,token.Access_token)
	
	var ide azurermagw.Identity
	ide = gw.Identity
	//printToFile(gw,"agw.txt")
	fmt.Printf("\nGateway Identity = %+v",ide)//.UserAssignedIdentities,
	fmt.Print("##################################",
				"\nGateway Name =",gw.Name,				
				"\nGateway Properties GatewayIPConfigurations[0].Name=",gw.Properties.GatewayIPConfigurations[0].Name,
				"\nGateway Name =",gw.Properties.SslCertificates[0].Name,
				"\nGateway Name =",gw.Properties.BackendHTTPSettingsCollection[0].Name,
				"\nGateway Name =",gw.Properties.BackendAddressPools[2].Name,"\n")
}*/
func printGWtoFile(gw azurermagw.ApplicationGateway, fileName string){
	payloadBytes, err := json.Marshal(gw)
	if err != nil {
		// handle err
	}

	////////print json gw
	rs := string(payloadBytes)
	//fmt.Printf(responseString)
	ress, err := PrettyString(rs)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("agw_before.json||||||||||||||||||||||||||||||||||||||||||||||||||||||||")
	printToFile(ress,fileName)
}
func printToFile(str string, fileName string){
	file, err := os.Create(fileName)
    if err != nil {
        log.Fatal(err)
    }
    mw := io.MultiWriter(os.Stdout, file)
    fmt.Fprintln(mw, str)
}
func getGW(subscriptionId string,resourceGroupName string,applicationGatewayName string, token string)(azurermagw.ApplicationGateway){
	requestURI := "https://management.azure.com/subscriptions/"+subscriptionId+"/resourceGroups/"+ 
	resourceGroupName+"/providers/Microsoft.Network/applicationGateways/"+applicationGatewayName+"?api-version=2021-08-01"
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
	var agw azurermagw.ApplicationGateway
	err = json.Unmarshal(responseData, &agw)
  
    if err != nil {  
        fmt.Println(err)
    }
	
	return agw
}
func PrettyString(str string) (string, error) {
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
        return "", err
    }
    return prettyJSON.String(), nil
}
func getToken(client_id string,client_secret string,tenant_id string)(azurermagw.Token){
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	params.Add("client_id", client_id)
	params.Add("client_secret",client_secret)
	params.Add("resource", `https://management.azure.com/`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/"+tenant_id+"/oauth2/token", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
	/*
    responseString := string(responseData)
    //fmt.Printf(responseString)

	res, err := PrettyString(responseString)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)*/

	var token azurermagw.Token
	err = json.Unmarshal(responseData, &token)
  
    if err != nil {  
        fmt.Println(err)
    }
	//fmt.Print("#########\ntoken.Token_type: ",token.Token_type,"\ntoken.Expires_in: ",token.Expires_in,"\ntoken.Ext_expires_in: ",token.Ext_expires_in,"\ntoken.Expires_on: ",token.Expires_on,"\ntoken.Not_before: ",token.Not_before,"\ntoken.Resource: ",token.Resource,"\n")
	return token
}
func restCall(token string){
	//token := os.Getenv("TOKEN")
		
	//fmt.Printf("Token = %s",token)

	// curl -X GET https://management.azure.com/subscriptions/*******************/resourcegroups?api-version=2020-09-01 -H "Authorization: Bearer {LONG_STRING_HERE}" -H "Content-type: application/json"
	requestURI := "https://management.azure.com/subscriptions/*******************/resourcegroups?api-version=2020-09-01"
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

    responseString := string(responseData)
    //fmt.Printf(responseString)

	res, err := PrettyString(responseString)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
}
