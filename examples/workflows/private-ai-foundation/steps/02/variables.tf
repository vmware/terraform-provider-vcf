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

variable "source_cluster_name" {
  description = "The name of the compute cluster"
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

variable "vcenter_root_password" {
  description = "root password for the vCenter Server Appliance (8-20 characters)"
}

variable "nsx_manager_admin_password" {
  description = "NSX Manager admin user password"
}

variable "esx_host1_fqdn" {
  description = "Fully qualified domain name of ESXi host 1"
}

variable "esx_host1_pass" {
  description = "Password to authenticate to the ESXi host 1"
}

variable "esx_host2_fqdn" {
  description = "Fully qualified domain name of ESXi host 2"
}

variable "esx_host2_pass" {
  description = "Password to authenticate to the ESXi host 2"
}

variable "esx_host3_fqdn" {
  description = "Fully qualified domain name of ESXi host 3"
}

variable "esx_host3_pass" {
  description = "Password to authenticate to the ESXi host 3"
}

variable "nsx_license_key" {
  description = "NSX license to be used"
}

variable "esx_license_key" {
  description = "License key for an ESXi host in the free pool."
}

variable "vsan_license_key" {
  description = "vSAN license key to be used"
}

variable "workload_domain_name" {
  description = "The name for the new workload domain"
}

variable "management_domain_id" {
  description = "The identifier off the management domain"
}

variable "workload_vcenter_address" {
  description = "The IP address for the vCenter Server in the new workload domain"
}

variable "workload_vcenter_fqdn" {
  description = "The fully qualified domain name for the vCenter Server in the new workload domain"
}

variable "nsx_manager_node1_address" {
  description = "The IP address for the first NSX Manager node"
}

variable "nsx_manager_node2_address" {
  description = "The IP address for the second NSX Manager node"
}

variable "nsx_manager_node3_address" {
  description = "The IP address for the third NSX Manager node"
}

variable "nsx_manager_vip_address" {
  description = "The virtual IP for the NSX Manager"
}