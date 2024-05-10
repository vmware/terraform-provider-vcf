# This step generates a custom host image with vGPU drivers
# on the vCenter for the management domain.
# The source for the offline software depot for this step has to
# contain the drivers.

terraform {
  required_providers {
    vsphere = {
      source = "hashicorp/vsphere"
    }
  }
}

# Connect to the vCenter Server backing the management domain
provider "vsphere" {
  user                 = var.vcenter_username
  password             = var.vcenter_password
  vsphere_server       = var.vcenter_server
}

# Read a datacenter. Can be any datacenter
data "vsphere_datacenter" "dc" {
  name = var.datacenter_name
}

# Retrieve the list of available host images from vLCM
# It is also valid to base your custom image on the build of a particular host
# but this scenario is not automated
data "vsphere_host_base_images" "base_images" {}

# Create an offline software depot
# The source for the depot should contain the vGPU drivers 
resource "vsphere_offline_software_depot" "depot" {
  location = var.depot_location
}

# Create a compute cluster
# It will remain empty and its sole purpose is to be used by vLCM to configure
# a custom image with the GPU drivers
resource "vsphere_compute_cluster" "image_source_cluster" {
  name            = var.cluster_name
  datacenter_id   = data.vsphere_datacenter.dc.id

  # The "host_image" block enables vLCM on the cluster and configures a custom image with the provided settings
  # It is recommended to add this block after you have configured your depot and retrieved the list of base images
  # so that you can select the correct values
  # This example uses the first available image and the first available component
  host_image {
    esx_version = data.vsphere_host_base_images.base_images.version.0
    component {
      key = vsphere_offline_software_depot.depot.component.0.key
      version = vsphere_offline_software_depot.depot.component.0.version.0
    }
  }
}