// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package datastores

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// VsanDatastoreSchema this helper function extracts the vSAN Datastore schema, so that
// it's made available for both workload domain and cluster creation.
func VsanDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "vSAN datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"license_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "vSAN license key to be used",
				ValidateFunc: validation.NoZeroValues,
			},
			"failures_to_tolerate": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Number of ESXi host failures to tolerate in the vSAN cluster. One of 0, 1, or 2.",
				ValidateFunc: validation.IntBetween(0, 2),
			},
			"dedup_and_compression_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable vSAN deduplication and compression",
			},
			"esa_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable vSAN ESA",
			},
		},
	}
}

func TryConvertToVsanDatastoreSpec(object map[string]interface{}) (*vcf.VsanDatastoreSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VSANDatastoreSpec, object is nil")
	}
	datastoreName := object["datastore_name"].(string)
	if len(datastoreName) == 0 {
		return nil, fmt.Errorf("cannot convert to VSANDatastoreSpec, datastore_name is required")
	}
	result := &vcf.VsanDatastoreSpec{}
	result.DatastoreName = datastoreName
	result.LicenseKey = utils.ToStringPointer(object["license_key"])
	if dedupAndCompressionEnabled, ok := object["dedup_and_compression_enabled"]; ok && !validationutils.IsEmpty(dedupAndCompressionEnabled) {
		result.DedupAndCompressionEnabled = utils.ToBoolPointer(dedupAndCompressionEnabled)
	}
	if esaEnabled, ok := object["esa_enabled"]; ok && !validationutils.IsEmpty(esaEnabled) {
		value := esaEnabled.(bool)
		esaConfig := vcf.EsaConfig{Enabled: value}
		result.EsaConfig = &esaConfig
	}
	if failuresToTolerate, ok := object["failures_to_tolerate"]; ok && !validationutils.IsEmpty(failuresToTolerate) {
		failuresToTolerateInt := int32(failuresToTolerate.(int))
		result.FailuresToTolerate = &failuresToTolerateInt
	}

	return result, nil
}
