terraform {
  required_providers {
    azurermagw = {
      source  = "citeo.com/edu/azurermagw"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.5.0"
    }
  }
}
provider "azurermagw" {
}
provider "azurerm" {
  features {}
}
data "azurerm_application_gateway" "appgw"{
  name = "app-gateway"
  resource_group_name = "shared-app-gateway"
}
locals {
  probe_name                  = "aaa-mahmoud-probe-1"
  ssl_certificate_name        = "aaa-mahmoud-certificat-1"
  https_listener_name         = "aaa-mahmoud-httplistener-1-https"
  secret_id                   = "https://kv-mahmoud.vault.azure.net/secrets/default-citeo-ddc-cert/6ae31f879c2b4e05af69d0ccd7d7aba6"
  frontend_ip_configuration   = "app-gateway-fe-ip-config"
  frontend_port_http          = "app-gateway-fe-port"
  frontend_port_https         = "app-gateway-fe-port-443"
  backend_address_pool_name   = "aaa-mahmoud-backendAddressPool-1"
  backend_http_settings_name  = "aaa-mahmoud-backendHttpSettings-1"
  http_listener_name          = "aaa-mahmoud-httplistener-1-https"
  redirect_configuration_name = "aaa-mahmoud-redirectconfiguration-1"

}
resource "azurermagw_webappBinding" "citeo-binding" {
  name = "mahmoud-backendAddressPool-resource-name"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name

  request_routing_rule = {
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.backend_http_settings_name
    http_listener_name         = local.http_listener_name
    name                       = "aaa-mahmoud-requestroutingrule-1"
    rule_type                  = "Basic"
 //   redirect_configuration_name= local.redirect_configuration_name
 //   rewrite_rule_set_name      = ""
  //  url_path_map_name		       = ""
  }

  redirect_configuration = {
    name                 = "aaa-mahmoud-redirectconfiguration-1"
    redirect_type        = "Permanent"
    target_listener_name = local.https_listener_name #"default-citeo-plus-listener-443"
//    target_url           = "https://www.mahmoud.com"
  }

  ssl_certificate = {
    name                = local.ssl_certificate_name #azurerm_key_vault_certificate.default_citeo_ddc_cert.name
    key_vault_secret_id = local.secret_id #azurerm_key_vault_certificate.default_citeo_ddc_cert.secret_id
//    data                = "azerazerazerazerazerazerazerazer"
//    password            = "mahmoudpass"
  }

  http_listener = {
    frontend_ip_configuration_name = local.frontend_ip_configuration #data.azurerm_application_gateway.appgw.frontend_ip_configuration.name #"appGwPublicFrontendIp"
    frontend_port_name             = local.frontend_port_http #data.azurerm_application_gateway.appgw.frontend_port.name # "port_80"
    host_name                      = "www.mahmoud.com"
    name                           = "aaa-mahmoud-httplistener-1-http"
   protocol                       = "Http"
  }

  https_listener = {
    frontend_ip_configuration_name =local.frontend_ip_configuration # data.azurerm_application_gateway.appgw.frontend_ip_configuration.name # "appGwPublicFrontendIp"
    frontend_port_name             = local.frontend_port_https #data.azurerm_application_gateway.appgw.frontend_port.name # "port_443"
    host_name                      = "www.mahmoud.com"
//    host_names                      = ["www.mahmoud.com","ddd.eee.tn"]
    name                           = "aaa-mahmoud-httplistener-1-https"
    protocol                       = "Https"
    require_sni                    = true
    ssl_certificate_name           = local.ssl_certificate_name #azurerm_key_vault_certificate.default_citeo_gdtr_iframe_cert.name
  }

  backend_address_pool = {
    name = "aaa-mahmoud-backendAddressPool-1"
    fqdns = ["fqdn.mahmoud"]
//    ip_addresses=["10.2.3.3"]
  }

  backend_http_settings = {
    name                                = "aaa-mahmoud-backendHttpSettings-1"
    affinity_cookie_name                = "ApplicationGatewayAffinity"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 443
    probe_name                          = local.probe_name #"mahmoud-probe-2"
    protocol                            = "Https"
    request_timeout                     = 667
  }

  probe = {
    name                                      = local.probe_name #"mahmoud-probe-2"
    interval                                  = 30
    protocol                                  = "Https"
    path                                      = "/"
    timeout                                   = 30
    unhealthy_threshold                       = 3
    pick_host_name_from_backend_http_settings = true
    minimum_servers         = 1
    match = {
      body        = ""
      status_code = ["200-399","400","500"]
    }
  }
}
