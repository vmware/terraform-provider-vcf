---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vcf_host Resource - terraform-provider-vcf"
subcategory: ""
description: |-
  
---

# vcf_host (Resource)

Prerequisites for commissioning Hosts

- The following data is required:

  - Username of each host
  - Password of each host
  - FQDN of each host
  - Network pool ID to which each host has to be associated with

- The host, if intended to be used for a vSAN domain, should be vSAN compliant and certified as per the VMware Hardware Compatibility Guide.
  BIOS, HBA, SSD, HDD, etc. of the host must match the VMware Hardware Compatibility Guide.
- The host must have the drivers and firmware versions specified in the VMware Hardware Compatibility Guide.
- The host must have the supported version of ESXi (i.e 6.7.0-13006603) pre-installed on it.
- SSH and syslog must be enabled on the host.
- The host must be configured with DNS server for forward and reverse lookup and FQDN.
- The host name must be same as the FQDN.
- The host must have a standard switch with two NIC ports with a minimum 10 Gbps speed.
- The management IP must be configured to the first NIC port.
- Ensure that the host has a standard switch and the default uplinks with 10Gb speed are configured starting with traditional numbering (e.g., vmnic0) and increasing sequentially.
- Ensure that the host hardware health status is healthy without any errors.
- All disk partitions on HDD / SSD are deleted.
- The hosts, if intended to be used for vSAN, domain must be associated with vSAN enabled network pool.
- The hosts, if intended to be used for NFS, domain must be associated with NFS enabled network pool.
- The hosts, if intended to be used for VMFS on FC, domain must be associated with either a NFS enabled or vMotion enabled network pool.
- The hosts, if intended to be used for VVOL, domain must be associated with either a NFS enabled or vMotion enabled network pool.
- The hosts, if intended to be used for vSAN HCI Mesh(VSAN_REMOTE), domain must be associated with vSAN enabled network pool.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `fqdn` (String) Fully qualified domain name of ESXi host
- `password` (String, Sensitive) Password to authenticate to the ESXi host
- `storage_type` (String) Storage Type. One among: VSAN, VSAN_ESA, VSAN_REMOTE, NFS, VMFS_FC, VVOL
- `username` (String) Username to authenticate to the ESXi host

### Optional

- `network_pool_id` (String) ID of the network pool to associate the ESXi host with
- `network_pool_name` (String) Name of the network pool to associate the ESXi host with
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String) Assignable status of the host.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
