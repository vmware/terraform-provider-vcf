/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package datastores

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
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
		},
	}
}

func TryConvertToVsanDatastoreSpec(object map[string]interface{}) (*models.VSANDatastoreSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VSANDatastoreSpec, object is nil")
	}
	datastoreName := object["datastore_name"].(string)
	if len(datastoreName) == 0 {
		return nil, fmt.Errorf("cannot convert to VSANDatastoreSpec, datastore_name is required")
	}
	result := &models.VSANDatastoreSpec{}
	result.DatastoreName = &datastoreName
	licenseKey := object["license_key"].(string)
	result.LicenseKey = licenseKey
	if dedupAndCompressionEnabled, ok := object["dedup_and_compression_enabled"]; ok && !validation_utils.IsEmpty(dedupAndCompressionEnabled) {
		result.DedupAndCompressionEnabled = dedupAndCompressionEnabled.(bool)
	}
	if failuresToTolerate, ok := object["failures_to_tolerate"]; ok && !validation_utils.IsEmpty(failuresToTolerate) {
		failuresToTolerateInt := int32(failuresToTolerate.(int))
		result.FailuresToTolerate = &failuresToTolerateInt
	}

	return result, nil
}
