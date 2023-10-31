terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}
provider "vcf" {
  cloud_builder_host = var.cloud_builder_host
  cloud_builder_username = var.cloud_builder_username
  cloud_builder_password = var.cloud_builder_password
}

resource "vcf_instance" "sddc_1" {
  instance_id = "sddcId-1001"
  dv_switch_version = "7.0.0"
  skip_esx_thumbprint_validation = true
  management_pool_name = "bringup-networkpool"
  ceip_enabled = false
  esx_license = ""
  task_name = "workflowconfig/workflowspec-ems.json"
  sddc_manager {
    ip_address = "10.0.0.4"
    netmask = "255.255.255.0"
    hostname = "sddc-manager"
    root_user_credentials {
      username = "root"
      password = var.sddc_manager_root_user_password
    }
    second_user_credentials {
      username = "vcf"
      password = var.sddc_manager_secondary_user_password
    }
  }
  ntp_servers = [
    "10.0.0.250"
  ]
  dns {
    domain = "vsphere.local"
    name_server = "10.0.0.250"
    secondary_name_server = "10.0.0.250"
  }
  network {
    subnet = "10.0.0.0/22"
    vlan_id = "0"
    mtu = "1500"
    network_type = "MANAGEMENT"
    gateway = "10.0.0.250"
  }
  network {
    subnet = "10.0.4.0/24"
    include_ip_address_ranges {
      start_ip_address = "10.0.4.7"
      end_ip_address = "10.0.4.48"
    }
    include_ip_address_ranges {
      start_ip_address = "10.0.4.3"
      end_ip_address = "10.0.4.6"
    }
    include_ip_address = [
      "10.0.4.50",
      "10.0.4.49"]
    vlan_id = "0"
    mtu = "8940"
    network_type = "VSAN"
    gateway = "10.0.4.253"
  }
  network {
    subnet = "10.0.8.0/24"
    include_ip_address_ranges {
      start_ip_address = "10.0.8.3"
      end_ip_address = "10.0.8.50"
    }
    vlan_id = "0"
    mtu = "8940"
    network_type = "VMOTION"
    gateway = "10.0.8.253"
  }
  nsx {
    nsx_manager_size = "medium"
    nsx_manager {
      hostname = "nsx-mgmt-1"
      ip = "10.0.0.31"
    }
    root_nsx_manager_password = var.nsx_manager_root_password
    nsx_admin_password = var.nsx_manager_admin_password
    nsx_audit_password = var.nsx_manager_audit_password
    overlay_transport_zone {
      zone_name = "overlay-tz"
      network_name = "net-overlay"
    }
    vlan_transport_zone {
      zone_name = "vlan-tz"
      switch_name = "mgmt-nvds"
      network_name = "net-vlan"
    }
    vip = "10.0.0.30"
    vip_fqdn = "vip-nsx-mgmt"
    license = var.nsx_license_key
    transport_vlan_id = 0
  }
  vsan {
    vsan_name = "vsan-1"
    license = var.vsan_license_key
    datastore_name = "sfo01-m01-vsan"
  }
  dvs {
    mtu = 8940
    nioc {
      traffic_type = "VSAN"
      value = "HIGH"
    }
    nioc {
      traffic_type = "VMOTION"
      value = "LOW"
    }
    nioc {
      traffic_type = "VDP"
      value = "LOW"
    }
    nioc {
      traffic_type = "VIRTUALMACHINE"
      value = "HIGH"
    }
    nioc {
      traffic_type = "MANAGEMENT"
      value = "NORMAL"
    }
    nioc {
      traffic_type = "NFS"
      value = "LOW"
    }
    nioc {
      traffic_type = "HBR"
      value = "LOW"
    }
    nioc {
      traffic_type = "FAULTTOLERANCE"
      value = "LOW"
    }
    nioc {
      traffic_type = "ISCSI"
      value = "LOW"
    }
    dvs_name = "SDDC-Dswitch-Private"
    vmnics = [
      "vmnic0",
      "vmnic1"
    ]
    networks = [
      "MANAGEMENT",
      "VSAN",
      "VMOTION"
    ]
  }
  cluster {
    cluster_name = "SDDC-Cluster1"
    cluster_evc_mode = ""
    resource_pool {
      name = "Mgmt-ResourcePool"
      type = "management"
    }
    resource_pool {
      name = "Network-ResourcePool"
      type = "network"
    }
    resource_pool {
      name = "Compute-ResourcePool"
      type = "compute"
    }
    resource_pool {
      name = "User-RP"
      type = "compute"
    }
  }
  psc {
    psc_sso_domain = "vsphere.local"
    admin_user_sso_password = "TestTest!"
  }
  vcenter {
    vcenter_ip = "10.0.0.6"
    vcenter_hostname = "vcenter-1"
    license = var.vcenter_license_key
    root_vcenter_password = var.vcenter_root_password
    vm_size = "tiny"
  }
  host {
    credentials {
      username = "root"
      password = var.esx_host1_pass
    }
    ip_address_private {
      subnet = "255.255.252.0"
      cidr = ""
      ip_address = "10.0.0.100"
      gateway = "10.0.0.250"
    }
    hostname = "esxi-1"
    vswitch = "vSwitch0"
    server_id = "host-0"
    association = "SDDC-Datacenter"
  }
  host {
    credentials {
      username = "root"
      password = var.esx_host2_pass
    }
    ip_address_private {
      subnet = "255.255.252.0"
      cidr = ""
      ip_address = "10.0.0.101"
      gateway = "10.0.0.250"
    }
    hostname = "esxi-2"
    vswitch = "vSwitch0"
    association = "SDDC-Datacenter"
  }
  host {
    credentials {
      username = "root"
      password = var.esx_host3_pass
    }
    ip_address_private {
      subnet = "255.255.255.0"
      cidr = ""
      ip_address = "10.0.0.102"
      gateway = "10.0.0.250"
    }
    hostname = "esxi-3"
    vswitch = "vSwitch0"
    association = "SDDC-Datacenter"
  }
  host {
    credentials {
      username = "root"
      password = var.esx_host4_pass
    }
    ip_address_private {
      subnet = "255.255.255.0"
      cidr = ""
      ip_address = "10.0.0.103"
      gateway = "10.0.0.250"
    }
    hostname = "esxi-4"
    vswitch = "vSwitch0"
    server_id = "host-3"
    association = "SDDC-Datacenter"
  }
}