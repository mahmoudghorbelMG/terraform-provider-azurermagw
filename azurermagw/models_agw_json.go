package azurermagw

type Token struct {
	// defining struct variables
	Token_type     string `json:"token_type"`
	Expires_in     string `json:"expires_in"`
	Ext_expires_in string `json:"ext_expires_in"`
	Expires_on     string `json:"expires_on"`
	Not_before     string `json:"not_before"`
	Resource       string `json:"resource"`
	Access_token   string `json:"access_token"`
}
type TokenLogin struct {
	// defining struct variables
	Access_token   	string `json:"accessToken"`
	Expires_on     	string `json:"expiresOn"`
	Subscription_id	string `json:"subscription"`
	Token_type     	string `json:"tokenType"`
	Tenant			string `json:"tenant"`
}

type ApplicationGateway struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Etag     string `json:"etag"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Tags     struct {
	} `json:"tags"`
	Identity *struct { //Identity `json:"identity,omitempty"`
		Type                   string      `json:"type,omitempty"`
		UserAssignedIdentities interface{} `json:"userAssignedIdentities,omitempty"`
	} `json:"identity"`
	Properties struct {
		AuthenticationCertificates []struct {
			Id         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				Data string `json:"data"`
			} `json:"properties"`
		} `json:"authenticationCertificates"`
		AutoscaleConfiguration *struct {
			MaxCapacity int `json:"maxCapacity"`
			MinCapacity int `json:"minCapacity"`
		} `json:"autoscaleConfiguration"`
		BackendAddressPools 			[]BackendAddressPool `json:"backendAddressPools,omitempty"` 
		BackendHTTPSettingsCollection 	[]BackendHTTPSettings `json:"backendHttpSettingsCollection,omitempty"`
		BackendSettingsCollection 		[]struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				HostName                       string `json:"hostName"`
				PickHostNameFromBackendAddress string `json:"pickHostNameFromBackendAddress"`
				Port                           string `json:"port"`
				Probe                          struct {
					ID string `json:"id"`
				} `json:"probe"`
				Protocol                string `json:"protocol"`
				Timeout                 string `json:"timeout"`
				TrustedRootCertificates []struct {
					ID string `json:"id"`
				} `json:"trustedRootCertificates"`
			} `json:"properties"`
		} `json:"backendSettingsCollection"`
		CustomErrorConfigurations []struct {
			CustomErrorPageURL string `json:"customErrorPageUrl"`
			StatusCode         string `json:"statusCode"`
		} `json:"customErrorConfigurations"`
		EnableFips     bool `json:"enableFips,omitempty"`
		EnableHTTP2    bool `json:"enableHttp2,omitempty"`
		FirewallPolicy *struct {
			ID string `json:"id"`
		} `json:"firewallPolicy"`
		ForceFirewallPolicyAssociation bool `json:"forceFirewallPolicyAssociation,omitempty"`
		FrontendIPConfigurations       []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Etag       string `json:"etag"`
			Type       string `json:"type"`
			Properties struct {
				PrivateIPAddress          string `json:"privateIPAddress,omitempty"`
				ProvisioningState         string `json:"provisioningState"`
				PrivateIPAllocationMethod string `json:"privateIPAllocationMethod"`
				PublicIPAddress           struct {
					ID string `json:"id"`
				} `json:"publicIPAddress"`
				HTTPListeners []struct {
					ID string `json:"id"`
				} `json:"httpListeners"`
				Subnet *struct {
					ID string `json:"id"`
				} `json:"subnet"`
			} `json:"properties"`
		} `json:"frontendIPConfigurations"`
		FrontendPorts []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Etag       string `json:"etag"`
			Properties struct {
				ProvisioningState string `json:"provisioningState"`
				Port              int    `json:"port"`
				HTTPListeners     []struct {
					ID string `json:"id"`
				} `json:"httpListeners"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"frontendPorts"`
		GatewayIPConfigurations []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Etag       string `json:"etag"`
			Properties struct {
				ProvisioningState string `json:"provisioningState"`
				Subnet            struct {
					ID string `json:"id"`
				} `json:"subnet"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"gatewayIPConfigurations"`
		GlobalConfiguration *struct {
			EnableRequestBuffering  bool `json:"enableRequestBuffering"`
			EnableResponseBuffering bool `json:"enableResponseBuffering"`
		} `json:"globalConfiguration"`
		HTTPListeners [] HTTPListener `json:"httpListeners,omitempty"`
		Listeners []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				FrontendIPConfiguration struct {
					ID string `json:"id"`
				} `json:"frontendIPConfiguration"`
				FrontendPort struct {
					ID string `json:"id"`
				} `json:"frontendPort"`
				Protocol       string `json:"protocol"`
				SslCertificate *struct {
					ID string `json:"id"`
				} `json:"sslCertificate"`
				SslProfile *struct {
					ID string `json:"id"`
				} `json:"sslProfile"`
			} `json:"properties"`
		} `json:"listeners"`
		LoadDistributionPolicies []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				LoadDistributionAlgorithm string `json:"loadDistributionAlgorithm"`
				LoadDistributionTargets   []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					Properties struct {
						BackendAddressPool struct {
							ID string `json:"id"`
						} `json:"backendAddressPool"`
						WeightPerServer string `json:"weightPerServer"`
					} `json:"properties"`
				} `json:"loadDistributionTargets"`
			} `json:"properties"`
		} `json:"loadDistributionPolicies"`
		PrivateLinkConfigurations []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				IPConfigurations []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					Properties struct {
						Primary                   bool   `json:"primary"`
						PrivateIPAddress          string `json:"privateIPAddress,omitempty"`
						PrivateIPAllocationMethod string `json:"privateIPAllocationMethod"`
						Subnet                    struct {
							ID string `json:"id"`
						} `json:"subnet"`
					} `json:"properties"`
				} `json:"ipConfigurations"`
			} `json:"properties"`
		} `json:"privateLinkConfigurations"`
		Probes []Probe_json `json:"probes"`
		RedirectConfigurations []RedirectConfiguration `json:"redirectConfigurations,omitempty"`
		RequestRoutingRules []RequestRoutingRule `json:"requestRoutingRules,omitempty"`
		RewriteRuleSets []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Etag       string `json:"etag"`
			Properties struct {
				ProvisioningState string `json:"provisioningState"`
				RewriteRules      []struct {
					RuleSequence int           `json:"ruleSequence"`
					Conditions   []interface{} `json:"conditions"`
					Name         string        `json:"name"`
					ActionSet    struct {
						RequestHeaderConfigurations []struct {
							HeaderName  string `json:"headerName"`
							HeaderValue string `json:"headerValue"`
						} `json:"requestHeaderConfigurations"`
						ResponseHeaderConfigurations []struct {
							HeaderName  string `json:"headerName"`
							HeaderValue string `json:"headerValue"`
						} `json:"responseHeaderConfigurations"`
						URLConfiguration struct {
							ModifiedPath        string `json:"modifiedPath"`
							ModifiedQueryString string `json:"modifiedQueryString"`
							Reroute             bool   `json:"reroute"`
						} `json:"urlConfiguration"`
					} `json:"actionSet"`
				} `json:"rewriteRules"`
				RequestRoutingRules []struct {
					ID string `json:"id"`
				} `json:"requestRoutingRules"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"rewriteRuleSets"`
		RoutingRules []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				BackendAddressPool struct {
					ID string `json:"id"`
				} `json:"backendAddressPool"`
				BackendSettings struct {
					ID string `json:"id"`
				} `json:"backendSettings"`
				Listener struct {
					ID string `json:"id"`
				} `json:"listener"`
				RuleType string `json:"ruleType"`
			} `json:"properties"`
		} `json:"routingRules"`
		Sku struct {
			Name     string `json:"name"`
			Tier     string `json:"tier"`
			Capacity int    `json:"capacity"`
		} `json:"sku"`
		SslCertificates []SslCertificate `json:"sslCertificates"`
		SslPolicy struct {
			CipherSuites         []string `json:"cipherSuites,omitempty"`
			DisabledSslProtocols []string `json:"disabledSslProtocols,omitempty"`
			MinProtocolVersion   string   `json:"minProtocolVersion,omitempty"`
			PolicyName           string   `json:"policyName"`
			PolicyType           string   `json:"policyType"`
		} `json:"sslPolicy"`
		SslProfiles []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				ClientAuthConfiguration struct {
					VerifyClientCertIssuerDN bool `json:"verifyClientCertIssuerDN"`
				} `json:"clientAuthConfiguration"`
				SslPolicy struct {
					CipherSuites         []string `json:"cipherSuites"`
					DisabledSslProtocols []string `json:"disabledSslProtocols"`
					MinProtocolVersion   string   `json:"minProtocolVersion"`
					PolicyName           string   `json:"policyName"`
					PolicyType           string   `json:"policyType"`
				} `json:"sslPolicy"`
				TrustedClientCertificates []struct {
					ID string `json:"id"`
				} `json:"trustedClientCertificates"`
			} `json:"properties"`
		} `json:"sslProfiles"`
		TrustedClientCertificates []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				Data string `json:"data"`
			} `json:"properties"`
		} `json:"trustedClientCertificates"`
		TrustedRootCertificates []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				Data             string `json:"data"`
				KeyVaultSecretID string `json:"keyVaultSecretId"`
			} `json:"properties"`
		} `json:"trustedRootCertificates"`
		URLPathMaps []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				DefaultBackendAddressPool struct {
					ID string `json:"id"`
				} `json:"defaultBackendAddressPool"`
				DefaultBackendHTTPSettings struct {
					ID string `json:"id"`
				} `json:"defaultBackendHttpSettings"`
				DefaultLoadDistributionPolicy struct {
					ID string `json:"id"`
				} `json:"defaultLoadDistributionPolicy"`
				DefaultRedirectConfiguration struct {
					ID string `json:"id"`
				} `json:"defaultRedirectConfiguration"`
				DefaultRewriteRuleSet struct {
					ID string `json:"id"`
				} `json:"defaultRewriteRuleSet"`
				PathRules []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					Properties struct {
						BackendAddressPool struct {
							ID string `json:"id"`
						} `json:"backendAddressPool"`
						BackendHTTPSettings struct {
							ID string `json:"id"`
						} `json:"backendHttpSettings"`
						FirewallPolicy *struct {
							ID string `json:"id"`
						} `json:"firewallPolicy"`
						LoadDistributionPolicy struct {
							ID string `json:"id"`
						} `json:"loadDistributionPolicy"`
						Paths                 []string `json:"paths"`
						RedirectConfiguration struct {
							ID string `json:"id"`
						} `json:"redirectConfiguration"`
						RewriteRuleSet struct {
							ID string `json:"id"`
						} `json:"rewriteRuleSet"`
					} `json:"properties"`
				} `json:"pathRules"`
			} `json:"properties"`
		} `json:"urlPathMaps"`
		WebApplicationFirewallConfiguration *struct {
			Enabled            bool   `json:"enabled"`
			MaxRequestBodySize int    `json:"maxRequestBodySize,omitempty"`
			FirewallMode       string `json:"firewallMode,omitempty"`
			RuleSetType        string `json:"ruleSetType,omitempty"`
			RuleSetVersion     string `json:"ruleSetVersion"`
			DisabledRuleGroups []struct {
				RuleGroupName string   `json:"ruleGroupName"`
				Rules         []string `json:"rules"`
			} `json:"disabledRuleGroups"`
			Exclusions *[]struct {
				MatchVariable         string `json:"matchVariable"`
				Selector              string `json:"selector"`
				SelectorMatchOperator string `json:"selectorMatchOperator"`
			} `json:"exclusions"`
			RequestBodyCheck       bool `json:"requestBodyCheck"`
			MaxRequestBodySizeInKb int  `json:"maxRequestBodySizeInKb"`
			FileUploadLimitInMb    int  `json:"fileUploadLimitInMb"`
		} `json:"webApplicationFirewallConfiguration"`
		Zones                      []string      `json:"zones,omitempty"`
		ProvisioningState          string        `json:"provisioningState"`
		ResourceGUID               string        `json:"resourceGuid"`
		OperationalState           string        `json:"operationalState"`
		PrivateEndpointConnections []interface{} `json:"privateEndpointConnections"`
	} `json:"properties"`
}
