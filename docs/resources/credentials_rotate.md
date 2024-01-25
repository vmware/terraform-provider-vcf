---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vcf_credentials_rotate Resource - terraform-provider-vcf"
subcategory: ""
description: |-
  
---

# vcf_credentials_rotate (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credentials` (Block List, Min: 1) The credentials that should be rotated (see [below for nested schema](#nestedblock--credentials))
- `resource_name` (String) The name of the resource which credentials will be rotated
- `resource_type` (String) The type of the resource which credentials will be rotated

### Optional

- `once_only` (Boolean) If set to true operation is executed only once otherwise rotation is done each time.

### Read-Only

- `id` (String) The ID of this resource.
- `last_rotate_time` (String) The time of the last password rotation.

<a id="nestedblock--credentials"></a>
### Nested Schema for `credentials`

Required:

- `credential_type` (String) The type(s) of the account. One among: SSO, SSH, API, FTP, AUDIT
- `user_name` (String) The user name of the account.

Read-Only:

- `password` (String, Sensitive) The password for the account.