terraform {
  required_providers {
    hashicups = {
      source  = "citeo.com/edu/azurerm_agw"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.5.0"
    }
  }
}
provider "hashicups" {
}
provider "azurerm" {
  features {}
}
data "azurerm_application_gateway" "appgw"{
  name = "app-gateway"
  resource_group_name = "shared-app-gateway"
} 
 
resource "azurerm_agw_webappBinding" "citeo-binding" {
  name = "mahmoud-backendAddressPool-resource-name"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name"
    fqdns = ["fqdn.mahmoud"]
    ip_addresses=["10.2.3.3"]
  }
}
/*
resource "azurerm_agw_webappBinding" "citeo-binding4" {
  name = "mahmoud-backendAddressPool-resource-name4"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name4"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurerm_agw_webappBinding" "citeo-binding3" {
  name = "mahmoud-backendAddressPool-resource-name3"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name3"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurerm_agw_webappBinding" "citeo-binding2" {
  name = "mahmoud-backendAddressPool-resource-name2"
  agw_name              = data.azurerm_application_gateway.appgw.name
  agw_rg                = data.azurerm_application_gateway.appgw.resource_group_name
  backend_address_pool = {
    name = "mahmoud-backendAddressPool-name2"
    fqdns = ["fqdn.mahmoud.net"]
    ip_addresses=["100.0.0.100"]
  }
}
resource "azurerm_agw_webappBinding" "citeo-binding1" {
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
  value = azurerm_agw_webappBinding.citeo-binding
}
*/
