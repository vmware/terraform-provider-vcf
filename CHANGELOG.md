# CHANGELOG

## [v0.9.1](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.9.1)

> Release Date: June 12 2024

BUG FIXES:
* Tier 0 and Tier 1 routers are now optional for Edge Clusters [\#177](https://github.com/vmware/terraform-provider-vcf/issues/177)
* Accept VLAN "0" for network pools [\#175](https://github.com/vmware/terraform-provider-vcf/issues/175)
* New properties for management network configuration on edge nodes [\#147](https://github.com/vmware/terraform-provider-vcf/issues/147)

## [v0.9.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.9.0)

> Release Date: May 23 2024

FEATURES:
* Official support for VCF 5.1.1 [\#173](https://github.com/vmware/terraform-provider-vcf/pull/173)

## [v0.8.5](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.5)

> Release Date: Apr 26 2024

FEATURES:
* New resource for exporting cluster personality [\#143](https://github.com/vmware/terraform-provider-vcf/pull/143)
* Support configuring vSAN in stretched mode [\#154](https://github.com/vmware/terraform-provider-vcf/pull/154)

BUG FIXES:
* Fix cluster creation with vLCM image [\#148](https://github.com/vmware/terraform-provider-vcf/pull/148)
* Remove BGP Peer password requirements [\#150](https://github.com/vmware/terraform-provider-vcf/pull/150)

## [v0.8.1](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.1)

> Release Date: Feb 6 2024

BUG FIXES:
* Respect static IP pool configuration when configuring NSX [\#113](https://github.com/vmware/terraform-provider-vcf/issues/113)
* Fix Edge ASN upper boundary on 32-bit systems [\#120](https://github.com/vmware/terraform-provider-vcf/issues/120)

## [v0.8.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.0)

> Release Date: Jan 31 2024

FEATURES:
* NSX Edge Cluster

## [v0.7.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.7.0)

> Release Date: Jan 12 2024

FEATURES:
* Credentials data source
* Password update
* Password rotation
* Password auto-rotate policy configuration

## [v0.6.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.6.0)

> Release Date: Nov 23 2023

FEATURES:
* Support for CA configuration
* Support for CSR generation
* Support for replacing certificate of a resource in a Domain via configured CA
* Support for replacing certificate of a resource in a Domain via external CA

## [v0.5.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.5.0)

> Release Date: Oct 9th 2023

FEATURES:
* Add support for management domain deployment (bringup) [\#38](https://github.com/vmware/terraform-provider-vcf/issues/38)

**Note:** Provider has two mutually exclusive modes of operation: CloudBuilder mode (for bringup) and SDDCManager mode

## [v0.4.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.4.0)

> Release Date: Sep 11th 2023

BREAKING CHANGES:

* Removed the attribute "host_id" from the "vcf_host" resource and replaced it with just "id" as per Terraform standard practice. This way users can refer to the Host UUID (in cluster host spec for example) in the standard way, e.g. "vcf_host.host1.id"
* Replaced attribute "nsx_cluster_ref" from the "vcf_domain" datasource with a richer "nsx_configuration", that additionally contains IPs, Names and DNS Names of NSX-T Manager Nodes
* Renamed attribute "nsx_configuration.nsx_manager_node.dns_name" in "vcf_domain" to "nsx_configuration.nsx_manager_node.fqdn" for clarity
* Renamed attribute "vcenter" to "vcenter_configuration" in "vcf_domain" resource and "vcf_domain" datasource
* Replaced attribute "dns_name" in "vcenter_configuration" in "vcf_domain" resource with "fqdn"
* Replaced attribute "vcenter_fqdn" and "vcenter_id from the "vcf_domain" datasource with "vcenter_configuration" subresource, that contains "id" and "fqdn" attributes. 

FEATURES:
* Extend support for host resource: import [\#36](https://github.com/vmware/terraform-provider-vcf/issues/36)
* Add support for workload domain resource: import [\#35](https://github.com/vmware/terraform-provider-vcf/issues/35)
* Add support for configuration of NSX host TEP pool (static / DHCP) in r/vcf_domain [\#54](https://github.com/vmware/terraform-provider-vcf/issues/54)

**Note:** Management domain cannot be imported, but can be used as datasource

BUG FIXES:
* Include "domain_id" attribute to both imported cluster and cluster datasource [\#49](https://github.com/vmware/terraform-provider-vcf/issues/49)

## [v0.3.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.3.0)

> Release Date: Aug 22nd 2023

FEATURES:
* Add support for workload domain cluster resource: read, add, update, delete [\#32](https://github.com/vmware/terraform-provider-vcf/issues/32)
* Add support for workload domain cluster data source [\#32](https://github.com/vmware/terraform-provider-vcf/issues/34)
* Extend support for workload domain cluster resource: import [\#33](https://github.com/vmware/terraform-provider-vcf/issues/33)
* Extend support for workload domain cluster: expand and contract [\#37](https://github.com/vmware/terraform-provider-vcf/issues/37)

BUG FIXES:
* Fix IsEmpty not checking for boolean [\#45](https://github.com/vmware/terraform-provider-vcf/pull/45)



## [v0.2.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.2.0)

> Release Date: Jul 25th 2023

Add support for creating/deleting workload domains and being used as datasource.

## [v0.1.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.1.0)

> Release Date: Jun 5th 2023

Initial release, adding support for commissioning/decommissioning hosts, creating/destroying
network pools, creating/destroying SSO user, turning on/off the telemetry (CEIP).
