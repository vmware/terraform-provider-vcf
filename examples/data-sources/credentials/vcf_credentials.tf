terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

data "vcf_credentials" "creds" {

}

data "vcf_credentials" "creds_vc" {
  resource_type = "VCENTER"
}