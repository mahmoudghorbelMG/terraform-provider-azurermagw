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
 
resource "azurermagw_webappBinding" "citeo-binding" {
  name = "mahmoud-backendAddressPool-resource-name"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name"
    fqdns = ["fqdn.mahmoud"]
    ip_addresses=["10.2.3.3"]
  }
  backend_http_settings = {
    name                                = "mahmoud-backendHttpSettings-3"
    affinity_cookie_name                = "ApplicationGatewayAffinity"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 443
    probe_name                          = "mahmoud-probe-3"
    protocol                            = "Https"
    request_timeout                     = 667
  }
}
/*
resource "azurermagw_webappBinding" "citeo-binding4" {
  name = "mahmoud-backendAddressPool-resource-name4"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name4"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurermagw_webappBinding" "citeo-binding3" {
  name = "mahmoud-backendAddressPool-resource-name3"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name3"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurermagw_webappBinding" "citeo-binding2" {
  name = "mahmoud-backendAddressPool-resource-name2"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name2"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurermagw_webappBinding" "citeo-binding1" {
  name = "mahmoud-backendAddressPool-resource-name1"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name1"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }

}*/
/*
output "citeo-binding_out" {
  value = azurermagw_webappBinding.citeo-binding
}
*/
