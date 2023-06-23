/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// VMNicSchema this helper function extracts the VMNic Schema, so that
// it's made available for both Domain and Cluster creation.
func VMNicSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "VmNic ID of vSphere host to be associated with VDS, once added to cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"move_to_nvds": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This flag determines if the vmnic must be on N-VDS",
			},
			"uplink": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Uplink to be associated with vmnic",
				ValidateFunc: validation.NoZeroValues,
			},
			"vds_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "VDS name to associate with vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func TryConvertToVmNic(object map[string]interface{}) (*models.VMNic, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VMNic, object is nil")
	}
	id := object["id"].(string)
	if len(id) == 0 {
		return nil, fmt.Errorf("cannot convert to VMNic, id is required")
	}
	result := &models.VMNic{}
	result.ID = id
	if moveToNvds, ok := object["move_to_nvds"]; ok && !validation_utils.IsEmpty(moveToNvds) {
		result.MoveToNvds = moveToNvds.(bool)
	}
	if uplink, ok := object["uplink"]; ok && !validation_utils.IsEmpty(uplink) {
		result.Uplink = uplink.(string)
	}
	if vdsName, ok := object["vds_name"]; ok && !validation_utils.IsEmpty(vdsName) {
		result.VdsName = vdsName.(string)
	}
	return result, nil
}

func FlattenVMNic(vmNic *models.VMNicInfo) *map[string]interface{} {
	result := make(map[string]interface{})
	if vmNic == nil {
		return &result
	}

	return &result
}
