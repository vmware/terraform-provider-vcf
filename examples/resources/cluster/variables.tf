variable "sddc_manager_username" {
  description = "Username used to authenticate against an SDDC Manager instance"
  default = ""
}

variable "sddc_manager_password" {
  description = "Password used to authenticate against an SDDC Manager instance"
  default = ""
}

variable "sddc_manager_host" {
  description = "Fully qualified domain name of an SDDC Manager instance"
  default = ""
}

variable "esx_host1_fqdn" {
  description = "Fully qualified domain name of ESXi host 1"
  default = ""
}

variable "esx_host1_pass" {
  description = "Password to authenticate to the ESXi host 1"
  default = ""
}

variable "esx_host2_fqdn" {
  description = "Fully qualified domain name of ESXi host 2"
  default = ""
}

variable "esx_host2_pass" {
  description = "Password to authenticate to the ESXi host 2"
  default = ""
}

variable "esx_host3_fqdn" {
  description = "Fully qualified domain name of ESXi host 3"
  default = ""
}

variable "esx_host3_pass" {
  description = "Password to authenticate to the ESXi host 3"
  default = ""
}

variable "domain_id" {
  description = "Id of the domain in which the Cluster is to be created"
  default = ""
}

variable "esx_license_key" {
  description = "License key for an ESXi host in the free pool. This is required except in cases where the " +
  "ESXi host has already been licensed outside of the VMware Cloud Foundation system"
  default = ""
}

variable "vsan_license_key" {
  description = "vSAN license key to be used"
  default = ""
}
