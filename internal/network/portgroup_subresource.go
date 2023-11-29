// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
	"strings"
)

// PortgroupSchema this helper function extracts the Portgroup Schema, so that
// it's made available for both workload domain and cluster creation.
func PortgroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Port group name",
				ValidateFunc: validation.NoZeroValues,
			},
			"transport_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Port group transport type, One among: VSAN, VMOTION, MANAGEMENT, PUBLIC, " +
					"NFS, VREALIZE, ISCSI, EDGE_INFRA_OVERLAY_UPLINK",
				ValidateFunc: validation.StringInSlice([]string{
					"VSAN", "VMOTION", "MANAGEMENT", "PUBLIC", "NFS", "VREALIZE", "ISCSI", "EDGE_INFRA_OVERLAY_UPLINK",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"active_uplinks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of active uplinks associated with portgroup. This is only supported for VxRail.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func tryConvertToPortgroupSpec(object map[string]interface{}) (*models.PortgroupSpec, error) {
	result := &models.PortgroupSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to PortgroupSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to PortgroupSpec, name is required")
	}
	result.Name = &name
	if transportType, ok := object["transport_type"]; ok && !validationutils.IsEmpty(transportType) {
		transportTypeString := transportType.(string)
		result.TransportType = &transportTypeString
	}
	if activeUplinks, ok := object["active_uplinks"].([]string); ok && !validationutils.IsEmpty(activeUplinks) {
		result.ActiveUplinks = []string{}
		result.ActiveUplinks = append(result.ActiveUplinks, activeUplinks...)
	}

	return result, nil
}

func flattenPortgroupSpec(spec *models.PortgroupSpec) map[string]interface{} {
	result := make(map[string]interface{})
	if spec == nil {
		return result
	}
	result["name"] = *spec.Name
	result["transport_type"] = *spec.TransportType
	result["active_uplinks"] = spec.ActiveUplinks

	return result
}
