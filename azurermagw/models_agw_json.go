package azurermagw

type Token struct {
	// defining struct variables
	Token_type     string `json:"token_type,omitempty"`
	Expires_in     string `json:"expires_in,omitempty"`
	Ext_expires_in string `json:"ext_expires_in,omitempty"`
	Expires_on     string `json:"expires_on,omitempty"`
	Not_before     string `json:"not_before,omitempty"`
	Resource       string `json:"resource,omitempty"`
	Access_token   string `json:"access_token,omitempty"`
}
type TokenLogin struct {
	// defining struct variables
	Access_token   	string `json:"accessToken,omitempty"`
	Expires_on     	string `json:"expiresOn,omitempty"`
	Subscription_id	string `json:"subscription,omitempty"`
	Token_type     	string `json:"tokenType,omitempty"`
	Tenant			string `json:"tenant,omitempty"`
}

type ApplicationGateway struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Etag     string `json:"etag,omitempty"`
	Type     string `json:"type,omitempty"`
	Location string `json:"location,omitempty"`
	Tags     struct {
	} `json:"tags,omitempty"`
	Identity *struct { //Identity `json:"identity,omitempty"`
		Type                   string      `json:"type,omitempty"`
		UserAssignedIdentities interface{} `json:"userAssignedIdentities,omitempty"`
	} `json:"identity,omitempty"`
	Properties struct {
		AuthenticationCertificates []struct {
			Id         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				Data string `json:"data,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"authenticationCertificates,omitempty"`
		AutoscaleConfiguration *struct {
			MaxCapacity int `json:"maxCapacity,omitempty"`
			MinCapacity int `json:"minCapacity,omitempty"`
		} `json:"autoscaleConfiguration,omitempty"`
		BackendAddressPools 			[]BackendAddressPool `json:"backendAddressPools,omitempty"` 
		BackendHTTPSettingsCollection 	[]BackendHTTPSettings `json:"backendHttpSettingsCollection,omitempty"`
		BackendSettingsCollection 		[]struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				HostName                       string `json:"hostName,omitempty"`
				PickHostNameFromBackendAddress string `json:"pickHostNameFromBackendAddress,omitempty"`
				Port                           string `json:"port,omitempty"`
				Probe                          struct {
					ID string `json:"id,omitempty"`
				} `json:"probe,omitempty"`
				Protocol                string `json:"protocol,omitempty"`
				Timeout                 string `json:"timeout,omitempty"`
				TrustedRootCertificates []struct {
					ID string `json:"id,omitempty"`
				} `json:"trustedRootCertificates,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"backendSettingsCollection,omitempty"`
		CustomErrorConfigurations []struct {
			CustomErrorPageURL string `json:"customErrorPageUrl,omitempty"`
			StatusCode         string `json:"statusCode,omitempty"`
		} `json:"customErrorConfigurations,omitempty"`
		EnableFips     bool `json:"enableFips,omitempty"`
		EnableHTTP2    bool `json:"enableHttp2,omitempty"`
		FirewallPolicy *struct {
			ID string `json:"id,omitempty"`
		} `json:"firewallPolicy,omitempty"`
		ForceFirewallPolicyAssociation bool `json:"forceFirewallPolicyAssociation,omitempty"`
		FrontendIPConfigurations       []struct {
			Name       string `json:"name,omitempty"`
			ID         string `json:"id,omitempty"`
			Etag       string `json:"etag,omitempty"`
			Type       string `json:"type,omitempty"`
			Properties struct {
				PrivateIPAddress          string `json:"privateIPAddress,omitempty"`
				ProvisioningState         string `json:"provisioningState,omitempty"`
				PrivateIPAllocationMethod string `json:"privateIPAllocationMethod,omitempty"`
				PublicIPAddress           struct {
					ID string `json:"id,omitempty"`
				} `json:"publicIPAddress,omitempty"`
				HTTPListeners []struct {
					ID string `json:"id,omitempty"`
				} `json:"httpListeners,omitempty"`
				Subnet *struct {
					ID string `json:"id,omitempty"`
				} `json:"subnet,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"frontendIPConfigurations,omitempty"`
		FrontendPorts []struct {
			Name       string `json:"name,omitempty"`
			ID         string `json:"id,omitempty"`
			Etag       string `json:"etag,omitempty"`
			Properties struct {
				ProvisioningState string `json:"provisioningState,omitempty"`
				Port              int    `json:"port,omitempty"`
				HTTPListeners     []struct {
					ID string `json:"id,omitempty"`
				} `json:"httpListeners,omitempty"`
			} `json:"properties,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"frontendPorts,omitempty"`
		GatewayIPConfigurations []struct {
			Name       string `json:"name,omitempty"`
			ID         string `json:"id,omitempty"`
			Etag       string `json:"etag,omitempty"`
			Properties struct {
				ProvisioningState string `json:"provisioningState,omitempty"`
				Subnet            struct {
					ID string `json:"id,omitempty"`
				} `json:"subnet,omitempty"`
			} `json:"properties,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"gatewayIPConfigurations,omitempty"`
		GlobalConfiguration *struct {
			EnableRequestBuffering  bool `json:"enableRequestBuffering,omitempty"`
			EnableResponseBuffering bool `json:"enableResponseBuffering,omitempty"`
		} `json:"globalConfiguration,omitempty"`
		HTTPListeners [] HTTPListener `json:"httpListeners,omitempty"`
		Listeners []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				FrontendIPConfiguration struct {
					ID string `json:"id,omitempty"`
				} `json:"frontendIPConfiguration,omitempty"`
				FrontendPort struct {
					ID string `json:"id,omitempty"`
				} `json:"frontendPort,omitempty"`
				Protocol       string `json:"protocol,omitempty"`
				SslCertificate *struct {
					ID string `json:"id,omitempty"`
				} `json:"sslCertificate,omitempty"`
				SslProfile *struct {
					ID string `json:"id,omitempty"`
				} `json:"sslProfile,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"listeners,omitempty"`
		LoadDistributionPolicies []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				LoadDistributionAlgorithm string `json:"loadDistributionAlgorithm,omitempty"`
				LoadDistributionTargets   []struct {
					ID         string `json:"id,omitempty"`
					Name       string `json:"name,omitempty"`
					Properties struct {
						BackendAddressPool struct {
							ID string `json:"id,omitempty"`
						} `json:"backendAddressPool,omitempty"`
						WeightPerServer string `json:"weightPerServer,omitempty"`
					} `json:"properties,omitempty"`
				} `json:"loadDistributionTargets,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"loadDistributionPolicies,omitempty"`
		PrivateLinkConfigurations []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				IPConfigurations []struct {
					ID         string `json:"id,omitempty"`
					Name       string `json:"name,omitempty"`
					Properties struct {
						Primary                   bool   `json:"primary,omitempty"`
						PrivateIPAddress          string `json:"privateIPAddress,omitempty"`
						PrivateIPAllocationMethod string `json:"privateIPAllocationMethod,omitempty"`
						Subnet                    struct {
							ID string `json:"id,omitempty"`
						} `json:"subnet,omitempty"`
					} `json:"properties,omitempty"`
				} `json:"ipConfigurations,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"privateLinkConfigurations,omitempty"`
		Probes []Probe_json `json:"probes,omitempty"`
		RedirectConfigurations []RedirectConfiguration `json:"redirectConfigurations,omitempty"`
		RequestRoutingRules []RequestRoutingRule `json:"requestRoutingRules,omitempty"`
		RewriteRuleSets []struct {
			Name       string `json:"name,omitempty"`
			ID         string `json:"id,omitempty"`
			Etag       string `json:"etag,omitempty"`
			Properties struct {
				ProvisioningState string `json:"provisioningState,omitempty"`
				RewriteRules      []struct {
					RuleSequence int           `json:"ruleSequence,omitempty"`
					Conditions   []interface{} `json:"conditions,omitempty"`
					Name         string        `json:"name,omitempty"`
					ActionSet    struct {
						RequestHeaderConfigurations *[]struct {
							HeaderName  string `json:"headerName,omitempty"`
							HeaderValue string `json:"headerValue,omitempty"`
						} `json:"requestHeaderConfigurations,omitempty"`
						ResponseHeaderConfigurations *[]struct {
							HeaderName  string `json:"headerName,omitempty"`
							HeaderValue string `json:"headerValue,omitempty"`
						} `json:"responseHeaderConfigurations,omitempty"`
						URLConfiguration *struct {
							ModifiedPath        string `json:"modifiedPath,omitempty"`
							ModifiedQueryString string `json:"modifiedQueryString,omitempty"`
							Reroute             bool   `json:"reroute,omitempty"`
						} `json:"urlConfiguration,omitempty"`
					} `json:"actionSet,omitempty"`
				} `json:"rewriteRules,omitempty"`
				RequestRoutingRules *[]struct {
					ID string `json:"id,omitempty"`
				} `json:"requestRoutingRules,omitempty"`
			} `json:"properties,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"rewriteRuleSets,omitempty"`
		RoutingRules []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				BackendAddressPool struct {
					ID string `json:"id,omitempty"`
				} `json:"backendAddressPool,omitempty"`
				BackendSettings struct {
					ID string `json:"id,omitempty"`
				} `json:"backendSettings,omitempty"`
				Listener struct {
					ID string `json:"id,omitempty"`
				} `json:"listener,omitempty"`
				RuleType string `json:"ruleType,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"routingRules,omitempty"`
		Sku struct {
			Name     string `json:"name,omitempty"`
			Tier     string `json:"tier,omitempty"`
			Capacity int    `json:"capacity,omitempty"`
		} `json:"sku,omitempty"`
		SslCertificates []SslCertificate `json:"sslCertificates,omitempty"`
		SslPolicy struct {
			CipherSuites         []string `json:"cipherSuites,omitempty"`
			DisabledSslProtocols []string `json:"disabledSslProtocols,omitempty"`
			MinProtocolVersion   string   `json:"minProtocolVersion,omitempty"`
			PolicyName           string   `json:"policyName,omitempty"`
			PolicyType           string   `json:"policyType,omitempty"`
		} `json:"sslPolicy,omitempty"`
		SslProfiles []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				ClientAuthConfiguration struct {
					VerifyClientCertIssuerDN bool `json:"verifyClientCertIssuerDN,omitempty"`
				} `json:"clientAuthConfiguration,omitempty"`
				SslPolicy struct {
					CipherSuites         []string `json:"cipherSuites,omitempty"`
					DisabledSslProtocols []string `json:"disabledSslProtocols,omitempty"`
					MinProtocolVersion   string   `json:"minProtocolVersion,omitempty"`
					PolicyName           string   `json:"policyName,omitempty"`
					PolicyType           string   `json:"policyType,omitempty"`
				} `json:"sslPolicy,omitempty"`
				TrustedClientCertificates []struct {
					ID string `json:"id,omitempty"`
				} `json:"trustedClientCertificates,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"sslProfiles,omitempty"`
		TrustedClientCertificates []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				Data string `json:"data,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"trustedClientCertificates,omitempty"`
		TrustedRootCertificates []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				Data             string `json:"data,omitempty"`
				KeyVaultSecretID string `json:"keyVaultSecretId,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"trustedRootCertificates,omitempty"`
		URLPathMaps []struct {
			ID         string `json:"id,omitempty"`
			Name       string `json:"name,omitempty"`
			Properties struct {
				DefaultBackendAddressPool struct {
					ID string `json:"id,omitempty"`
				} `json:"defaultBackendAddressPool,omitempty"`
				DefaultBackendHTTPSettings struct {
					ID string `json:"id,omitempty"`
				} `json:"defaultBackendHttpSettings,omitempty"`
				DefaultLoadDistributionPolicy struct {
					ID string `json:"id,omitempty"`
				} `json:"defaultLoadDistributionPolicy,omitempty"`
				DefaultRedirectConfiguration struct {
					ID string `json:"id,omitempty"`
				} `json:"defaultRedirectConfiguration,omitempty"`
				DefaultRewriteRuleSet struct {
					ID string `json:"id,omitempty"`
				} `json:"defaultRewriteRuleSet,omitempty"`
				PathRules []struct {
					ID         string `json:"id,omitempty"`
					Name       string `json:"name,omitempty"`
					Properties struct {
						BackendAddressPool struct {
							ID string `json:"id,omitempty"`
						} `json:"backendAddressPool,omitempty"`
						BackendHTTPSettings struct {
							ID string `json:"id,omitempty"`
						} `json:"backendHttpSettings,omitempty"`
						FirewallPolicy *struct {
							ID string `json:"id,omitempty"`
						} `json:"firewallPolicy,omitempty"`
						LoadDistributionPolicy struct {
							ID string `json:"id,omitempty"`
						} `json:"loadDistributionPolicy,omitempty"`
						Paths                 []string `json:"paths,omitempty"`
						RedirectConfiguration struct {
							ID string `json:"id,omitempty"`
						} `json:"redirectConfiguration,omitempty"`
						RewriteRuleSet struct {
							ID string `json:"id,omitempty"`
						} `json:"rewriteRuleSet,omitempty"`
					} `json:"properties,omitempty"`
				} `json:"pathRules,omitempty"`
			} `json:"properties,omitempty"`
		} `json:"urlPathMaps,omitempty"`
		WebApplicationFirewallConfiguration *struct {
			Enabled            bool   `json:"enabled,omitempty"`
			MaxRequestBodySize int    `json:"maxRequestBodySize,omitempty"`
			FirewallMode       string `json:"firewallMode,omitempty"`
			RuleSetType        string `json:"ruleSetType,omitempty"`
			RuleSetVersion     string `json:"ruleSetVersion,omitempty"`
			DisabledRuleGroups []struct {
				RuleGroupName string   `json:"ruleGroupName,omitempty"`
				Rules         []string `json:"rules,omitempty"`
			} `json:"disabledRuleGroups,omitempty"`
			Exclusions *[]struct {
				MatchVariable         string `json:"matchVariable,omitempty"`
				Selector              string `json:"selector,omitempty"`
				SelectorMatchOperator string `json:"selectorMatchOperator,omitempty"`
			} `json:"exclusions,omitempty"`
			RequestBodyCheck       bool `json:"requestBodyCheck,omitempty"`
			MaxRequestBodySizeInKb int  `json:"maxRequestBodySizeInKb,omitempty"`
			FileUploadLimitInMb    int  `json:"fileUploadLimitInMb,omitempty"`
		} `json:"webApplicationFirewallConfiguration,omitempty"`
		Zones                      []string      `json:"zones,omitempty"`
		ProvisioningState          string        `json:"provisioningState,omitempty"`
		ResourceGUID               string        `json:"resourceGuid,omitempty"`
		OperationalState           string        `json:"operationalState,omitempty"`
		PrivateEndpointConnections []interface{} `json:"privateEndpointConnections,omitempty"`
	} `json:"properties,omitempty"`
}
