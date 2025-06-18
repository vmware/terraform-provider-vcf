// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/installer"
)

var sharesLevelValues = []string{"custom", "high", "low", "normal"}
var resourcePoolTypeValues = []string{"management", "compute", "network"}

func GetSddcClusterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"datacenter_name": {
					Type:        schema.TypeString,
					Description: "vCenter Datacenter Name",
					Required:    true,
				},
				"cluster_name": {
					Type:        schema.TypeString,
					Description: "vCenter Cluster Name",
					Required:    true,
				},
				"cluster_evc_mode": {
					Type:        schema.TypeString,
					Description: "vCenter cluster EVC mode",
					Optional:    true,
				},
				"resource_pool": getResourcePoolSchema(),
			},
		},
	}
}

func getResourcePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "Resource Pool name",
					Required:    true,
				},
				"type": {
					Type:         schema.TypeString,
					Description:  "Type of resource pool, possible values: \"management\", \"compute\", \"network\"",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(resourcePoolTypeValues, false),
				},
				"cpu_limit": {
					Type:         schema.TypeFloat, // There is no int64 type in the Terraform, so using Float as the next best thing
					Description:  "CPU limit, default -1 (unlimited)",
					Optional:     true,
					Default:      -1,
					ValidateFunc: validation_utils.ValidateParsingFloatToInt,
				},
				"cpu_reservation_expandable": {
					Type:        schema.TypeBool,
					Description: "Is CPU reservation expandable, default true",
					Optional:    true,
					Default:     true,
				},
				"cpu_reservation_mhz": {
					Type:         schema.TypeFloat,
					Description:  "CPU reservation in Mhz",
					Optional:     true,
					ValidateFunc: validation_utils.ValidateParsingFloatToInt,
				},
				"cpu_reservation_percentage": {
					Type:         schema.TypeInt,
					Description:  "CPU reservation percentage, from 0 to 100, default 0",
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 100),
				},
				"cpu_shares_level": {
					Type:         schema.TypeString,
					Description:  "CPU shares level, default 'normal', possible values: \"custom\", \"high\", \"low\", \"normal\"",
					Optional:     true,
					Default:      "normal",
					ValidateFunc: validation.StringInSlice(sharesLevelValues, false),
				},
				"cpu_shares_value": {
					Type:        schema.TypeInt,
					Description: "CPU shares value, only required when shares level is 'normal'",
					Optional:    true,
					Default:     0,
				},
				"memory_limit": {
					Type:         schema.TypeFloat,
					Description:  "Memory limit, default -1 (unlimited)",
					Optional:     true,
					Default:      -1,
					ValidateFunc: validation_utils.ValidateParsingFloatToInt,
				},
				"memory_reservation_expandable": {
					Type:        schema.TypeBool,
					Description: "Is Memory reservation expandable, default true",
					Optional:    true,
					Default:     true,
				},
				"memory_reservation_mb": {
					Type:         schema.TypeFloat,
					Description:  "Memory reservation in MB",
					Optional:     true,
					ValidateFunc: validation_utils.ValidateParsingFloatToInt,
				},
				"memory_reservation_percentage": {
					Type:         schema.TypeInt,
					Description:  "Memory reservation percentage, from 0 to 100, default 0",
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 100),
				},
				"memory_shares_level": {
					Type:         schema.TypeString,
					Description:  "Memory shares level, default 'normal', possible values: \"custom\", \"high\", \"low\", \"normal\"",
					Optional:     true,
					Default:      "normal",
					ValidateFunc: validation.StringInSlice(sharesLevelValues, false),
				},
				"memory_shares_value": {
					Type:        schema.TypeInt,
					Description: "Memory shares value, only required when shares level is 'normal'",
					Optional:    true,
					Default:     0,
				},
			},
		},
	}
}

func GetSddcClusterSpecFromSchema(rawData []interface{}) *installer.SddcClusterSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	clusterName := utils.ToStringPointer(data["cluster_name"])
	datacenterName := utils.ToStringPointer(data["datacenter_name"])
	clusterEvcMode := data["cluster_evc_mode"].(string)

	clusterSpecBinding := &installer.SddcClusterSpec{
		ClusterEvcMode: &clusterEvcMode,
		ClusterName:    clusterName,
		DatacenterName: datacenterName,
	}

	if resourcePoolSpecs := getResourcePoolSpecsFromSchema(
		data["resource_pool"].([]interface{})); len(resourcePoolSpecs) > 0 {
		clusterSpecBinding.ResourcePoolSpecs = &resourcePoolSpecs
	}

	return clusterSpecBinding
}

func getResourcePoolSpecsFromSchema(rawData []interface{}) []installer.ResourcePoolSpec {
	var resourcePoolSpecs []installer.ResourcePoolSpec
	for _, resourcePool := range rawData {
		data := resourcePool.(map[string]interface{})
		cpuLimit := int64(data["cpu_limit"].(float64))
		cpuReservationExpandable := data["cpu_reservation_expandable"].(bool)
		cpuReservationMhz := int64(data["cpu_reservation_mhz"].(float64))
		cpuReservationPercentage := utils.ToInt32Pointer(data["cpu_reservation_percentage"])
		cpuSharesLevel := installer.ResourcePoolSpecCpuSharesLevel(data["cpu_shares_level"].(string))
		cpuSharesValue := int32(data["cpu_shares_value"].(int))
		memoryLimit := int64(data["memory_limit"].(float64))
		memoryReservationPercentage := utils.ToInt32Pointer(data["memory_reservation_percentage"])
		memoryReservationExpandable := utils.ToBoolPointer(data["memory_reservation_expandable"])
		memoryReservationMB := int64(data["memory_reservation_mb"].(float64))
		memorySharesLevel := installer.ResourcePoolSpecMemorySharesLevel(data["memory_shares_level"].(string))
		memorySharesValue := int32(data["memory_shares_value"].(int))
		name := utils.ToStringPointer(data["name"])
		resourcePoolType := installer.ResourcePoolSpecType(data["type"].(string))

		resourcePoolSpec := &installer.ResourcePoolSpec{
			CpuLimit:                    &cpuLimit,
			CpuReservationExpandable:    &cpuReservationExpandable,
			CpuReservationMhz:           &cpuReservationMhz,
			CpuReservationPercentage:    cpuReservationPercentage,
			CpuSharesValue:              &cpuSharesValue,
			CpuSharesLevel:              &cpuSharesLevel,
			MemoryLimit:                 &memoryLimit,
			MemorySharesLevel:           &memorySharesLevel,
			MemoryReservationPercentage: memoryReservationPercentage,
			MemoryReservationExpandable: memoryReservationExpandable,
			MemoryReservationMb:         &memoryReservationMB,
			MemorySharesValue:           &memorySharesValue,
			Name:                        name,
			Type:                        &resourcePoolType,
		}
		resourcePoolSpecs = append(resourcePoolSpecs, *resourcePoolSpec)
	}
	return resourcePoolSpecs
}
