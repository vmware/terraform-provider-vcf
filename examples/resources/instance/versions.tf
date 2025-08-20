# Terraform and Provider Version Constraints
terraform {
  required_version = ">= 1.0"

  required_providers {
    vcf = {
      source  = "vmware/vcf"
      version = "0.17.0"
    }
  }
}