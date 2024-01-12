terraform {
  required_providers {
    vcf = {
      source  = "vmware/vcf"
    }
  }
}

data "vcf_credentials" "sddc_creds" {
  resource_type = "VCENTER"
  account_type = "USER"
}

resource "vcf_credentials_rotate" "vc_0_rotate" {
  resource_name = data.vcf_credentials.sddc_creds.credentials[0].resource[0].name
  resource_type = data.vcf_credentials.sddc_creds.credentials[0].resource[0].type
  credentials {
    credential_type = data.vcf_credentials.sddc_creds.credentials[0].credential_type
    user_name = data.vcf_credentials.sddc_creds.credentials[0].user_name
  }
}