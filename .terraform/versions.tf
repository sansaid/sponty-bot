terraform {
  cloud {
    organization = local.tf_org

    workspaces {
      name = [local.app_name]
    }
  }
  
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.10.0"
    }
  }
}