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

variable "cluster_id" {
  description = "The identifier (in SDDC manager) of the compute cluster"
}

variable "sddc_manager_username" {
  description = "Username used to authenticate against an SDDC Manager instance"
}

variable "sddc_manager_password" {
  description = "Password used to authenticate against an SDDC Manager instance"
}

variable "sddc_manager_host" {
  description = "FQDN or IP Address of an SDDC Manager instance"
}

variable "edge_cluster_name" {
  description = "The name of the edge cluster"
}

variable "edge_cluster_root_pass" {
  description = "The root user password for the edge cluster"
}

variable "edge_cluster_admin_pass" {
  description = "The administrator password for the edge cluster"
}

variable "edge_cluster_audit_pass" {
  description = "The audit user password for the edge cluster"
}

variable "edge_node1_cidr" {
  description = "The IP address of the first edge node (in CIDR format, e.g. 10.0.0.12/24)"
}

variable "edge_node1_root_pass" {
  description = "The root user password for the first edge node"
}

variable "edge_node1_admin_pass" {
  description = "The administrator password for the first edge node"
}

variable "edge_node1_audit_pass" {
  description = "The audit user password for the first edge node"
}

variable "edge_node2_cidr" {
  description = "The IP address of the second edge node (in CIDR format, e.g. 10.0.0.12/24)"
}

variable "edge_node2_root_pass" {
  description = "The root user password for the second edge node"
}

variable "edge_node2_admin_pass" {
  description = "The administrator password for the second edge node"
}

variable "edge_node2_audit_pass" {
  description = "The audit user password for the second edge node"
}

variable "bgp_peer_password" {
  description = "The password for the bgp peers"
}