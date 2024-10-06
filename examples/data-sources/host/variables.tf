variable "sddc_manager_host" {
  type        = string
  description = "The fully qualified domain name of the SDDC Manager instance."
}

variable "sddc_manager_username" {
  type        = string
  description = "The username to authenticate to the SDDC Manager instance."
  sensitive   = true
}

variable "sddc_manager_password" {
  type        = string
  description = "The password to authenticate to the SDDC Manager instance."
  sensitive   = true
}

variable "host_fqdn" {
  type        = string
  description = "The fully qualified domain name of the ESXi host."
  default     = "sfo-w01-esx01.sfo.rainpole.io"
}
