/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// CommissionedHostSchema this helper function extracts the Host
// schema, so that it's made available for both Domain and Cluster creation.
func CommissionedHostSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a vSphere host in the free pool",
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Host name of the vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
			"availability_zone_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Availability Zone Name. This is required while performing a stretched cluster expand operation",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IP address of the vSphere host",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"license_key": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "License key of a vSphere host in the free pool. This is required except in cases where the " +
					"ESXi host has already been licensed outside of the VMware Cloud Foundation system",
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Username of the vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "SSH password of the vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
			"serial_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Serial Number of the vSphere hosts",
				ValidateFunc: validation.NoZeroValues,
			},
			"ssh_thumbprint": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "SSH thumbprint(fingerprint) of the vSphere host. Note:This field will be mandatory in future releases.",
				ValidateFunc: validation.NoZeroValues,
			},
			"vmnics": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Contains vmnic configurations for vSphere host",
				Elem:        network.VMNicSchema(),
			},
		},
	}
}
