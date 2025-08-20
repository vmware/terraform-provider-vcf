# VCF Provider Configuration
variable "installer_host" {
  description = "The hostname or IP address of the VCF installer"
  type        = string
}

variable "installer_username" {
  description = "Username for VCF installer authentication"
  type        = string
}

variable "installer_password" {
  description = "Password for VCF installer authentication"
  type        = string
  sensitive   = true
}

variable "allow_unverified_tls" {
  description = "Allow unverified TLS connections"
  type        = bool
  default     = true
}

variable "ceip_enabled" {
  description = "Enable VMware Customer Experience Improvement Program"
  type        = bool
  default     = false
}

variable "fips_enabled" {
  description = "Enable FIPS mode"
  type        = bool
  default     = false
}

# VCF Instance Configuration
variable "instance_id" {
  description = "Unique identifier for the VCF instance"
  type        = string
}

variable "vcf_version" {
  description = "VCF version"
  type        = string
}

variable "skip_esx_thumbprint_validation" {
  description = "Whether to skip ESX thumbprint validation"
  type        = bool
  default     = false
}

variable "network_pool_name" {
  description = "Name of the network pool for management"
  type        = string
}

# SDDC Manager Configuration
variable "sddc_manager_info" {
  description = "SDDC Manager configuration details"
  type = object({
    hostname            = string
    ssh_password        = string
    root_user_password  = string
    local_user_password = string
  })
  sensitive = true
}

# NTP Configuration
variable "ntp_servers" {
  description = "List of NTP servers"
  type        = list(string)
  default     = ["pool.ntp.org"]
}

# DNS Configuration
variable "dns_info" {
  description = "DNS configuration details"
  type = object({
    domain               = string
    primary_nameserver   = string
    secondary_nameserver = optional(string)
  })
}

# Network Configurations
variable "management_network" {
  description = "Management network configuration"
  type = object({
    subnet         = string
    vlan_id        = number
    mtu            = number
    network_type   = string
    gateway        = string
    active_uplinks = list(string)
  })
}

variable "vmotion_network" {
  description = "vMotion network configuration"
  type = object({
    subnet         = string
    vlan_id        = number
    mtu            = number
    network_type   = string
    gateway        = string
    active_uplinks = list(string)
  })
}

variable "vsan_network" {
  description = "vSAN network configuration"
  type = object({
    subnet         = string
    vlan_id        = number
    mtu            = number
    network_type   = string
    gateway        = string
    active_uplinks = list(string)
  })
}

variable "vm_management_network" {
  description = "VM Management network configuration"
  type = object({
    subnet         = string
    vlan_id        = number
    mtu            = number
    network_type   = string
    gateway        = string
    active_uplinks = list(string)
  })
}

# IP Address Range Configurations
variable "vmotion_ip_address_range" {
  description = "vMotion IP address range"
  type = object({
    ip_start = string
    ip_end   = string
  })
}

variable "vsan_ip_address_range" {
  description = "vSAN IP address range"
  type = object({
    ip_start = string
    ip_end   = string
  })
}

variable "overlay_ip_address_range" {
  description = "Overlay network IP address range"
  type = object({
    ip_start = string
    ip_end   = string
  })
}

# NSX Configuration
variable "nsx_manager_size" {
  description = "Size of NSX Manager VMs"
  type        = string
  default     = "medium"
  validation {
    condition     = contains(["medium", "large", "xlarge"], var.nsx_manager_size)
    error_message = "NSX Manager size must be medium, large, or xlarge."
  }
}

variable "nsx_manager_hostname" {
  description = "List of NSX Manager hostnames"
  type        = list(string)
}

variable "nsx_passwords" {
  description = "NSX password configuration"
  type = object({
    root_nsx_manager_password = string
    nsx_admin_password        = string
    nsx_audit_password        = string
  })
  sensitive = true
}

variable "tep_ip_address_pool" {
  description = "TEP (Tunnel Endpoint) IP address pool configuration"
  type = object({
    name              = string
    subnet            = string
    gateway           = string
    transport_vlan_id = number
  })
}

variable "nsx_vip_fqdn" {
  description = "NSX VIP fully qualified domain name"
  type        = string
}

variable "nsx_vlan_transportzone" {
  description = "NSX VLAN transport zone"
  type        = string
}

variable "nsx_overlay_transportzone" {
  description = "NSX overlay transport zone"
  type        = string
}

# vSAN Configuration
variable "vsan_config" {
  description = "vSAN configuration details"
  type = object({
    datastore_name       = string
    license_key          = optional(string)
    failures_to_tolerate = optional(number, 1)
    esa_enabled          = bool
  })
}


# DVS Configuration
variable "dvs_name" {
  description = "Name of the Distributed Virtual Switch"
  type        = string
  default     = "VCF-DVS"
}

variable "dvs_mtu" {
  description = "MTU size for the Distributed Virtual Switch"
  type        = number
  default     = 1500
}

variable "vmnic_mappings" {
  description = "List of vmnic to uplink mappings"
  type = list(object({
    vmnic  = string
    uplink = string
  }))
}

# Cluster Configuration
variable "datacenter_name" {
  description = "Name of the datacenter"
  type        = string
}

variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
}

# vCenter Configuration
variable "vcenter_config" {
  description = "vCenter Server configuration details"
  type = object({
    hostname      = string
    root_password = string
    vm_size       = optional(string, "medium")
    storage_size  = optional(string, "lstorage")
  })
  sensitive = true

  validation {
    condition     = contains(["tiny", "small", "medium", "large", "xlarge"], var.vcenter_config.vm_size)
    error_message = "vCenter VM size must be tiny, small, medium, large, or xlarge."
  }
}

# Host Configuration
variable "hosts" {
  description = "List of ESXi hosts to be added to the cluster"
  type = list(object({
    hostname = string
    credentials = object({
      username = string
      password = string
    })
  }))
  sensitive = true
}

# Operations Fleet Management Configuration
variable "operations_fleet_management" {
  description = "VCF Operations Fleet Management configuration"
  type = object({
    hostname            = string
    root_user_password  = optional(string)
    admin_user_password = optional(string)
  })
  sensitive = true
  default   = null
}

# Operations Configuration
variable "operations" {
  description = "VCF Operations configuration"
  type = object({
    node = list(object({
      hostname           = string
      root_user_password = optional(string)
      type               = string
    }))
    admin_user_password = optional(string)
    appliance_size      = optional(string, "small")
    load_balancer_fqdn  = optional(string)
  })
  sensitive = true
  default   = null

  validation {
    condition     = var.operations == null || contains(["xsmall", "small", "medium", "large", "xlarge"], var.operations.appliance_size)
    error_message = "Operations appliance size must be xsmall, small, medium, large, or xlarge."
  }
}

# Operations Collector Configuration
variable "operations_collector" {
  description = "VCF Operations Collector configuration"
  type = object({
    hostname           = string
    root_user_password = optional(string)
    appliance_size     = optional(string, "small")
  })
  sensitive = true
  default   = null

  validation {
    condition     = var.operations_collector == null || contains(["small", "standard"], var.operations_collector.appliance_size)
    error_message = "Operations Collector appliance size must be small or standard."
  }
}

# Automation Configuration
variable "automation" {
  description = "VCF Automation configuration"
  type = object({
    hostname              = string
    admin_user_password   = optional(string)
    internal_cluster_cidr = string
    node_prefix           = optional(string)
    ip_pool               = list(string)
  })
  sensitive = true
  default   = null

  validation {
    condition     = var.automation == null || contains(["198.18.0.0/15", "240.0.0.0/15", "250.0.0.0/15"], var.automation.internal_cluster_cidr)
    error_message = "Internal cluster CIDR must be one of: 198.18.0.0/15, 240.0.0.0/15, 250.0.0.0/15."
  }
}