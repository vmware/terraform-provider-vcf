/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// NsxTSchema this helper function extracts the NSX-T schema, which
// contains the parameters required to install and configure NSX-T in a workload domain.
func NsxTSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Virtual IP address which would act as proxy/alias for NSX-T Managers",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vip_fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "FQDN for VIP so that common SSL certificates can be installed across all managers",
				ValidateFunc: validation.NoZeroValues,
			},
			"license_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "NSX license value",
				ValidateFunc: validation.NoZeroValues,
			},
			"form_factor": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "NSX manager form factor",
				ValidateFunc: validation.NoZeroValues,
			},
			"nsx_manager_admin_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "NSX manager admin password (basic auth and SSH)",
				ValidateFunc: validation.NoZeroValues,
			},
			"nsx_manager_audit_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "NSX manager Audit password",
				ValidateFunc: validation.NoZeroValues,
			},
			"nsx_managers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification details of the NSX Manager virtual machines",
				Elem:        NsxtManagerSchema(),
			},
		},
	}
}
