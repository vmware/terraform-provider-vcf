variable "vcenter_username" {
  description = "Username used to authenticate against the vCenter Server"
}

variable "vcenter_password" {
  description = "Password used to authenticate against the vCenter Server"
}

variable "vcenter_server" {
  description = "FQDN or IP Address of the vCenter Server"
}

variable "datacenter_name" {
  description = "The name of the datacenter"
}

variable "vsan_datastore_name" {
  description = "The name of the vsan datastore"
}

variable "content_library_name" {
  description = "The name of the content library"
}

variable "content_library_url" {
  description = "The subscription URL for the content library"
}