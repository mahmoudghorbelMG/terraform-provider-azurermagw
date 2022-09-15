package azurermagw

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	token                 *Token
	AZURE_SUBSCRIPTION_ID string
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{ 
			"azure_client_id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Client ID which should be used (only Service principal is supported actually).",
			},
			"azure_client_secret": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				MarkdownDescription: "The Client Secret which should be used.",
			},
			"azure_tenant_id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Tenant ID which should be used.",
			},
			"azure_subscription_id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Subscription ID which should be used.",
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct { 
	AZURE_CLIENT_ID       types.String `tfsdk:"azure_client_id"`
	AZURE_CLIENT_SECRET   types.String `tfsdk:"azure_client_secret"`
	AZURE_TENANT_ID       types.String `tfsdk:"azure_tenant_id"`
	AZURE_SUBSCRIPTION_ID types.String `tfsdk:"azure_subscription_id"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a AZURE_CLIENT_ID to the provider
	var AZURE_CLIENT_ID string
	if config.AZURE_CLIENT_ID.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create Azure client",
			"Cannot use unknown value as AZURE_CLIENT_ID",
		)
		return
	}

	if config.AZURE_CLIENT_ID.Null {
		AZURE_CLIENT_ID = os.Getenv("AZURE_CLIENT_ID")
	} else {
		AZURE_CLIENT_ID = config.AZURE_CLIENT_ID.Value
	}


	// User must provide a AZURE_CLIENT_SECRET to the provider
	var AZURE_CLIENT_SECRET string
	if config.AZURE_CLIENT_SECRET.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create Azure client",
			"Cannot use unknown value as AZURE_CLIENT_SECRET",
		)
		return
	}

	if config.AZURE_CLIENT_SECRET.Null {
		AZURE_CLIENT_SECRET = os.Getenv("AZURE_CLIENT_SECRET")
	} else {
		AZURE_CLIENT_SECRET = config.AZURE_CLIENT_SECRET.Value
	}


	// User must provide a AZURE_TENANT_ID to the provider
	var AZURE_TENANT_ID string
	if config.AZURE_TENANT_ID.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create Azure client",
			"Cannot use unknown value as AZURE_TENANT_ID",
		)
		return
	}

	if config.AZURE_TENANT_ID.Null {
		AZURE_TENANT_ID = os.Getenv("AZURE_TENANT_ID")
	} else {
		AZURE_TENANT_ID = config.AZURE_TENANT_ID.Value
	}

	// User must provide a AZURE_SUBSCRIPTION_ID to the provider
	var AZURE_SUBSCRIPTION_ID string
	if config.AZURE_SUBSCRIPTION_ID.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create Azure client",
			"Cannot use unknown value as AZURE_SUBSCRIPTION_ID",
		)
		return
	}

	if config.AZURE_SUBSCRIPTION_ID.Null {
		AZURE_SUBSCRIPTION_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	} else {
		AZURE_SUBSCRIPTION_ID = config.AZURE_SUBSCRIPTION_ID.Value
	}

	/*if AZURE_CLIENT_ID == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find AZURE_CLIENT_ID",
			"AZURE_CLIENT_ID cannot be an empty string",
		)
		return
	}	
	if AZURE_CLIENT_SECRET == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find AZURE_CLIENT_SECRET",
			"AZURE_CLIENT_SECRET cannot be an empty string",
		)
		return
	}
	if AZURE_TENANT_ID == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find AZURE_TENANT_ID",
			"AZURE_TENANT_ID cannot be an empty string",
		)
		return
	}
	if AZURE_SUBSCRIPTION_ID == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find AZURE_SUBSCRIPTION_ID",
			"AZURE_SUBSCRIPTION_ID cannot be an empty string",
		)
		return
	}*/
	
	//check if we can get token after a successful login
	cmd := exec.Command("az", "account", "get-access-token")
	token_json, err := cmd.Output()
	var login bool
	if err!= nil {
		fmt.Println("\n================== login = false =======================")
		login = false
	}else{
		fmt.Println("\n================== login = true =======================\n",string(token_json))		
		login = true 
	}

	var token Token		
	if login {
		err = json.Unmarshal(token_json, &token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nuuuuuuuuuuuuuuuu token with login uuuuuuuuuuuuuuuu\n",token)
	}else{
		if AZURE_SUBSCRIPTION_ID == "" || AZURE_TENANT_ID == "" || AZURE_CLIENT_SECRET == "" || AZURE_CLIENT_ID == "" {
			// Error vs warning - empty value must stop execution
			resp.Diagnostics.AddError(
			"Unable to find AZURE_SUBSCRIPTION_ID or AZURE_TENANT_ID or AZURE_CLIENT_SECRET or AZURE_CLIENT_ID",
			"Please provide all these values (in tf file or via Environment variables) or connect to azure through az cli.",)
			return
		}else{
			token = getToken(AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)	
			fmt.Println("\nuuuuuuuuuuuuuuuu token with ENV uuuuuuuuuuuuuuuu\n",token)
		}
		
	}
	// create Token
	p.token = &token
	p.AZURE_SUBSCRIPTION_ID = AZURE_SUBSCRIPTION_ID
	//resp.Diagnostics.AddWarning("################TOKEN############### : ",p.token.Access_token)

	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		//"hashicups_order": resourceOrderType{},
		"azurermagw_binding_service": resourceBindingServiceType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		//"hashicups_coffees": dataSourceCoffeesType{},
	}, nil
}

// Get Token to call Azure Rest API
func getToken(client_id string, client_secret string, tenant_id string) Token {
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	params.Add("client_id", client_id)
	params.Add("client_secret", client_secret)
	params.Add("resource", `https://management.azure.com/`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/"+tenant_id+"/oauth2/token", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read and put the json response in byte format
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal the json format response to a token struct
	var token Token
	err = json.Unmarshal(responseData, &token)
	if err != nil {
		log.Fatal(err)
	}
	return token
}
