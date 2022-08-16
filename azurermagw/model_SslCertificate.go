package azurermagw

import (
	//"fmt"
	//"fmt"
	//"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SslCertificate struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState string `json:"provisioningState,omitempty"`
		PublicCertData    string `json:"publicCertData,omitempty"`
		KeyVaultSecretID  string `json:"keyVaultSecretId,omitempty"`
		Password          string `json:"password,omitempty"`
		HTTPListeners     *[]struct {
			ID string `json:"id"`
		} `json:"httpListeners"`
	} `json:"properties"`
	Type string `json:"type"`
} 

type Ssl_certificate struct {
	Name         								types.String	`tfsdk:"name"`	
	Id           								types.String	`tfsdk:"id"`
	Key_vault_secret_id							types.String	`tfsdk:"key_vault_secret_id"`
	Data	                       				types.String	`tfsdk:"data"`								
	Password									types.String	`tfsdk:"password"`
}

func createSslCertificate(sslCertificate_plan Ssl_certificate,AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) (SslCertificate, string, string){	
	sslCertificate_json := SslCertificate{
		Name:       sslCertificate_plan.Name.Value,
		//ID:         "",
		//Etag:       "",
		Properties: struct{
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			PublicCertData string "json:\"publicCertData,omitempty\""; 
			KeyVaultSecretID string "json:\"keyVaultSecretId,omitempty\""; 
			Password string "json:\"password,omitempty\""; 
			HTTPListeners *[]struct{
				ID string "json:\"id\""
			} "json:\"httpListeners\""
		}{
			//KeyVaultSecretID: sslCertificate_plan.Key_vault_secret_id.Value,
		},
		Type: "Microsoft.Network/applicationGateways/sslCertificates",
	}
	var error_exclusivity string
	var error_password string
	//there is 2 constraints for SSLCertificate we have to check
	//   1) Data and Key_vault_secret_id are optional but one of them has to be provided
	//   2) If Data is provided, Password is required
	if sslCertificate_plan.Key_vault_secret_id.Value != "" {
		if sslCertificate_plan.Data.Value != "" {
			//both are provided
			error_exclusivity = "fatal-both-exist"
		}else{
			//only Key_vault_secret_id is provided
			sslCertificate_json.Properties.KeyVaultSecretID = sslCertificate_plan.Key_vault_secret_id.Value
		}
	}else{
		if sslCertificate_plan.Data.Value != "" {
			//only data is provided. check the password
			if sslCertificate_plan.Password.Value != "" {
				sslCertificate_json.Properties.Password = sslCertificate_plan.Password.Value 
			}else{
				error_password = "fatal"
			}
		}else{
			//both are empty
			error_exclusivity = "fatal-both-miss"
		}
	}
	return sslCertificate_json,error_exclusivity,error_password
}
func generateSslCertificateState(gw ApplicationGateway, SslCertificateName string) Ssl_certificate {
	//retrieve json element from gw
	index := getSslCertificateElementKey(gw, SslCertificateName)
	sslCertificate_json := gw.Properties.SslCertificates[index]
	
	// Map response body to resource schema attribute
	var sslCertificate_state Ssl_certificate
	sslCertificate_state = Ssl_certificate{
		Name:                types.String{Value: sslCertificate_json.Name},
		Id:                  types.String{Value: sslCertificate_json.ID},
		//Key_vault_secret_id: types.String{},
		//Data:                types.String{},
		//Password:            types.String{},
	}
	if sslCertificate_json.Properties.KeyVaultSecretID != "" {
		sslCertificate_state.Key_vault_secret_id = types.String{Value: sslCertificate_json.Properties.KeyVaultSecretID}
		sslCertificate_state.Data.Null = true
		sslCertificate_state.Password.Null = true
	}else{
		sslCertificate_state.Key_vault_secret_id.Null = true
		sslCertificate_state.Data = types.String{Value: sslCertificate_json.Properties.PublicCertData}
		sslCertificate_state.Password = types.String{Value: sslCertificate_json.Properties.Password}
	}
	
	return sslCertificate_state
}
func getSslCertificateElementKey(gw ApplicationGateway, SslCertificateName string) int {
	key := -1
	for i := len(gw.Properties.SslCertificates) - 1; i >= 0; i-- {
		if gw.Properties.SslCertificates[i].Name == SslCertificateName {
			key = i
		}
	}
	return key
}
func checkSslCertificateElement(gw ApplicationGateway, SslCertificateName string) bool {
	exist := false
	for i := len(gw.Properties.SslCertificates) - 1; i >= 0; i-- {
		if gw.Properties.SslCertificates[i].Name == SslCertificateName {
			exist = true
		}
	}
	return exist
}
func removeSslCertificateElement(gw *ApplicationGateway, SslCertificateName string) {
	for i := len(gw.Properties.SslCertificates) - 1; i >= 0; i-- {
		if gw.Properties.SslCertificates[i].Name == SslCertificateName {
			gw.Properties.SslCertificates = append(gw.Properties.SslCertificates[:i], gw.Properties.SslCertificates[i+1:]...)
		}
	}
}