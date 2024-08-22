variable "sddc_manager_host" {
  description = "The fully qualified domain name of the SDDC Manager instance."
}

variable "sddc_manager_username" {
  description = "The username to authenticate to the SDDC Manager instance."
  sensitive   = true
}

variable "sddc_manager_password" {
  description = "The password to authenticate to the SDDC Manager instance."
  sensitive   = true
}

variable "domain_id" {
  description = "The ID of the workload domain."
}

variable "domain_name" {
  description = "The name of the workload domain."
  name        = "sfo-w01"
}
