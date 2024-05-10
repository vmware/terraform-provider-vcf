terraform {
  required_providers {
    vsphere = {
      source = "hashicorp/vsphere"
    }
  }
}

# Connect to the vCenter Server backing the workload domain
provider "vsphere" {
  user                 = var.vcenter_username
  password             = var.vcenter_password
  vsphere_server       = var.vcenter_server
}

# Request the Datacenter you created in step 2 as part of your workload domain
data "vsphere_datacenter" "dc" {
  name = var.datacenter_name
}

# Request the Cluster you created in step 2 as part of your workload domain
data "vsphere_compute_cluster" "vi_cluster" {
  name            = var.cluster_name
  datacenter_id   = data.vsphere_datacenter.dc.id
}

# Request the default vSAN storage policy created as part of your workload domain
data vsphere_storage_policy image_policy {
	name = var.storage_policy_name
}

# Request a management network. This has to be a distributed portgroup on the DVS which the edge nodes are connected to
data vsphere_network mgmt_net {
	name = var.management_network_name
	datacenter_id = data.vsphere_datacenter.dc.id
}

# Request the content library you created in step 4
data vsphere_content_library subscribed_lib {
	name = var.contenty_library_name
}

# Request the DVS which the edge nodes are connected to
data vsphere_distributed_virtual_switch dvs {
	name = var.dvs_name
	datacenter_id = data.vsphere_datacenter.dc.id
}

# Request any one of the hosts in your cluster
# This sample assumes that every host on your cluster has the same vGPU
data "vsphere_host" "host1" {
    name = var.host1_fqdn
    datacenter_id = data.vsphere_datacenter.dc.id
}

# Request the list of available vGPU profiles
# This sample assumes that the vGPU settings are the same on all host in the cluster
data "vsphere_host_vgpu_profile" "vgpus" {
  host_id = data.vsphere_host.host1.id
}

# Create a virtual machine class with the vGPU profile
# It is recommended to run this after you have retrieved and read the list of vGPU profiles
# in order to determine which one you want to use.
# This sample just uses the first available profile
resource "vsphere_virtual_machine_class" "vgpu_class" {
	name = var.virtual_machine_class_name
	# These resource settings will be insufficient for actual AI computations
	# Increase these according to your needs
	cpus = 4
	memory = 4096
	# 100% memory reservation is required when using vGPUs
	memory_reservation = 100
    vgpu_devices = [ data.vsphere_host_vgpu_profile.vgpus.vgpu_profiles[0].vgpu  ]
}

# Enable Supervisor on your cluster
resource "vsphere_supervisor" "supervisor" {
	cluster = data.vsphere_compute_cluster.vi_cluster.id
	storage_policy = data.vsphere_storage_policy.image_policy.id
	content_library = data.vsphere_content_library.subscribed_lib.id
	main_dns = "10.0.0.250"
	worker_dns = "10.0.0.250"
	edge_cluster = var.edge_cluster
	dvs_uuid = data.vsphere_distributed_virtual_switch.dvs.id
	sizing_hint = "MEDIUM"
	
	management_network {
		network = data.vsphere_network.mgmt_net.id
		subnet_mask = "255.255.255.0"
		starting_address = "10.0.0.150"
		gateway = "10.0.0.250"
		address_count = 5
	}

	ingress_cidr {
		address = "10.10.10.0"
		prefix = 24
	}

	egress_cidr {
		address = "10.10.11.0"
		prefix = 24
	}

	pod_cidr {
		address = "10.244.10.0"
		prefix = 23
	}

	service_cidr {
		address = "10.10.12.0"
		prefix = 24
	}

	search_domains = [ "vrack.vsphere.local" ]

	namespace {
		name = var.namespace_name
		content_libraries = [ data.vsphere_content_library.subscribed_lib.id ]
        vm_classes = [ vsphere_virtual_machine_class.vgpu_class.id ]
	}
}