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

variable "vcf_cluster_id" {
  description = "Id of the cluster that is to be used as a data source. Note: management domain default cluster ID can be used to refer to some of it's attributes"
  default = ""
}