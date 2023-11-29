// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package datastores

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"
)

// VmfsDatastoreSchema this helper function extracts the VMFS Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func VmfsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_names": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "VMFS datastore names used for VMFS on FC for cluster creation",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func TryConvertToVmfsDatastoreSpec(object map[string]interface{}) (*models.VmfsDatastoreSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VmfsDatastoreSpec, object is nil")
	}
	datastoreNames := object["datastore_names"].([]string)
	if len(datastoreNames) == 0 {
		return nil, fmt.Errorf("cannot convert to VmfsDatastoreSpec, datastore_names is required")
	}
	result := &models.VmfsDatastoreSpec{}
	result.FcSpec = []*models.FcSpec{}
	for _, datastoreName := range datastoreNames {
		datastoreNameRef := &datastoreName
		result.FcSpec = append(result.FcSpec, &models.FcSpec{DatastoreName: datastoreNameRef})
	}
	return result, nil
}
