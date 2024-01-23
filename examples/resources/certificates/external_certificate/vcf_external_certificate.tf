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

resource "vcf_csr" "vcenter_csr" {
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

resource "vcf_external_certificate" "vcenter_cert" {
  csr_id               = vcf_csr.vcenter_csr.id
  resource_certificate = var.new_certificate
  ca_certificate       = var.ca_certificate
}