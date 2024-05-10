variable "vcenter_username" {
  description = "Username used to authenticate against the vCenter Server"
}

variable "vcenter_password" {
  description = "Password used to authenticate against the vCenter Server"
}

variable "vcenter_server" {
  description = "FQDN or IP Address of the vCenter Server for the workload domain"
}

variable "datacenter_name" {
  description = "The name of the datacenter"
}

variable "cluster_name" {
  description = "The name of the compute cluster"
}

variable "storage_policy_name" {
  description = "The name of the storage policy"
}

variable "management_network_name" {
  description = "The name of the management network. This should be a distributed portgroup on the DVS which the edge nodes are are connected to"
}

variable "contenty_library_name" {
  description = "The name of the subscribed content library"
}

variable "dvs_name" {
  description = "The name of the distributed switch"
}

variable "edge_cluster" {
  description = "The identifier of the edge cluster"
}

variable "host1_fqdn" {
  description = "The fully qualified domain name of one of the hosts in the cluster"
}

variable "namespace_name" {
  description = "The name for the new vSphere Namespace"
}

variable "virtual_machine_class_name" {
  description = "The name of the new Virtual Machine Class"
}