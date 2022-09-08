package azurermagw

import (
	//"fmt"
	//"fmt"
	//"strings"

	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SslCertificate struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Etag       string `json:"etag,omitempty"`
	Properties struct {
		ProvisioningState string `json:"provisioningState,omitempty"`
		Data 			  string `json:"data,omitempty"`
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

func createSslCertificate(sslCertificate_plan Ssl_certificate,AZURE_SUBSCRIPTION_ID string, rg_name string, agw_name string) (SslCertificate){	
	sslCertificate_json := SslCertificate{
		Name:       sslCertificate_plan.Name.Value,
		//ID:         "",
		//Etag:       "",
		Properties: struct{
			ProvisioningState string "json:\"provisioningState,omitempty\""; 
			Data string "json:\"data,omitempty\"";
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
	if sslCertificate_plan.Key_vault_secret_id.Value != "" {		
		//only Key_vault_secret_id is provided
		sslCertificate_json.Properties.KeyVaultSecretID = sslCertificate_plan.Key_vault_secret_id.Value
		//}
	}
	if sslCertificate_plan.Data.Value != "" {
		//only data is provided. check the password	
		// data must be base64 encoded (from azurerm provider)
		// output.ApplicationGatewaySslCertificatePropertiesFormat.Data = utils.String(utils.Base64EncodeIfNot(data))
		sslCertificate_json.Properties.Data = base64EncodeIfNot(sslCertificate_plan.Data.Value)
		sslCertificate_json.Properties.Password = sslCertificate_plan.Password.Value 
	}
	return sslCertificate_json
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
func checkSslCertificateCreate(plan BindingService, gw ApplicationGateway, resp *tfsdk.CreateResourceResponse) bool {
	//there is 2 constraints we have to check for SSLCertificate 
	//   1) Data and Key_vault_secret_id are optional but one of them has to be provided
	//   2) If Data is provided, Password is required
	
	//fatal-both-exist
	if plan.Ssl_certificate.Key_vault_secret_id.Value != "" && plan.Ssl_certificate.Data.Value != "" {
		resp.Diagnostics.AddError(
			"Unable to create binding. In the SSL Certificate ("+plan.Ssl_certificate.Name.Value+") configuration, 2 optional parameters mutually exclusive "+ 
			"are declared: Data and Key_vault_secret_id. Only one has to be set. ",
			"Please, change configuration then retry.",				)
		return true
	}	
	//fatal-both-miss
	if plan.Ssl_certificate.Key_vault_secret_id.Value == "" && plan.Ssl_certificate.Data.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to create binding. In the SSL Certificate  ("+plan.Ssl_certificate.Name.Value+") configuration, both optional parameters mutually exclusive "+ 
			"are missing: Data and Key_vault_secret_id. At least and only one has to be set. ",
			"Please, change configuration then retry.",
			)
		return true
	}
	// miss password
	if plan.Ssl_certificate.Data.Value != "" && plan.Ssl_certificate.Password.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to create binding. In the SSL Certificate  ("+plan.Ssl_certificate.Name.Value+") configuration, Data parameter (pfx file content) "+ 
			"is provided without password. ",
			"Please, add password then retry.",
			)
		return true
	}
	return false
}
func checkSslCertificateUpdate(plan BindingService, gw ApplicationGateway, resp *tfsdk.UpdateResourceResponse) bool {
	//there is 2 constraints we have to check for SSLCertificate 
	//   1) Data and Key_vault_secret_id are optional but one of them has to be provided
	//   2) If Data is provided, Password is required
	
	//fatal-both-exist
	if plan.Ssl_certificate.Key_vault_secret_id.Value != "" && plan.Ssl_certificate.Data.Value != "" {
		resp.Diagnostics.AddError(
			"Unable to update binding. In the SSL Certificate ("+plan.Ssl_certificate.Name.Value+") configuration, 2 optional parameters mutually exclusive "+ 
			"are declared: Data and Key_vault_secret_id. Only one has to be set. ",
			"Please, change configuration then retry.",				)
		return true
	}	
	//fatal-both-miss
	if plan.Ssl_certificate.Key_vault_secret_id.Value == "" && plan.Ssl_certificate.Data.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to update binding. In the SSL Certificate  ("+plan.Ssl_certificate.Name.Value+") configuration, both optional parameters mutually exclusive "+ 
			"are missing: Data and Key_vault_secret_id. At least and only one has to be set. ",
			"Please, change configuration then retry.",
			)
		return true
	}
	// miss password
	if plan.Ssl_certificate.Data.Value != "" && plan.Ssl_certificate.Password.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to update binding. In the SSL Certificate  ("+plan.Ssl_certificate.Name.Value+") configuration, Data parameter (pfx file content) "+ 
			"is provided without password. ",
			"Please, add password then retry.",
			)
		return true
	}
	return false
}
func base64EncodeIfNot(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if base64IsEncoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}
func base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}