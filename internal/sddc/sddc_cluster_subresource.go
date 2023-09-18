/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/models"
)

var sharesLevelValues = []string{"custom", "high", "low", "normal"}
var resourcePoolTypeValues = []string{"management", "compute", "network"}

// var vmFolders = []string{"MANAGEMENT", "NETWORKING", "EDGENODES"}.

func GetSddcClusterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cluster_evc_mode": {
					Type:        schema.TypeString,
					Description: "vCenter cluster EVC mode",
					Optional:    true,
				},
				"cluster_name": {
					Type:        schema.TypeString,
					Description: "vCenter Cluster Name",
					Required:    true,
				},
				"host_failures_to_tolerate": {
					Type:         schema.TypeInt,
					Description:  "Host failures to tolerate. In between 0 and 3",
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 3),
				},
				"host_profile_compliance_check_hour": {
					Type:         schema.TypeInt,
					Description:  "Hour at which the scheduled compliance check runs. In between 0 and 23",
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 23),
				},
				"host_profile_compliance_check_minute": {
					Type:         schema.TypeInt,
					Description:  "Minute at which the scheduled compliance check run. In between 0 and 59",
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 59),
				},
				"host_id": {
					Type:        schema.TypeList,
					Description: "vCenter Cluster Host IDs",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"personality_name": {
					Type:        schema.TypeString,
					Description: "Cluster Personality Name",
					Optional:    true,
				},
				"resource_pool": getResourcePoolSchema(),
				// TODO Implement VM Folders
				"vm_folders": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
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
				"cpu_limit": {
					Type:        schema.TypeInt,
					Description: "CPU limit, default -1 (unlimited)",
					Optional:    true,
					Default:     -1,
				},
				"cpu_reservation_expandable": {
					Type:        schema.TypeBool,
					Description: "Is CPU reservation expandable, default true",
					Optional:    true,
					Default:     true,
				},
				"cpu_reservation_mhz": {
					Type:        schema.TypeInt,
					Description: "CPU reservation in Mhz",
					Optional:    true,
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
					Type:        schema.TypeInt,
					Description: "Memory limit, default -1 (unlimited)",
					Optional:    true,
					Default:     -1,
				},
				"memory_reservation_expandable": {
					Type:        schema.TypeBool,
					Description: "Is Memory reservation expandable, default true",
					Optional:    true,
					Default:     true,
				},
				"memory_reservation_mb": {
					Type:        schema.TypeInt,
					Description: "Memory reservation in MB",
					Optional:    true,
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
			},
		},
	}
}

func GetSddcClusterSpecFromSchema(rawData []interface{}) *models.SDDCClusterSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	clusterEvcMode := data["cluster_evc_mode"].(string)
	clusterName := utils.ToStringPointer(data["cluster_name"])
	hostFailuresToTolerate := utils.ToInt32Pointer(data["host_failures_to_tolerate"])
	hostProfileComplianceCheckHour := utils.ToInt32Pointer(data["host_profile_compliance_check_hour"])
	hostProfileComplianceCheckMinute := utils.ToInt32Pointer(data["host_profile_compliance_check_minute"])
	personalityName := data["personality_name"].(string)

	clusterSpecBinding := &models.SDDCClusterSpec{
		ClusterEvcMode:                   clusterEvcMode,
		ClusterName:                      clusterName,
		HostFailuresToTolerate:           hostFailuresToTolerate,
		HostProfileComplianceCheckHour:   hostProfileComplianceCheckHour,
		HostProfileComplianceCheckMinute: hostProfileComplianceCheckMinute,
		PersonalityName:                  personalityName,
	}

	if hosts, ok := data["host_id"].([]interface{}); ok {
		clusterSpecBinding.Hosts = utils.ToStringSlice(hosts)
	}
	if resourcePoolSpecs := getResourcePoolSpecsFromSchema(
		data["resource_pool"].([]interface{})); len(resourcePoolSpecs) > 0 {
		clusterSpecBinding.ResourcePoolSpecs = resourcePoolSpecs
	}
	return clusterSpecBinding
}

func getResourcePoolSpecsFromSchema(rawData []interface{}) []*models.ResourcePoolSpec {
	var resourcePoolSpecs []*models.ResourcePoolSpec
	for _, resourcePool := range rawData {
		data := resourcePool.(map[string]interface{})
		cpuLimit := data["cpu_limit"].(int64)
		cpuReservationExpandable := data["cpu_reservation_expandable"].(bool)
		cpuReservationMhz := data["cpu_reservation_mhz"].(int64)
		cpuReservationPercentage := utils.ToInt32Pointer(data["cpu_reservation_percentage"])
		cpuSharesLevel := data["cpu_shares_level"].(string)
		cpuSharesValue := data["cpu_shares_value"].(int32)
		memoryLimit := data["memory_limit"].(int64)
		memoryReservationPercentage := utils.ToInt32Pointer(data["memory_reservation_percentage"])
		memoryReservationExpandable := utils.ToBoolPointer(data["memory_reservation_expandable"])
		memoryReservationMB := data["memory_reservation_mb"].(int64)
		memorySharesLevel := data["memory_shares_level"].(string)
		memorySharesValue := data["memory_shares_value"].(int32)
		name := utils.ToStringPointer(data["name"])
		resourcePoolType := data["type"].(string)

		resourcePoolSpec := &models.ResourcePoolSpec{
			CPULimit:                    cpuLimit,
			CPUReservationExpandable:    cpuReservationExpandable,
			CPUReservationMhz:           cpuReservationMhz,
			CPUReservationPercentage:    cpuReservationPercentage,
			CPUSharesValue:              cpuSharesValue,
			CPUSharesLevel:              cpuSharesLevel,
			MemoryLimit:                 memoryLimit,
			MemorySharesLevel:           memorySharesLevel,
			MemoryReservationPercentage: memoryReservationPercentage,
			MemoryReservationExpandable: memoryReservationExpandable,
			MemoryReservationMb:         memoryReservationMB,
			MemorySharesValue:           memorySharesValue,
			Name:                        name,
			Type:                        resourcePoolType,
		}
		resourcePoolSpecs = append(resourcePoolSpecs, resourcePoolSpec)
	}
	return resourcePoolSpecs
}
