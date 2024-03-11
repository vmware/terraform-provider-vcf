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

variable "cluster_id" {
  description = "Id of the cluster in the vCenter server (e.g. domain-c21)"
  default     = ""
}

variable "domain_id" {
  description = "Id of the domain in which the cluster is to be created"
  default     = ""
}
