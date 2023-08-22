# CHANGELOG

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
