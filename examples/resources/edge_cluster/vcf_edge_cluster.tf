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

resource "vcf_edge_cluster" "cluster_1" {
  name              = var.cluster_name
  root_password     = var.cluster_root_pass
  admin_password    = var.cluster_admin_pass
  audit_password    = var.cluster_audit_pass
  tier0_name        = "T0_cluster_1"
  tier1_name        = "T1_cluster_1"
  form_factor       = "MEDIUM"
  profile_type      = "DEFAULT"
  routing_type      = "EBGP"
  high_availability = "ACTIVE_ACTIVE"
  mtu               = 9000
  asn               = "65004"

  edge_node {
    name               = var.edge_node_1_name
    compute_cluster_id = var.compute_cluster_id
    root_password      = var.edge_node_1_root_pass
    admin_password     = var.edge_node_1_admin_pass
    audit_password     = var.edge_node_1_audit_pass
    management_ip      = "10.0.0.52/24"
    management_gateway = "10.0.0.250"
    tep_gateway        = "192.168.52.1"
    tep1_ip            = "192.168.52.12/24"
    tep2_ip            = "192.168.52.13/24"
    tep_vlan           = 1252
    inter_rack_cluster = false

    uplink {
      vlan         = 2083
      interface_ip = "192.168.18.2/24"
      bgp_peer {
        ip       = "192.168.18.10/24"
        password = "VMware1!"
        asn      = "65001"
      }
    }

    uplink {
      vlan         = 2084
      interface_ip = "192.168.19.2/24"
      bgp_peer {
        ip       = "192.168.19.10/24"
        password = "VMware1!"
        asn      = "65001"
      }
    }
  }

  edge_node {
    name               = var.edge_node_2_name
    compute_cluster_id = var.compute_cluster_id
    root_password      = var.edge_node_2_root_pass
    admin_password     = var.edge_node_2_admin_pass
    audit_password     = var.edge_node_2_audit_pass
    management_ip      = "10.0.0.53/24"
    management_gateway = "10.0.0.250"
    tep_gateway        = "192.168.52.1"
    tep1_ip            = "192.168.52.14/24"
    tep2_ip            = "192.168.52.15/24"
    tep_vlan           = 1252
    inter_rack_cluster = false

    uplink {
      vlan         = 2083
      interface_ip = "192.168.18.3/24"
      bgp_peer {
        ip       = "192.168.18.10/24"
        password = "VMware1!"
        asn      = "65001"
      }
    }

    uplink {
      vlan         = 2084
      interface_ip = "192.168.19.3/24"
      bgp_peer {
        ip       = "192.168.19.10/24"
        password = "VMware1!"
        asn      = "65001"
      }
    }
  }
}
