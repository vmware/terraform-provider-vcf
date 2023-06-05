variable "sddc_manager_username" {
  description = "Username used to authenticate against an SDDC Manager instance"
  default = ""
}

variable "sddc_manager_password" {
  description = "Password used to authenticate against an SDDC Manager instance"
  default = ""
}

variable "sddc_manager_host" {
  description = "FQDN of an SDDC Manager instance"
  default = ""
}

variable "host_fqdn" {
  description = "FQDN of an ESXi host that is to be commissioned"
  default = ""
}

variable "host_ssh_user" {
  description = "SSH user in ESXi host that is to be commissioned"
  default = ""
}

variable "host_ssh_pass" {
  description = "SSH pass in ESXi host that is to be commissioned"
  default = ""
}