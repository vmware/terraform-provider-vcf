# This step creates a new workload domain with the custom image
# from step 1.

terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
    vsphere = {
      source = "hashicorp/vsphere"
    }
  }
}

# Connect to the SDDC Manager
provider "vcf" {
  sddc_manager_host     = var.sddc_manager_host
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
}

# Connect to the vCenter Server backing the management domain
provider "vsphere" {
  user                 = var.vcenter_username
  password             = var.vcenter_password
  vsphere_server       = var.vcenter_server
}

# Request the same datacenter which you created your cluster on in step 1
data "vsphere_datacenter" "dc" {
  name = var.datacenter_name
}

# Request the compute cluster which you created in step 1
data "vsphere_compute_cluster" "image_source_cluster" {
  name            = var.source_cluster_name
  datacenter_id   = data.vsphere_datacenter.dc.id
}

# Configure a network pool for your hosts
resource "vcf_network_pool" "domain_pool" {
  name = "engineering-pool"
  network {
    gateway = "192.168.10.1"
    mask    = "255.255.255.0"
    mtu     = 9000
    subnet  = "192.168.10.0"
    type    = "VSAN"
    vlan_id = 100
    ip_pools {
      start = "192.168.10.5"
      end   = "192.168.10.50"
    }
  }
  network {
    gateway = "192.168.11.1"
    mask    = "255.255.255.0"
    mtu     = 9000
    subnet  = "192.168.11.0"
    type    = "vMotion"
    vlan_id = 100
    ip_pools {
      start = "192.168.11.5"
      end   = "192.168.11.50"
    }
  }
}

# Commission 3 hosts for the new domain
resource "vcf_host" "host1" {
  fqdn            = var.esx_host1_fqdn
  username        = "root"
  password        = var.esx_host1_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type    = "VSAN"
}
resource "vcf_host" "host2" {
  fqdn            = var.esx_host2_fqdn
  username        = "root"
  password        = var.esx_host2_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type    = "VSAN"
}
resource "vcf_host" "host3" {
  fqdn            = var.esx_host3_fqdn
  username        = "root"
  password        = var.esx_host3_pass
  network_pool_id = vcf_network_pool.domain_pool.id
  storage_type    = "VSAN"
}

# Extract a vLCM personality (a cluster image) from the cluster you created in step 1
# This will be applied on the new workload domain
# It is crucial that you do this before creating the domain as it is not possible to enable vLCM afterwards
resource "vcf_cluster_personality" "custom_image" {
    name      = "custom-image-3"
    cluster_id = data.vsphere_compute_cluster.image_source_cluster.id
    domain_id = var.management_domain_id
}

# Create a workload domain
resource "vcf_domain" "wld01" {
  name = var.workload_domain_name

  vcenter_configuration {
    datacenter_name = "${var.workload_domain_name}-datacenter"
    fqdn            = var.workload_vcenter_fqdn
    gateway         = "10.0.0.250"
    ip_address      = var.workload_vcenter_address
    name            = "${var.workload_domain_name}-vcenter"
    root_password   = var.vcenter_root_password
    subnet_mask     = "255.255.252.0"
  }

  cluster {
    name = "vi-cluster"

    host {
      id          = vcf_host.host1.id
      license_key = var.esx_license_key
    }

    host {
      id          = vcf_host.host2.id
      license_key = var.esx_license_key
    }

    host {
      id          = vcf_host.host3.id
      license_key = var.esx_license_key
    }

    vds {
      name = "${var.workload_domain_name}-vds01"

      portgroup {
        name           = "${var.workload_domain_name}-vds01-PortGroup-Mgmt"
        transport_type = "MANAGEMENT"
      }

      portgroup {
        name           = "${var.workload_domain_name}-vds01-PortGroup-vMotion"
        transport_type = "VMOTION"
      }

      portgroup {
        name           = "${var.workload_domain_name}-vds01-PortGroup-VSAN"
        transport_type = "VSAN"
      }
    }

    vsan_datastore {
      datastore_name = "${var.workload_domain_name}-vsan"
      license_key    = var.vsan_license_key
    }

    geneve_vlan_id = "112"
    cluster_image_id = vcf_cluster_personality.custom_image.id
  }

  nsx_configuration {
    license_key                = var.nsx_license_key
    nsx_manager_admin_password = var.nsx_manager_admin_password

    # You need to prepare the DNS entries for these hostnames before running Terraform
    # You are free to modify the FQDNs
    nsx_manager_node {
      fqdn        = "nsx-mgmt-wld-1.vrack.vsphere.local"
      gateway     = "10.0.0.250"
      ip_address  = var.nsx_manager_node1_address
      name        = "nsx-mgmt-wld-1"
      subnet_mask = "255.255.252.0"
    }

    nsx_manager_node {
      fqdn        = "nsx-mgmt-wld-2.vrack.vsphere.local"
      gateway     = "10.0.0.250"
      ip_address  = var.nsx_manager_node2_address
      name        = "nsx-mgmt-wld-2"
      subnet_mask = "255.255.252.0"
    }

    nsx_manager_node {
      fqdn        = "nsx-mgmt-wld-3.vrack.vsphere.local"
      gateway     = "10.0.0.250"
      ip_address  = var.nsx_manager_node3_address
      name        = "nsx-mgmt-wld-3"
      subnet_mask = "255.255.252.0"
    }
    vip      = var.nsx_manager_vip_address
    vip_fqdn = "nsx-manager-wld.vrack.vsphere.local"
  }
}
