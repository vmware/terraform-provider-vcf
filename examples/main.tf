# Required Provider
terraform {
  required_providers {
    vcf = {
      source  = "vmware/vcf"
    }
  }
}

# Provider Configuration
provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

# Resources
resource "vcf_ceip" "ceip" {
  status    = "DISABLE"
}

resource "vcf_user" "testuser1" {
  name      = "testuser1@vrack.vsphere.local"
  domain    = "vrack.vsphere.local"
  type      = "USER"
  role_name = "VIEWER"
}

resource "vcf_user" "serviceuser1" {
  name      = "serviceuser1"
  domain    = "Nil"
  type      = "SERVICE"
  role_name = "ADMIN"
}

output "serviceuser1_apikey" {
  value = "${vcf_user.serviceuser1.api_key}"
}

resource "vcf_network_pool" "eng_pool" {
  name    = "engineering-pool2"
  network {
    gateway   = "192.168.8.1"
    mask      = "255.255.255.0"
    mtu       = 9000
    subnet    = "192.168.8.0"
    type      = "VSAN"
    vlan_id   = 100
    ip_pools {
      start = "192.168.8.5"
      end   = "192.168.8.50"
    }
  }
  network {
    gateway   = "192.168.9.1"
    mask      = "255.255.255.0"
    mtu       = 9000
    subnet    = "192.168.9.0"
    type      = "vMotion"
    vlan_id   = 100
    ip_pools {
      start = "192.168.9.5"
      end   = "192.168.9.50"
    }
  }
}

resource "vcf_host" "esxi_1" {
  fqdn      = "esxi-1.vrack.vsphere.local"
  username  = var.esxi_1_user
  password  = var.esxi_1_pass
  network_pool_name = "bringup-networkpool"
  storage_type = "VSAN"
}

resource "vcf_host" "esxi_2" {
  fqdn      = "esxi-2.vrack.vsphere.local"
  username  = var.esxi_2_user
  password  = var.esxi_2_pass
  network_pool_name = "bringup-networkpool"
  storage_type = "VSAN"
}

resource "vcf_host" "esxi_3" {
  fqdn      = "esxi-3.vrack.vsphere.local"
  username  = var.esxi_3_user
  password  = var.esxi_3_pass
  network_pool_name = "bringup-networkpool"
  storage_type = "VSAN"
}