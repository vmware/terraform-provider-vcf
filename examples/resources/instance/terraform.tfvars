# VCF Provider Configuration
installer_host       = "sfo-ins01.sfo.rainpole.io"
installer_username   = "admin@local"
installer_password   = "VMware123!VMware123!"
allow_unverified_tls = true

# VCF Instance Configuration
instance_id                    = "sfo-vcf01"
skip_esx_thumbprint_validation = true
network_pool_name              = "sfo-m01-np01"
ceip_enabled                   = false
fips_enabled                   = false
vcf_version                    = "9.0.0.0"

# SDDC Manager Configuration
sddc_manager_info = {
  hostname            = "sfo-vcf01.sfo.rainpole.io"
  ssh_password        = "VMware123!VMware123!"
  root_user_password  = "VMware123!VMware123!"
  local_user_password = "VMware123!VMware123!"
}

# NTP Configuration
ntp_servers = [
  "10.11.10.4",
  "10.11.10.5"
]

# DNS Configuration
dns_info = {
  domain               = "sfo.rainpole.io"
  primary_nameserver   = "10.11.10.4"
  secondary_nameserver = "10.11.10.5"
}

# Management Network Configuration
management_network = {
  subnet         = "10.11.11.0/24"
  vlan_id        = 1111
  mtu            = 1500
  network_type   = "MANAGEMENT"
  gateway        = "10.11.11.1"
  active_uplinks = ["uplink1", "uplink2"]
}

# vMotion Network Configuration
vmotion_network = {
  subnet         = "10.11.12.0/24"
  vlan_id        = 1112
  mtu            = 8900
  network_type   = "VMOTION"
  gateway        = "10.11.12.1"
  active_uplinks = ["uplink1", "uplink2"]
}

# vSAN Network Configuration
vsan_network = {
  subnet         = "10.11.13.0/24"
  vlan_id        = 1113
  mtu            = 8900
  network_type   = "VSAN"
  gateway        = "10.11.13.1"
  active_uplinks = ["uplink1", "uplink2"]
}

# VM Management Network Configuration
vm_management_network = {
  subnet         = "10.11.10.0/24"
  vlan_id        = 1110
  mtu            = 1500
  network_type   = "VM_MANAGEMENT"
  gateway        = "10.11.10.1"
  active_uplinks = ["uplink1", "uplink2"]
}

# IP Address Range Configurations
vmotion_ip_address_range = {
  ip_start = "10.11.12.101"
  ip_end   = "10.11.12.108"
}

vsan_ip_address_range = {
  ip_start = "10.11.13.101"
  ip_end   = "10.11.13.108"
}

overlay_ip_address_range = {
  ip_start = "10.11.14.10"
  ip_end   = "10.11.14.50"
}

# NSX Configuration
nsx_manager_size = "medium"

nsx_manager_hostname = [
  "sfo-m01-nsx01a.sfo.rainpole.io"
]

nsx_passwords = {
  root_nsx_manager_password = "VMware123!VMware123!"
  nsx_admin_password        = "VMware123!VMware123!"
  nsx_audit_password        = "VMware123!VMware123!"
}

tep_ip_address_pool = {
  name              = "tep-ip-pool"
  subnet            = "10.11.14.0/24"
  gateway           = "10.11.14.1"
  transport_vlan_id = 1114
}

nsx_vip_fqdn              = "sfo-m01-nsx01.sfo.rainpole.io"
nsx_vlan_transportzone    = "nsx-vlan-transportzone"
nsx_overlay_transportzone = "overlay-tz-sfo-m01-nsx01"

# vSAN Configuration
vsan_config = {
  datastore_name       = "sfo-m01-vsan01"
  license_key          = null
  failures_to_tolerate = 1
  esa_enabled          = true
}

# DVS Configuration
dvs_name = "sfo-m01-dvs01"
dvs_mtu  = 8900

vmnic_mappings = [
  {
    vmnic  = "vmnic0"
    uplink = "uplink1"
  },
  {
    vmnic  = "vmnic1"
    uplink = "uplink2"
  }
]

# Cluster Configuration
datacenter_name = "sfo-m01-dc01"
cluster_name    = "sfo-m01-cl01"

# vCenter Configuration
vcenter_config = {
  hostname      = "sfo-m01-vc01.sfo.rainpole.io"
  root_password = "VMware123!VMware123!"
  vm_size       = "medium"
  storage_size  = "lstorage"
}

# Host Configuration - Your 4 ESXi hosts
hosts = [
  {
    hostname = "sfo01-m01-r01-esx01.sfo.rainpole.io"
    credentials = {
      username = "root"
      password = "VMware123!"
    }
  },
  {
    hostname = "sfo01-m01-r01-esx02.sfo.rainpole.io"
    credentials = {
      username = "root"
      password = "VMware123!"
    }
  },
  {
    hostname = "sfo01-m01-r01-esx03.sfo.rainpole.io"
    credentials = {
      username = "root"
      password = "VMware123!"
    }
  },
  {
    hostname = "sfo01-m01-r01-esx04.sfo.rainpole.io"
    credentials = {
      username = "root"
      password = "VMware123!"
    }
  }
]

# Operations Fleet Management Configuration
operations_fleet_management = {
  hostname            = "flt-fm01.rainpole.io"
  root_user_password  = "VMware123!VMware123!"
  admin_user_password = "VMware123!VMware123!"
}

# Operations Configuration
operations = {
  node = [
    {
      hostname           = "flt-ops01a.rainpole.io"
      root_user_password = "VMware123!VMware123!"
      type               = "master"
    },
    {
      hostname           = "flt-ops01b.rainpole.io"
      root_user_password = "VMware123!VMware123!"
      type               = "replica"
    },
    {
      hostname           = "flt-ops01c.rainpole.io"
      root_user_password = "VMware123!VMware123!"
      type               = "data"
    }
  ]
  admin_user_password = "VMware123!VMware123!"
  appliance_size      = "medium"
  load_balancer_fqdn  = "flt-ops01.rainpole.io"
}

# Operations Collector Configuration
operations_collector = {
  hostname           = "sfo-opsc01.sfo.rainpole.io"
  root_user_password = "VMware123!VMware123!"
  appliance_size     = "small"
}

# Automation Configuration
automation = {
  hostname              = "flt-auto01.rainpole.io"
  admin_user_password   = "VMware123!VMware123!"
  internal_cluster_cidr = "198.18.0.0/15"
  node_prefix           = "sfo-automation"
  ip_pool               = ["10.11.10.106", "10.11.10.107", "10.11.10.108", "10.11.10.109"]
}