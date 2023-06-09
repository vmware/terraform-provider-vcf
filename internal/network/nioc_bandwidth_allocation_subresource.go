/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NiocBandwidthAllocationSchema this helper function extracts the NiocBandwidthAllocation
// Schema, so that it's made available for both Domain and Cluster creation.
func NiocBandwidthAllocationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DvsHostInfrastructureTrafficResource resource type",
				ValidateFunc: validation.NoZeroValues,
			},
			"nioc_traffic_resource_allocation": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "DvsHostInfrastructureTrafficResourceAllocation",
				MaxItems:    1,
				Elem:        NiocTrafficResourceAllocationSchema(),
			},
		},
	}
}
