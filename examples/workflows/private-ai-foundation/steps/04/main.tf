# This step creates a subscribed content library.
# You can use a custom subscription URL or you can use
# the content published by Broadcom

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

# Request the vSAN datastore you created in step 2 as part of your workload domain
data "vsphere_datastore" "vsan_ds" {
	datacenter_id = "${data.vsphere_datacenter.dc.id}"
	name = var.vsan_datastore_name
}

# Create a subscribed Content Library
resource "vsphere_content_library" "library" {
  name            = var.content_library_name
  storage_backing = [ data.vsphere_datastore.vsan_ds.id ]
	subscription {
	  subscription_url = var.content_library_url
	}
}