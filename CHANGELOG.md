# CHANGELOG

## [v0.12.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.12.0)

> Release Date: 2024-11-18

FEATURES:

- Added support for VCF 5.2.0 [#246](https://github.com/vmware/terraform-provider-vcf/pull/246)

CHORES:

- Updated `actions/setup-go` to 5.1.0 [#249](https://github.com/vmware/terraform-provider-vcf/pull/249)
- Updated `crazy-max/ghaction-import-gpg` to 6.2.0 [#250](https://github.com/vmware/terraform-provider-vcf/pull/250)
- Updated `terraform-plugin-framework-validators` to 0.14.0 [#248](https://github.com/vmware/terraform-provider-vcf/pull/248)
- Updated `terraform-plugin-mux` to 0.17.0 [#255](https://github.com/vmware/terraform-provider-vcf/pull/255)
- Updated `terraform-plugin-sdk/v2` to 2.35.0 [#260](https://github.com/vmware/terraform-provider-vcf/pull/260)
- Updated `goreleaser-action` to 6.1.0 [#262](https://github.com/vmware/terraform-provider-vcf/pull/262)
- Updated `terraform-plugin-docs` to 0.20.0 [#261](https://github.com/vmware/terraform-provider-vcf/pull/261)
- Updated `terraform-plugin-framework` to 1.13.0 [#258](https://github.com/vmware/terraform-provider-vcf/pull/258)
- Updated `golangci-lint` configuration [#247](https://github.com/vmware/terraform-provider-vcf/pull/247)
- Updated codeowners configuration [#257](https://github.com/vmware/terraform-provider-vcf/pull/257)

## [v0.11.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.11.0)

> Release Date: 2024-10-07

BUG FIXES:

- `resource/vcf_credentials`: Fixed the missing `resourceType` definition for `VROPS`.
- `data/vcf_credentials`: Fixed the `resourceType` definition for NSX Edges from `NSX_EDGE` to `NSXT_EDGE` per the API.
- `data/vcf_credentials`: Fixed a segmentation violation.
- `provider`: Updated the provider configuration not to assume the configuration is for Cloud Builder if the SDDC Manager username is not provided and added a check to ensure that at least one of the configurations is provided. If neither is provided, returns an appropriate error message. [\#239](https://github.com/vmware/terraform-provider-vcf/pull/239)

FEATURES:

- `data/vcf_network_pool`: Added network pool data source. [\#225](https://github.com/vmware/terraform-provider-vcf/pull/225)
- `data/vcf_domain`: Updated to support `name` alongside the existing `domain_id`. [\#228](https://github.com/vmware/terraform-provider-vcf/pull/228)
- `provider`: Added a task tracker to log the messages for each subtask. When integrated into a resource logging can be enabled with `TF_LOG_PROVIDER_VCF` with valid a log level. [\#227](https://github.com/vmware/terraform-provider-vcf/pull/227)

REFACTOR:

- Refactored instances of `apiClient` to be more concise, where applicable. This is preferred in Go for its brevity and clarity. [\#231](https://github.com/vmware/terraform-provider-vcf/pull/231)

CHORES:

- Added CodeQL Analysis. [\#221](https://github.com/vmware/terraform-provider-vcf/pull/221)
- Updated Go to v1.22.6 [\#221](https://github.com/vmware/terraform-provider-vcf/pull/221)
- Updated `vmware/vcf-sdk-go` to 0.3.3. [\#203](https://github.com/vmware/terraform-provider-vcf/pull/203)
- Updated `hashicorp/terraform-plugin-framework` to 1.13.0. [\#221](https://github.com/vmware/terraform-provider-vcf/pull/221), [\#240](https://github.com/vmware/terraform-provider-vcf/pull/240)
- Updated `hashicorp/terraform-plugin-framework-validators` to 0.13.0 [\#220](https://github.com/vmware/terraform-provider-vcf/pull/220)
- Updated `hashicorp/terraform-plugin-go` to 0.24.0. [\#241](https://github.com/vmware/terraform-provider-vcf/pull/241)
- Updated `golangci/golangci-lint-action` to 6.1.0. [\#205](https://github.com/vmware/terraform-provider-vcf/pull/205)
- Addressed linting errors identified by `golangci` to 6.1.0. [\#222](https://github.com/vmware/terraform-provider-vcf/pull/222)
- Updated to uses Go's idiomatic conventions group imports. [\#223](https://github.com/vmware/terraform-provider-vcf/pull/223)
- Updated and/or added the copyright and SPDX-License-Identifier, as needed. [\#229](https://github.com/vmware/terraform-provider-vcf/pull/229)
- Updated the provider parameter descriptions; otherwise `make documentation` failed. [\#233](https://github.com/vmware/terraform-provider-vcf/pull/233)

## [v0.10.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.10.0)

> Release Date: 2024-07-09

BREAKING CHANGES:

- The identifier for `r/csr` has changed. Existing resources will become invalid.

FEATURES:

- Allow deployment of VUM-based management domain. [\#151](https://github.com/vmware/terraform-provider-vcf/issues/151)
- Add support for vSAN ESA enablement. [\#182](https://github.com/vmware/terraform-provider-vcf/issues/182)
- Accept names instead of identifiers for several resources. [\#91](https://github.com/vmware/terraform-provider-vcf/issues/91), [\#191](https://github.com/vmware/terraform-provider-vcf/issues/191), [\#193](https://github.com/vmware/terraform-provider-vcf/issues/193)

BUG FIXES:

- `r/csr` fails for `SDDC_MANAGER` and `NSXT_MANAGER`. [\#195](https://github.com/vmware/terraform-provider-vcf/issues/195)

## [v0.9.1](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.9.1)

> Release Date: 2024-06-12

BUG FIXES:

- Tier 0 and Tier 1 routers are now optional for Edge Clusters. [\#177](https://github.com/vmware/terraform-provider-vcf/issues/177)
- Accept VLAN "0" for network pools. [\#175](https://github.com/vmware/terraform-provider-vcf/issues/175)
- New properties for management network configuration on edge nodes. [\#147](https://github.com/vmware/terraform-provider-vcf/issues/147)

## [v0.9.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.9.0)

> Release Date: 2024-05-23

FEATURES:

- Official support for VCF 5.1.1. [\#173](https://github.com/vmware/terraform-provider-vcf/pull/173)

## [v0.8.5](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.5)

> Release Date: April 26, 2024-04-26

FEATURES:

- New resource for exporting cluster personality. [\#143](https://github.com/vmware/terraform-provider-vcf/pull/143)
- Support configuring vSAN in stretched mode. [\#154](https://github.com/vmware/terraform-provider-vcf/pull/154)

BUG FIXES:

- Fix cluster creation with vLCM image. [\#148](https://github.com/vmware/terraform-provider-vcf/pull/148)
- Remove BGP Peer password requirements. [\#150](https://github.com/vmware/terraform-provider-vcf/pull/150)

## [v0.8.1](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.1)

> Release Date: 2024-02-06

BUG FIXES:

- Respect static IP pool configuration when configuring NSX. [\#113](https://github.com/vmware/terraform-provider-vcf/issues/113)
- Fix Edge ASN upper boundary on 32-bit systems. [\#120](https://github.com/vmware/terraform-provider-vcf/issues/120)

## [v0.8.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.8.0)

> Release Date: 2024-01-31

FEATURES:

- NSX Edge Cluster

## [v0.7.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.7.0)

> Release Date: Jan 12 2024

FEATURES:

- Credentials data source.
- Password update.
- Password rotation.
- Password auto-rotate policy configuration.

## [v0.6.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.6.0)

> Release Date: 2023-11-23

FEATURES:

- Support for CA configuration.
- Support for CSR generation.
- Support for replacing certificate of a resource in a Domain via configured CA.
- Support for replacing certificate of a resource in a Domain via external CA.

## [v0.5.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.5.0)

> Release Date: 2023-10-09

FEATURES:

- Add support for management domain deployment (bring-up). [\#38](https://github.com/vmware/terraform-provider-vcf/issues/38)

**Note:** Provider has two mutually exclusive modes of operation: CloudBuilder mode (for bring-up) and SDDC Manager mode.

## [v0.4.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.4.0)

> Release Date: 2023-09-11

BREAKING CHANGES:

- Removed the attribute "host_id" from the "vcf_host" resource and replaced it with just "id" as per Terraform standard practice. This way users can refer to the Host UUID (in cluster host spec for example) in the standard way, e.g. "vcf_host.host1.id".
- Replaced attribute "nsx_cluster_ref" from the "vcf_domain" data source with a richer "nsx_configuration", that additionally contains IPs, Names and DNS Names of NSX-T Manager Nodes.
- Renamed attribute "nsx_configuration.nsx_manager_node.dns_name" in "vcf_domain" to "nsx_configuration.nsx_manager_node.fqdn" for clarity.
- Renamed attribute "vcenter" to "vcenter_configuration" in "vcf_domain" resource and "vcf_domain" data source.
- Replaced attribute "dns_name" in "vcenter_configuration" in "vcf_domain" resource with "fqdn".
- Replaced attribute "vcenter_fqdn" and "vcenter_id from the "vcf_domain" data source with "vcenter_configuration" subresource, that contains "id" and "fqdn" attributes.

FEATURES:

- Extend support for host resource: import. [\#36](https://github.com/vmware/terraform-provider-vcf/issues/36)
- Add support for workload domain resource: import. [\#35](https://github.com/vmware/terraform-provider-vcf/issues/35)
- Add support for configuration of NSX host TEP pool (static / DHCP) in `r/vcf_domain`. [\#54](https://github.com/vmware/terraform-provider-vcf/issues/54)

**Note:** Management domain cannot be imported, but can be used as data source.

BUG FIXES:

- Include "domain_id" attribute to both imported cluster and cluster data source. [\#49](https://github.com/vmware/terraform-provider-vcf/issues/49)

## [v0.3.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.3.0)

> Release Date: 2023-08-22

FEATURES:

- Add support for workload domain cluster resource: read, add, update, delete. [\#32](https://github.com/vmware/terraform-provider-vcf/issues/32)
- Add support for workload domain cluster data source. [\#32](https://github.com/vmware/terraform-provider-vcf/issues/34)
- Extend support for workload domain cluster resource: import. [\#33](https://github.com/vmware/terraform-provider-vcf/issues/33)
- Extend support for workload domain cluster: expand and contract. [\#37](https://github.com/vmware/terraform-provider-vcf/issues/37)

BUG FIXES:

- Fix `IsEmpty` not checking for boolean. [\#45](https://github.com/vmware/terraform-provider-vcf/pull/45)

## [v0.2.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.2.0)

> Release Date: 2023-07-25

Add support for creating/deleting workload domains and being used as data source.

## [v0.1.0](https://github.com/vmware/terraform-provider-vcf/releases/tag/v0.1.0)

> Release Date: 2023-06-05

Initial release, adding support for commissioning/decommissioning hosts, creating/destroying network
pools, creating/destroying SSO user, turning on/off the telemetry (CEIP).
