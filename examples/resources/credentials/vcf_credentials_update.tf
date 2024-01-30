terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

resource "vcf_credentials_update" "vc_0_update" {
  resource_name = var.esxi_host_1
  resource_type = "ESXI"
  credentials {
    credential_type = "SSH"
    user_name       = "root"
    password        = var.esxi_host_1_pass
  }
}
