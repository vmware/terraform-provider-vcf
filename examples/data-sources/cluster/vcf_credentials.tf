terraform {
  required_providers {
    vcf = {
      source  = "vmware/vcf"
    }
  }
}

data "vcf_credentials" "creds" {

}

data "vcf_credentials" "vcenter_credentials" {
  resource_type = "VCENTER"
}