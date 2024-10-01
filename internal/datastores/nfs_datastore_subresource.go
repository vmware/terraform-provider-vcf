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

// NfsDatastoreSchema this helper function extracts the NFS Datastore schema, so that
// it's made available for both workload domain and cluster creation.
func NfsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "NFS datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Shared directory path used for NFS based cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"read_only": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Readonly is used to identify whether to mount the directory as readOnly or not",
			},
			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified domain name or IP address of the NFS endpoint",
				ValidateFunc: validation.NoZeroValues,
			},
			"user_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "User tag used to annotate NFS share",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func TryConvertToNfsDatastoreSpec(object map[string]interface{}) (*vcf.NfsDatastoreSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to NfsDatastoreSpec, object is nil")
	}
	datastoreName := object["datastore_name"].(string)
	if len(datastoreName) == 0 {
		return nil, fmt.Errorf("cannot convert to NfsDatastoreSpec, datastore_name is required")
	}
	path := object["path"].(string)
	if len(path) == 0 {
		return nil, fmt.Errorf("cannot convert to NfsDatastoreSpec, path is required")
	}
	result := &vcf.NfsDatastoreSpec{}
	result.DatastoreName = datastoreName
	result.NasVolume = vcf.NasVolumeSpec{}
	result.NasVolume.Path = path
	if readOnly, ok := object["read_only"]; ok && !validationutils.IsEmpty(readOnly) {
		result.NasVolume.ReadOnly = readOnly.(bool)
	} else {
		return nil, fmt.Errorf("cannot convert to NfsDatastoreSpec, read_only is required")
	}
	if serverName, ok := object["server_name"]; ok && !validationutils.IsEmpty(serverName) {
		result.NasVolume.ServerName = []string{}
		result.NasVolume.ServerName = append(result.NasVolume.ServerName, serverName.(string))
	} else {
		return nil, fmt.Errorf("cannot convert to NfsDatastoreSpec, server_name is required")
	}
	if userTag, ok := object["user_tag"]; ok && !validationutils.IsEmpty(userTag) {
		result.NasVolume.UserTag = utils.ToStringPointer(userTag)
	}
	return result, nil
}
