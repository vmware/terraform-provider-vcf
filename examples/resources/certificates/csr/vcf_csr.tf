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

resource "vcf_csr" "csr1" {
  domain_id         = var.vcf_domain_id
  country           = "BG"
  email             = "admin@vmware.com"
  key_size          = "3072"
  locality          = "Sofia"
  state             = "Sofia-grad"
  organization      = "VMware Inc."
  organization_unit = "VCF"
  resource          = "VCENTER"
}