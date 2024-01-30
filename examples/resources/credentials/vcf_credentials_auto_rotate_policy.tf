terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

data "vcf_credentials" "sddc_creds" {
  resource_type = "VCENTER"
}

resource "vcf_credentials_auto_rotate_policy" "vc_0_autorotate" {
  resource_id          = data.vcf_credentials.sddc_creds.credentials[0].resource[0].id
  resource_type        = data.vcf_credentials.sddc_creds.credentials[0].resource[0].type
  user_name            = data.vcf_credentials.sddc_creds.credentials[0].user_name
  enable_auto_rotation = true
  auto_rotate_days     = 7
}