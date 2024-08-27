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

data "vcf_host" "example" {
  fqdn = var.host_fqdn
}

output "host_id" {
  value = data.vcf_host.example.id
}
