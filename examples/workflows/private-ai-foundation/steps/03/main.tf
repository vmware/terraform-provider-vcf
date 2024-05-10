terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

# Connect to the SDDC Manager
provider "vcf" {
  sddc_manager_host     = var.sddc_manager_host
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
}

# Create the Edge Cluster
# This sample will create a cluster with the minimum number of nodes - 2
# You can add more nodes if necessary. The VCF Terraform provider supports Edge Cluster expansion and shrinkage
# which means you can add more nodes later too
resource "vcf_edge_cluster" "edge_cluster" {
  name              = var.edge_cluster_name
  root_password     = var.edge_cluster_root_pass
  admin_password    = var.edge_cluster_admin_pass
  audit_password    = var.edge_cluster_audit_pass
  tier0_name        = "T0_myCluster2"
  tier1_name        = "T1_myCluster2"
  form_factor       = "MEDIUM"
  profile_type      = "DEFAULT"
  routing_type      = "EBGP"
  high_availability = "ACTIVE_ACTIVE"
  mtu               = 9000
  asn               = 65004
  skip_tep_routability_check = true

  edge_node {
    name               = "${var.edge_cluster_name}-node-1.vrack.vsphere.local"
    # The ID of the compute cluster is available on the domain resource from step 2
    compute_cluster_id = var.cluster_id
    root_password      = var.edge_node1_root_pass
    admin_password     = var.edge_node1_admin_pass
    audit_password     = var.edge_node1_audit_pass
    management_ip      = var.edge_node1_cidr
    management_gateway = "10.0.0.250"
    tep_gateway        = "192.168.52.1"
    tep1_ip            = "192.168.52.10/24"
    tep2_ip            = "192.168.52.11/24"
    tep_vlan           = 100
    inter_rack_cluster = false

    uplink {
      vlan         = 2083
      interface_ip = "192.168.18.2/24"
      bgp_peer {
        ip       = "192.168.18.10/24"
        password = var.bgp_peer_password
        asn      = 65001
      }
    }

    uplink {
      vlan         = 2084
      interface_ip = "192.168.19.2/24"
      bgp_peer {
        ip       = "192.168.19.10/24"
        password = var.bgp_peer_password
        asn      = 65001
      }
    }
  }

  edge_node {
    name               = "${var.edge_cluster_name}-node-2.vrack.vsphere.local"
    compute_cluster_id = var.cluster_id
    root_password      = var.edge_node2_root_pass
    admin_password     = var.edge_node2_admin_pass
    audit_password     = var.edge_node2_audit_pass
    management_ip      = var.edge_node2_cidr
    management_gateway = "10.0.0.250"
    tep_gateway        = "192.168.52.1"
    tep1_ip            = "192.168.52.12/24"
    tep2_ip            = "192.168.52.13/24"
    tep_vlan           = 100
    inter_rack_cluster = false

    uplink {
      vlan         = 2083
      interface_ip = "192.168.18.3/24"
      bgp_peer {
        ip       = "192.168.18.10/24"
        password = var.bgp_peer_password
        asn      = 65001
      }
    }

    uplink {
      vlan         = 2084
      interface_ip = "192.168.19.3/24"
      bgp_peer {
        ip       = "192.168.19.10/24"
        password = var.bgp_peer_password
        asn      = 65001
      }
    }
  }
}