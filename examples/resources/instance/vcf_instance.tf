terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

# If you wish to automate your Cloud Builder setup
# you can use the vSphere Terraform Provider - https://github.com/hashicorp/terraform-provider-vsphere
provider "vcf" {
  installer_host     = var.installer_host
  installer_username = var.installer_username
  installer_password = var.installer_password
}

resource "vcf_instance" "sddc_1" {
  instance_id = "sddcId-1001"
  skip_esx_thumbprint_validation = true
  management_pool_name = "bringup-networkpool"
  ceip_enabled = false

  sddc_manager {
    hostname = "sddc-manager"
    ssh_password = "MnogoSl0jn@P@rol@!"
    root_user_password = "MnogoSl0jn@P@rol@!"
    local_user_password = "MnogoSl0jn@P@rol@!"
  }

  ntp_servers = [
    "10.0.0.250"
  ]

  dns {
    domain = "vrack.vsphere.local"
    name_server = "10.0.0.250"
  }

  network {
    subnet = "10.0.0.0/22"
    vlan_id = "0"
    mtu = "1500"
    network_type = "MANAGEMENT"
    gateway = "10.0.0.250"
    active_uplinks = [
      "uplink1",
      "uplink2"
    ]
  }

  network {
    active_uplinks = [
      "uplink1",
      "uplink2"
    ]

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
      "10.0.4.49"
    ]

    vlan_id = "0"
    mtu = "8940"
    network_type = "VSAN"
    gateway = "10.0.4.253"
  }

  network {
    active_uplinks = [
      "uplink1",
      "uplink2"
    ]

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
    }

    root_nsx_manager_password = "MnogoSl0jn@P@rol@!"
    nsx_admin_password = "MnogoSl0jn@P@rol@!"
    nsx_audit_password = "MnogoSl0jn@P@rol@!"
    vip_fqdn = "vip-nsx-mgmt"
    transport_vlan_id = 0
  }

  vsan {
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

    vmnic_mapping {
      vmnic = "vmnic0"
      uplink = "uplink1"
    }

    vmnic_mapping {
      vmnic = "vmnic1"
      uplink = "uplink2"
    }

    networks = [
      "MANAGEMENT",
      "VSAN",
      "VMOTION"
    ]
  }

  cluster {
    datacenter_name = "SDDC-Datacenter"
    cluster_name = "SDDC-Cluster1"

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

  vcenter {
    vcenter_hostname = "vcenter-1"
    root_vcenter_password = "MnogoSl0jn@P@rol@!"
    vm_size = "tiny"
  }

  host {
    credentials {
      username = "root"
      password = var.esx_host1_pass
    }

    hostname = "esxi-1"
  }

  host {
    credentials {
      username = "root"
      password = var.esx_host2_pass
    }

    hostname = "esxi-2"
  }

  host {
    credentials {
      username = "root"
      password = var.esx_host3_pass
    }

    hostname = "esxi-3"
  }

  host {
    credentials {
      username = "root"
      password = var.esx_host4_pass
    }

    hostname = "esxi-4"
  }
}