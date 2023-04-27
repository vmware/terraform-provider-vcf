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

variable "esxi_1_user" {
  description = "SSH user of ESXi1 server"
  default = ""
}

variable "esxi_1_pass" {
  description = "SSH password of ESXi1 server"
  default = ""
}

variable "esxi_2_user" {
  description = "SSH user of ESXi2 server"
  default = ""
}

variable "esxi_2_pass" {
  description = "SSH password of ESXi2 server"
  default = ""
}

variable "esxi_3_user" {
  description = "SSH user of ESXi3 server"
  default = ""
}

variable "esxi_3_pass" {
  description = "SSH password of ESXi3 server"
  default = ""
}