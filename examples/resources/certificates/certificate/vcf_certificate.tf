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

resource "vcf_certificate_authority" "ca" {
  open_ssl {
    common_name       = "test.openssl.eng.vmware.com"
    country           = "BG"
    state             = "Sofia-grad"
    locality          = "Sofia"
    organization      = "VMware"
    ogranization_unit = "VCF"
  }
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

resource "vcf_certificate" "vcenter_cert" {
  csr_id = vcf_csr.vcenter_csr.id
  ca_id  = vcf_certificate_authority.ca.id
}