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
    ogranization_unit = "CIBG"
  }
}

