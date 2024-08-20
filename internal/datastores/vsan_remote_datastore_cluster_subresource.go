// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package datastores

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"
)

// VsanRemoteDatastoreClusterSchema this helper function extracts the VSAN Datastore Cluster
// schema, so that it's made available for both Domain and Cluster creation.
func VsanRemoteDatastoreClusterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_uuids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "vSAN HCI Mesh remote datastore UUIDs",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func TryConvertToVSANRemoteDatastoreClusterSpec(object map[string]interface{}) (*models.VSANRemoteDatastoreClusterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VSANRemoteDatastoreClusterSpec, object is nil")
	}
	datastoreUuids := object["datastore_uuids"].([]string)
	if len(datastoreUuids) == 0 {
		return nil, fmt.Errorf("cannot convert to VSANRemoteDatastoreClusterSpec, datastore_uuids is required")
	}
	result := &models.VSANRemoteDatastoreClusterSpec{}
	result.VSANRemoteDatastoreSpec = []*models.VSANRemoteDatastoreSpec{}
	for _, datastoreUuid := range datastoreUuids {
		datastoreUuidRef := &datastoreUuid
		result.VSANRemoteDatastoreSpec = append(result.VSANRemoteDatastoreSpec,
			&models.VSANRemoteDatastoreSpec{DatastoreUUID: datastoreUuidRef})
	}
	return result, nil
}
