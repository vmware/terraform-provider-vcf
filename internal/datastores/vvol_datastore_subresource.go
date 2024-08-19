// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package datastores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// VvolDatastoreSchema this helper function extracts the VVOL Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func VvolDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "vVol datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"storage_container_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage container",
				ValidateFunc: validation.IsUUID,
			},
			"storage_protocol_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the VASA storage protocol. One among: ISCSI, NFS, FC.",
				ValidateFunc: validation.StringInSlice(
					[]string{"ISCSI", "NFS", "FC"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage user",
				ValidateFunc: validation.IsUUID,
			},
			"vasa_provider_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage provider",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func TryConvertToVvolDatastoreSpec(object map[string]interface{}) (*models.VvolDatastoreSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, object is nil")
	}
	result := &models.VvolDatastoreSpec{}
	result.VasaProviderSpec = &models.VasaProviderSpec{}

	datastoreName := object["datastore_name"].(string)
	if len(datastoreName) == 0 {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, datastore_name is required")
	}
	result.Name = &datastoreName

	storageContainerId := object["storage_container_id"].(string)
	if len(storageContainerId) == 0 {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, storage_container_id is required")
	}
	result.VasaProviderSpec.StorageContainerID = &storageContainerId

	storageContainerProtocolType := object["storage_protocol_type"].(string)
	if len(storageContainerProtocolType) == 0 {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, storage_protocol_type is required")
	}
	result.VasaProviderSpec.StorageProtocolType = &storageContainerProtocolType

	userId := object["user_id"].(string)
	if len(userId) == 0 {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, userId is required")
	}
	result.VasaProviderSpec.UserID = &userId

	vasaProviderId := object["vasa_provider_id"].(string)
	if len(vasaProviderId) == 0 {
		return nil, fmt.Errorf("cannot convert to VvolDatastoreSpec, vasa_provider_id is required")
	}
	result.VasaProviderSpec.VasaProviderID = &vasaProviderId

	return result, nil
}
