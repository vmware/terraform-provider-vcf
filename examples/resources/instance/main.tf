provider "vcf" {
  installer_host       = var.installer_host
  installer_username   = var.installer_username
  installer_password   = var.installer_password
  allow_unverified_tls = var.allow_unverified_tls
}

resource "vcf_instance" "main" {
  instance_id                    = var.instance_id
  skip_esx_thumbprint_validation = var.skip_esx_thumbprint_validation
  management_pool_name           = var.network_pool_name
  ceip_enabled                   = var.ceip_enabled
  fips_enabled                   = var.fips_enabled
  version                        = var.vcf_version

  sddc_manager {
    hostname            = var.sddc_manager_info["hostname"]
    ssh_password        = var.sddc_manager_info["ssh_password"]
    root_user_password  = var.sddc_manager_info["root_user_password"]
    local_user_password = var.sddc_manager_info["local_user_password"]
  }

  ntp_servers = var.ntp_servers

  dns {
    domain                = var.dns_info.domain
    name_server           = var.dns_info.primary_nameserver
    secondary_name_server = var.dns_info.secondary_nameserver
  }

  network {
    subnet         = var.management_network.subnet
    vlan_id        = var.management_network.vlan_id
    mtu            = var.management_network.mtu
    network_type   = var.management_network.network_type
    gateway        = var.management_network.gateway
    active_uplinks = var.management_network.active_uplinks
  }

  network {
    subnet         = var.vmotion_network.subnet
    vlan_id        = var.vmotion_network.vlan_id
    mtu            = var.vmotion_network.mtu
    network_type   = var.vmotion_network.network_type
    gateway        = var.vmotion_network.gateway
    active_uplinks = var.vmotion_network.active_uplinks

    include_ip_address_ranges {
      start_ip_address = var.vmotion_ip_address_range.ip_start
      end_ip_address   = var.vmotion_ip_address_range.ip_end
    }
  }

  network {
    subnet         = var.vsan_network.subnet
    vlan_id        = var.vsan_network.vlan_id
    mtu            = var.vsan_network.mtu
    network_type   = var.vsan_network.network_type
    gateway        = var.vsan_network.gateway
    active_uplinks = var.vsan_network.active_uplinks

    include_ip_address_ranges {
      start_ip_address = var.vsan_ip_address_range.ip_start
      end_ip_address   = var.vsan_ip_address_range.ip_end
    }
  }

  network {
    subnet         = var.vm_management_network.subnet
    vlan_id        = var.vm_management_network.vlan_id
    mtu            = var.vm_management_network.mtu
    network_type   = var.vm_management_network.network_type
    gateway        = var.vm_management_network.gateway
    active_uplinks = var.vm_management_network.active_uplinks
  }

  nsx {
    nsx_manager_size = var.nsx_manager_size

    dynamic "nsx_manager" {
      for_each = var.nsx_manager_hostname
      content {
        hostname = nsx_manager.value
      }
    }

    root_nsx_manager_password = var.nsx_passwords.root_nsx_manager_password
    nsx_admin_password        = var.nsx_passwords.nsx_admin_password
    nsx_audit_password        = var.nsx_passwords.nsx_audit_password

    ip_address_pool {
      name = var.tep_ip_address_pool.name

      subnet {
        cidr    = var.tep_ip_address_pool.subnet
        gateway = var.tep_ip_address_pool.gateway

        ip_address_pool_range {
          start = var.overlay_ip_address_range.ip_start
          end   = var.overlay_ip_address_range.ip_end
        }
      }
    }

    vip_fqdn          = var.nsx_vip_fqdn
    transport_vlan_id = var.tep_ip_address_pool.transport_vlan_id
  }

  vsan {
    datastore_name       = var.vsan_config.datastore_name
    failures_to_tolerate = var.vsan_config.failures_to_tolerate
    esa_enabled          = var.vsan_config.esa_enabled
  }

  dvs {
    dvs_name = var.dvs_name
    mtu      = var.dvs_mtu

    dynamic "vmnic_mapping" {
      for_each = var.vmnic_mappings
      content {
        vmnic  = vmnic_mapping.value.vmnic
        uplink = vmnic_mapping.value.uplink
      }
    }

    nsxt_switch_config {
      transport_zones {
        name           = var.nsx_vlan_transportzone
        transport_type = "VLAN"
      }
      transport_zones {
        name           = var.nsx_overlay_transportzone
        transport_type = "OVERLAY"
      }
      host_switch_operational_mode = "ENS_INTERRUPT"
    }

    nsx_teaming {
      policy = "LOADBALANCE_SRCID"
      active_uplinks = [
        "uplink1",
        "uplink2"
      ]
    }

    networks = [
      "MANAGEMENT",
      "VSAN",
      "VMOTION",
      "VM_MANAGEMENT"
    ]
  }

  cluster {
    datacenter_name = var.datacenter_name
    cluster_name    = var.cluster_name
  }

  vcenter {
    vcenter_hostname      = var.vcenter_config.hostname
    root_vcenter_password = var.vcenter_config.root_password
    vm_size               = var.vcenter_config.vm_size
    storage_size          = var.vcenter_config.storage_size
  }

  dynamic "host" {
    for_each = var.hosts
    content {
      hostname = host.value.hostname
      credentials {
        username = host.value.credentials.username
        password = host.value.credentials.password
      }
    }
  }

  dynamic "operations_fleet_management" {
    for_each = var.operations_fleet_management != null ? [var.operations_fleet_management] : []
    content {
      hostname            = operations_fleet_management.value.hostname
      root_user_password  = operations_fleet_management.value.root_user_password
      admin_user_password = operations_fleet_management.value.admin_user_password
    }
  }

  dynamic "operations" {
    for_each = var.operations != null ? [var.operations] : []
    content {
      dynamic "node" {
        for_each = operations.value.node
        content {
          hostname           = node.value.hostname
          root_user_password = node.value.root_user_password
          type               = node.value.type
        }
      }
      admin_user_password = operations.value.admin_user_password
      appliance_size      = operations.value.appliance_size
      load_balancer_fqdn  = operations.value.load_balancer_fqdn
    }
  }

  dynamic "operations_collector" {
    for_each = var.operations_collector != null ? [var.operations_collector] : []
    content {
      hostname           = operations_collector.value.hostname
      root_user_password = operations_collector.value.root_user_password
      appliance_size     = operations_collector.value.appliance_size
    }
  }

  dynamic "automation" {
    for_each = var.automation != null ? [var.automation] : []
    content {
      hostname              = automation.value.hostname
      admin_user_password   = automation.value.admin_user_password
      internal_cluster_cidr = automation.value.internal_cluster_cidr
      node_prefix           = automation.value.node_prefix
      ip_pool               = automation.value.ip_pool
    }
  }
}