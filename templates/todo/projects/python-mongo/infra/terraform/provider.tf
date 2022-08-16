#Set the terraform required version, and Configure the Azure Provider.Use remote storage

# Configure the Azure Provider
terraform {
  required_version = "= v1.1.7"
  required_providers {
    azurerm = {
      version = "~>3.18.0"
      source  = "hashicorp/azurerm"
    }
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "1.2.15"
    }
  }
}

provider "azurerm" {
  features {}
}

# Make client_id, tenant_id, subscription_id and object_id variables
data "azurerm_client_config" "current" {}