terraform {
  required_providers {
    vcf = {
      source  = "vmware/vcf"
    }
  }
}

provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

resource "vcf_user" "user1" {
  name      = var.sso_username
  domain    = var.sso_domain
  type      = "USER"
  role_name = "ADMIN"
}

# Service users have api_key output associated with them
resource "vcf_user" "service_user1" {
  name      = var.sso_service_username
  domain    = var.sso_domain
  type      = "SERVICE"
  role_name = "VIEWER"
}