terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

data "vcf_domain" "sfo-m01" {
  domain_id = var.domain_id
}

data "vcf_domain" "sfo-w01" {
  name = var.domain_name
}