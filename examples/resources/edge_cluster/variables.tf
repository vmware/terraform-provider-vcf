variable "sddc_manager_username" {
  description = "Username used to authenticate against an SDDC Manager instance."
  default     = ""
}

variable "sddc_manager_password" {
  description = "Password used to authenticate against an SDDC Manager instance."
  default     = ""
}

variable "sddc_manager_host" {
  description = "Fully qualified domain name of an SDDC Manager instance."
  default     = ""
}

variable "cluster_name" {
  description = "The display name of the edge cluster."
  default     = ""
}

variable "cluster_root_pass" {
  description = "The root user password for the edge cluster."
  default     = ""
}

variable "cluster_admin_pass" {
  description = "The admin user password for the edge cluster."
  default     = ""
}

variable "cluster_audit_pass" {
  description = "The audit user password for the edge cluster."
  default     = ""
}

variable "compute_cluster_id" {
  description = "The identifier of the compute cluster where the edge nodes will be deployed."
  default     = ""
}

variable "edge_node_1_name" {
  description = "The display name of the first edge node."
  default     = ""
}

variable "edge_node_1_root_pass" {
  description = "The root user password for the first edge node."
  default     = ""
}

variable "edge_node_1_admin_pass" {
  description = "The admin user password for the first edge node."
  default     = ""
}

variable "edge_node_1_audit_pass" {
  description = "The audit user password for the first edge node."
  default     = ""
}

variable "edge_node_2_name" {
  description = "The display name of the second edge node."
  default     = ""
}

variable "edge_node_2_root_pass" {
  description = "The root user password for the second edge node."
  default     = ""
}

variable "edge_node_2_admin_pass" {
  description = "The admin user password for the second edge node."
  default     = ""
}

variable "edge_node_2_audit_pass" {
  description = "The audit user password for the second edge node."
  default     = ""
}
