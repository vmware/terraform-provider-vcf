terraform {
  required_providers {
    vcf = {
      source  = "vmware/vcf"
    }
  }
}

provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

resource "vcf_network_pool" "domain_pool" {
  name    = "cluster-pool"
  network {
    gateway   = "192.168.12.1"
    mask      = "255.255.255.0"
    mtu       = 9000
    subnet    = "192.168.12.0"
    type      = "VSAN"
    vlan_id   = 100
    ip_pools {
      start = "192.168.12.5"
      end   = "192.168.12.50"
    }
  }
  network {
    gateway   = "192.168.13.1"
    mask      = "255.255.255.0"
    mtu       = 9000
    subnet    = "192.168.13.0"
    type      = "vMotion"
    vlan_id   = 100
    ip_pools {
      start = "192.168.13.5"
      end   = "192.168.13.50"
    }
  }
}

resource "vcf_host" "host1" {
  fqdn      = var.esx_host1_fqdn
  username  = "root"
  password  = var.esx_host1_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type = "VSAN"
}

resource "vcf_host" "host2" {
  fqdn      = var.esx_host2_fqdn
  username  = "root"
  password  = var.esx_host2_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type = "VSAN"
}
resource "vcf_host" "host3" {
  fqdn      = var.esx_host3_fqdn
  username  = "root"
  password  = var.esx_host3_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type = "VSAN"
}

resource "vcf_cluster" "cluster1" {
  // Here you can reference a terraform managed domain's id
  domain_id = var.domain_id
  name = "sfo-m01-cl01"
  host {
    id = vcf_host.host1.id
    license_key = var.esx_license_key
    vmnic {
      id = "vmnic0"
      vds_name = "sfo-m01-cl01-vds01"
    }
    vmnic {
      id = "vmnic1"
      vds_name = "sfo-m01-cl01-vds01"
    }
  }
  host {
    id = vcf_host.host2.id
    license_key = var.esx_license_key
    vmnic {
      id = "vmnic0"
      vds_name = "sfo-m01-cl01-vds01"
    }
    vmnic {
      id = "vmnic1"
      vds_name = "sfo-m01-cl01-vds01"
    }
  }
  host {
    id = vcf_host.host3.id
    license_key = var.esx_license_key
    vmnic {
      id = "vmnic0"
      vds_name = "sfo-m01-cl01-vds01"
    }
    vmnic {
      id = "vmnic1"
      vds_name = "sfo-m01-cl01-vds01"
    }
  }
  vds {
    name = "sfo-m01-cl01-vds01"
    portgroup {
      name = "sfo-m01-cl01-vds01-pg-mgmt"
      transport_type = "MANAGEMENT"
    }
    portgroup {
      name = "sfo-m01-cl01-vds01-pg-vsan"
      transport_type = "VSAN"
    }
    portgroup {
      name = "sfo-m01-cl01-vds01-pg-vmotion"
      transport_type = "VMOTION"
    }
  }
  vsan_datastore {
    datastore_name = "sfo-m01-cl01-ds-vsan01"
    failures_to_tolerate = 1
    license_key = var.vsan_license_key
  }
  geneve_vlan_id = 3
}