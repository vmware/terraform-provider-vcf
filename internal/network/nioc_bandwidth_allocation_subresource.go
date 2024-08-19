// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"

	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// NiocBandwidthAllocationSchema this helper function extracts the NiocBandwidthAllocation
// Schema, so that it's made available for both Domain and Cluster creation.
func NiocBandwidthAllocationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Host infrastructure traffic type. " +
					"Example: management, faultTolerance, vmotion, " +
					"virtualMachine, iSCSI, nfs, hbr, vsan, vdp etc.",
				ValidateFunc: validation.NoZeroValues,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The maximum allowed usage for a traffic class belonging to this resource pool per host " +
					"physical NIC. The utilization of a traffic class will not exceed the specified limit even if " +
					"there are available network resources. If this value is unset or set to -1 in an update " +
					"operation, then there is no limit on the network resource usage (only bounded by available " +
					"resource and shares). Units are in Mbits/sec",
				ValidateFunc: validation.NoZeroValues,
			},
			"reservation": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Amount of bandwidth resource that is guaranteed available to the host infrastructure traffic " +
					"class. If the utilization is less than the reservation, the extra bandwidth is used for other " +
					"host infrastructure traffic class types. Unit is Mbits/sec",
				ValidateFunc: validation.NoZeroValues,
			},
			"shares": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The number of shares allocated. Used to determine resource allocation in case of resource " +
					"contention. This value is only set if level is set to custom. If level is not set to custom, " +
					"this value is ignored. Therefore, only shares with custom values can be compared. " +
					"There is no unit for this value. It is a relative measure based on the settings for other resource pools.",
				ValidateFunc: validation.NoZeroValues,
			},
			"shares_level": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The allocation level. The level is a simplified view of shares. Levels map to a " +
					"pre-determined set of numeric values for shares. If the shares value does not map to a " +
					"predefined size, then the level is set as custom. One among: low, normal, high, custom",
				ValidateFunc: validation.StringInSlice([]string{
					"low", "normal", "high", "custom",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
		},
	}
}

func tryConvertToNiocBandwidthAllocationSpec(object map[string]interface{}) (*models.NiocBandwidthAllocationSpec, error) {
	result := &models.NiocBandwidthAllocationSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to NiocBandwidthAllocationSpec, object is nil")
	}
	typeParam := object["type"].(string)
	if len(typeParam) == 0 {
		return nil, fmt.Errorf("cannot convert to NiocBandwidthAllocationSpec, type is required")
	}
	result.Type = &typeParam
	result.NiocTrafficResourceAllocation = &models.NiocTrafficResourceAllocation{}
	if limit, ok := object["limit"]; ok && !validationutils.IsEmpty(limit) {
		limitRef := limit.(int64)
		result.NiocTrafficResourceAllocation.Limit = &limitRef
	}
	if reservation, ok := object["reservation"]; ok && !validationutils.IsEmpty(reservation) {
		reservationRef := reservation.(int64)
		result.NiocTrafficResourceAllocation.Reservation = &reservationRef
	}
	if shares, ok := object["shares"]; ok && !validationutils.IsEmpty(shares) {
		if result.NiocTrafficResourceAllocation.SharesInfo == nil {
			result.NiocTrafficResourceAllocation.SharesInfo = &models.SharesInfo{}
		}
		result.NiocTrafficResourceAllocation.SharesInfo.Shares = shares.(int32)
	}
	if sharesLevel, ok := object["shares_level"]; ok && !validationutils.IsEmpty(sharesLevel) {
		if result.NiocTrafficResourceAllocation.SharesInfo == nil {
			result.NiocTrafficResourceAllocation.SharesInfo = &models.SharesInfo{}
		}
		result.NiocTrafficResourceAllocation.SharesInfo.Level = sharesLevel.(string)
	}
	return result, nil
}

func flattenNiocBandwidthAllocationSpec(spec *models.NiocBandwidthAllocationSpec) map[string]interface{} {
	result := make(map[string]interface{})
	if spec == nil {
		return result
	}
	result["type"] = *spec.Type
	result["limit"] = *spec.NiocTrafficResourceAllocation.Limit
	if spec.NiocTrafficResourceAllocation != nil {
		result["reservation"] = *spec.NiocTrafficResourceAllocation.Reservation
		result["shares"] = spec.NiocTrafficResourceAllocation.SharesInfo.Shares
		result["shares_level"] = spec.NiocTrafficResourceAllocation.SharesInfo.Level
	}

	return result
}
