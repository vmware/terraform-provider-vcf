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
  description = "The name of the datacenter where the new cluster will be created"
}

variable "cluster_name" {
  description = "The name of the compute cluster"
}

variable "depot_location" {
  description = "The URL where the contents for the offline software depot are hosted"
}