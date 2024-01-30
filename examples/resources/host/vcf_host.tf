terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

resource "vcf_network_pool" "eng_pool" {
  name = "engineering-pool"
  network {
    gateway = "192.168.8.1"
    mask    = "255.255.255.0"
    mtu     = 9000
    subnet  = "192.168.8.0"
    type    = "VSAN"
    vlan_id = 100
    ip_pools {
      start = "192.168.8.5"
      end   = "192.168.8.50"
    }
  }
  network {
    gateway = "192.168.9.1"
    mask    = "255.255.255.0"
    mtu     = 9000
    subnet  = "192.168.9.0"
    type    = "vMotion"
    vlan_id = 100
    ip_pools {
      start = "192.168.9.5"
      end   = "192.168.9.50"
    }
  }
}

resource "vcf_host" "host1" {
  fqdn            = var.host_fqdn
  username        = var.host_ssh_user
  password        = var.host_ssh_pass
  network_pool_id = vcf_network_pool.eng_pool.id
  storage_type    = "VSAN"
}