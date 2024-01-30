variable "cloud_builder_username" {
  description = "Username to authenticate to CloudBuilder"
  default     = ""
}

variable "cloud_builder_password" {
  description = "Password to authenticate to CloudBuilder"
  default     = ""
}

variable "cloud_builder_host" {
  description = "Fully qualified domain name or IP address of the CloudBuilder"
  default     = ""
}

variable "sddc_manager_root_user_password" {
  description = "Root user password for the SDDC Manager VM. Password needs to be a strong password with at least one alphabet and one special character and at least 8 characters in length"
  default     = ""
}

variable "sddc_manager_secondary_user_password" {
  description = "Second user (vcf) password for the SDDC Manager VM.  Password needs to be a strong password with at least one alphabet and one special character and at least 8 characters in length."
  default     = ""
}

variable "vcenter_root_password" {
  description = "root password for the vCenter Server Appliance (8-20 characters)"
  default     = ""
}

variable "nsx_manager_admin_password" {
  description = "NSX admin password. The password must be at least 12 characters long. Must contain at-least 1 uppercase, 1 lowercase, 1 special character and 1 digit. In addition, a character cannot be repeated 3 or more times consecutively."
  default     = ""
}

variable "nsx_manager_audit_password" {
  description = "NSX audit password. The password must be at least 12 characters long. Must contain at-least 1 uppercase, 1 lowercase, 1 special character and 1 digit. In addition, a character cannot be repeated 3 or more times consecutively."
  default     = ""
}

variable "nsx_manager_root_password" {
  description = " NSX Manager root password. Password should have 1) At least eight characters, 2) At least one lower-case letter, 3) At least one upper-case letter 4) At least one digit 5) At least one special character, 6) At least five different characters , 7) No dictionary words, 6) No palindromes"
  default     = ""
}

variable "esx_host1_pass" {
  description = "Password to authenticate to the ESXi host 1"
  default     = ""
}

variable "esx_host2_pass" {
  description = "Password to authenticate to the ESXi host 2"
  default     = ""
}

variable "esx_host3_pass" {
  description = "Password to authenticate to the ESXi host 3"
  default     = ""
}

variable "esx_host4_pass" {
  description = "Password to authenticate to the ESXi host 4"
  default     = ""
}

variable "nsx_license_key" {
  description = "NSX license to be used"
  default     = ""
}

variable "vcenter_license_key" {
  description = "vCenter license to be used"
  default     = ""
}

variable "vsan_license_key" {
  description = "vSAN license key to be used"
  default     = ""
}