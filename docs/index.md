---
page_title: "Terraform Provider for VMware Cloud Foundation"
subcategory: ""
description: |-
---

<img src="https://raw.githubusercontent.com/vmware/terraform-provider-vcf/main/docs/images/icon-color.svg" alt="VMware Cloud Foundation" width="150">

# Terraform Provider for VMware Cloud Foundation

The following table lists the supported platforms for this provider.

| Platform                     | Support     |
|------------------------------|-------------|
| VMware Cloud Foundation 5.1+ | `≥ v0.9.0`  |
| VMware Cloud Foundation 5.0  | `≥ v0.9.0`  |
| VMware Cloud Foundation 4.5  | `≤ v0.8.0`  |
| VMware Cloud Foundation 4.4  | `≤ v0.8.0`  |

[^1]: VMware Cloud Foundation on Dell VxRAIL is **not supported** by this provider.

The plugin supports versions in accordance with the [Broadcom Product Lifecycle][product-lifecycle]. [^1]

See the VMware Cloud Foundation [release notes](https://docs.vmware.com/en/VMware-Cloud-Foundation/) for the individual build numbers.

[product-lifecycle]: https://support.broadcom.com/group/ecx/productlifecycle

## Example Usage

```hcl
terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
      version = "x.y.z"
    }
  }
}

provider "vcf" {
  sddc_manager_host     = var.sddc_manager_host
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  allow_unverified_tls  = var.allow_unverified_tls
}
```

Refer to the provider documentation for information on all of the resources
and data sources supported by this provider. Each includes a detailed
description of the purpose and how to use it.

## Argument Reference

The following arguments are used to configure the provider:

- `sddc_manager_host` - (Optional) Fully qualified domain name or IP address of the SDDC Manager.
- `sddc_manager_password` - (Optional) Password to authenticate to SDDC Manager.
- `sddc_manager_username` - (Optional) Username to authenticate to SDDC Manager.
- `cloud_builder_host` - (Optional) Fully qualified domain name or IP address of the Cloud Builder.
- `cloud_builder_password` - (Optional) Password to authenticate to Cloud Builder.
- `cloud_builder_username` - (Optional) Username to authenticate to Cloud Builder.
- `allow_unverified_tls` (Boolean) If enabled, this allows the use of TLS certificates that cannot be verified.
