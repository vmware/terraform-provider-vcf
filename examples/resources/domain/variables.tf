variable "sddc_manager_username" {
  description = "Username used to authenticate against an SDDC Manager instance"
  default     = ""
}

variable "sddc_manager_password" {
  description = "Password used to authenticate against an SDDC Manager instance"
  default     = ""
}

variable "sddc_manager_host" {
  description = "Fully qualified domain name of an SDDC Manager instance"
  default     = ""
}

variable "vcenter_root_password" {
  description = "root password for the vCenter Server Appliance (8-20 characters)"
  default     = ""
}

variable "nsx_manager_admin_password" {
  description = "NSX Manager admin user password"
  default     = ""
}

variable "esx_host1_fqdn" {
  description = "Fully qualified domain name of ESXi host 1"
  default     = ""
}

variable "esx_host1_pass" {
  description = "Password to authenticate to the ESXi host 1"
  default     = ""
}

variable "esx_host2_fqdn" {
  description = "Fully qualified domain name of ESXi host 2"
  default     = ""
}

variable "esx_host2_pass" {
  description = "Password to authenticate to the ESXi host 2"
  default     = ""
}

variable "esx_host3_fqdn" {
  description = "Fully qualified domain name of ESXi host 3"
  default     = ""
}

variable "esx_host3_pass" {
  description = "Password to authenticate to the ESXi host 3"
  default     = ""
}

variable "nsx_license_key" {
  description = "NSX license to be used"
  default     = ""
}

variable "esx_license_key" {
  description = "License key for an ESXi host in the free pool. This is required except in cases where the ESXi host has already been licensed outside of the VMware Cloud Foundation system"
  default     = ""
}

variable "vsan_license_key" {
  description = "vSAN license key to be used"
  default     = ""
}