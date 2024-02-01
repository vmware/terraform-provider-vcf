---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vcf_edge_cluster Resource - terraform-provider-vcf"
subcategory: ""
description: |-
  
---

# vcf_edge_cluster (Resource)


Prerequisites for creating an NSX edge cluster
* The following conditions must be met:
    * Separate VLANs must be available for Host TEP VLAN and Edge TEP VLAN use
    * Host TEP VLAN and Edge TEP VLAN need to be routed
    * If dynamic routing is desired, two BGP peers (on TORs or infra ESG) with an interface IP, ASN and BGP password are required
    * An ASN has to be reserved for the NSX Edge cluster's Tier-0 interfaces
    * DNS entries for NSX Edge components should be populated in a customer-managed DNS server
    * The vSphere clusters hosting the Edge clusters should be L2 Uniform. All host nodes in a hosting vSphere cluster need to have identical management, uplink, Edge and host TEP networks
    * The vSphere clusters hosting the NSX Edge node VMs must have the same pNIC speed for NSX enabled VDS uplinks chosen for Edge overlay (e.g., either 10G or 25G but not both)
    * All nodes of an NSX Edge cluster must use the same set of NSX enabled VDS uplinks. The selected uplinks must be prepared for overlay use
    * If the vSphere cluster hosting the Edge nodes has hosts with a DPU device then enable SR-IOV in the BIOS and in the vSphere Client (if required by your DPU vendor)

Do not attempt to add and remove edge nodes in a single configuration change. You can either shrink or expand a cluster, but you cannot run both operations
simultaneously.

Review the documentation for VMware Cloud Foundation for more information about NSX Edge Clusters.


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `admin_password` (String) Administrator password for the NSX manager
- `audit_password` (String) Audit user password for the NSX manager
- `edge_node` (Block List, Min: 1) The nodes in the edge cluster (see [below for nested schema](#nestedblock--edge_node))
- `form_factor` (String) One among: XLARGE, LARGE, MEDIUM, SMALL
- `high_availability` (String) One among: ACTIVE_ACTIVE, ACTIVE_STANDBY
- `mtu` (Number) Maximum transmission unit size for the cluster
- `name` (String) The name of the edge cluster
- `profile_type` (String) One among: DEFAULT, CUSTOM. If set to CUSTOM a 'profile' must be provided
- `root_password` (String) Root user password for the NSX manager
- `routing_type` (String) One among: EBGP, STATIC
- `tier0_name` (String) Name for the Tier-0 gateway
- `tier1_name` (String) Name for the Tier-1 gateway

### Optional

- `asn` (Number) ASN for the cluster
- `internal_transit_subnets` (List of String) Subnet addresses in CIDR notation that are used to assign addresses to logical links connecting service routers and distributed routers
- `profile` (Block List, Max: 1) The specification for the edge cluster profile (see [below for nested schema](#nestedblock--profile))
- `skip_tep_routability_check` (Boolean) Set to true to bypass normal ICMP-based check of Edge TEP / host TEP routability (default is false, meaning do check)
- `tier1_unhosted` (Boolean) Select whether Tier-1 being created per this spec is hosted on the new Edge cluster or not (default value is false, meaning hosted)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `transit_subnets` (List of String) Transit subnet addresses in CIDR notation that are used to assign addresses to logical links connecting Tier-0 and Tier-1s

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--edge_node"></a>
### Nested Schema for `edge_node`

Required:

- `admin_password` (String) The administrator password for the edge node
- `audit_password` (String) The audit password for the edge node
- `compute_cluster_id` (String) The id of the compute cluster
- `inter_rack_cluster` (Boolean) Whether or not this is an inter-rack cluster. True for L2 non-uniform and L3, false for L2 uniform
- `management_gateway` (String) The gateway address for the management network
- `management_ip` (String) The IP address (CIDR) for the management network
- `name` (String) The name of the edge node
- `root_password` (String) The root user password for the edge node
- `tep1_ip` (String) The IP address (CIDR) of the first tunnel endpoint
- `tep2_ip` (String) The IP address (CIDR) of the second tunnel endpoint
- `tep_gateway` (String) The gateway for the tunnel endpoints
- `tep_vlan` (Number) The VLAN ID for the tunnel endpoint

Optional:

- `first_nsx_vds_uplink` (String) The name of the first NSX-enabled VDS uplink
- `second_nsx_vds_uplink` (String) The name of the second NSX-enabled VDS uplink
- `uplink` (Block List) Specifications of Tier-0 uplinks for the edge node (see [below for nested schema](#nestedblock--edge_node--uplink))

<a id="nestedblock--edge_node--uplink"></a>
### Nested Schema for `edge_node.uplink`

Required:

- `interface_ip` (String) The IP address (CIDR) for the distributed switch uplink
- `vlan` (Number) The VLAN ID for the distributed switch uplink

Optional:

- `bgp_peer` (Block List) List of BGP Peer configurations (see [below for nested schema](#nestedblock--edge_node--uplink--bgp_peer))

<a id="nestedblock--edge_node--uplink--bgp_peer"></a>
### Nested Schema for `edge_node.uplink.bgp_peer`

Required:

- `asn` (Number) ASN
- `ip` (String) IP address
- `password` (String) Password




<a id="nestedblock--profile"></a>
### Nested Schema for `profile`

Required:

- `bfd_allowed_hop` (Number) BFD allowed hop
- `bfd_declare_dead_multiple` (Number) BFD declare dead multiple
- `bfd_probe_interval` (Number) BFD probe interval
- `name` (String) The name of the profile
- `standby_relocation_threshold` (Number) Standby relocation threshold


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `update` (String)