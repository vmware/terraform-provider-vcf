// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/network"
)

func DataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of the Cluster to be used as data source",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the domain",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a workload domain that the cluster belongs to",
			},
			"host": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ESXi host information present in the Cluster",
				Elem:        cluster.HostSpecSchema(),
			},
			"vds": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "vSphere Distributed Switches to add to the Cluster",
				Elem:        network.VdsSchema(),
			},
			"primary_datastore_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the primary datastore",
			},
			"primary_datastore_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage type of the primary datastore",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of the cluster if default or not",
			},
			"is_stretched": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of the cluster if stretched or not",
			},
		},
	}
}

func dataSourceClusterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient
	clusterId := data.Get("cluster_id").(string)
	_, err := cluster.ImportCluster(ctx, data, apiClient, clusterId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
