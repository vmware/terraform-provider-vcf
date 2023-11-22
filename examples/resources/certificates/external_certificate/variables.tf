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

variable "vcf_domain_id" {
  description = "Id of the domain for whose resources CSRs are to be generated"
  default = ""
}

variable "new_certificate" {
  description = "PEM encoded certificate for a resource within a Domain that is to replace the old one. Can be issued by an external to VCF CA."
  default = ""
}

variable "ca_certificate" {
  description = "PEM encoded certificate of the CA that issued the above resource certificate"
  default = ""
}