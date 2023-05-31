variable "sddc_manager_username" {
  description = "username used to authenticate against an SDDC Manager instance"
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

variable "sso_domain" {
  description = "The SSO domain in which an SSO user is to be created"
  default = "vrack.vsphere.local"
}

variable "sso_username" {
  description = "Username of an SSO user to be created"
  default = ""
}

variable "sso_service_username" {
  description = "Username of an SSO service user to be created"
  default = ""
}