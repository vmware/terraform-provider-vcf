/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"
	"sort"
	"time"
)

func DataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of the Domain to be used as data source",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the domain",
			},
			"cluster": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specification representing the clusters in the workload domain",
				Elem:        clusterSubresourceSchema(),
			},
			"nsx_cluster_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Represents NSX Manager cluster references associated with the domain",
				Elem:        network.NsxClusterRefSchema(),
			},
			"vcenter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vCenter Server instance",
			},
			"vcenter_fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fully qualified domain name of the vCenter Server instance",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the workload domain",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the workload domain",
			},
			"sso_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the SSO domain associated with the workload domain",
			},
			"sso_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the SSO domain associated with the workload domain",
			},
			"is_management_sso_domain": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Shows whether the domain is joined to the management domain SSO",
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	getDomainParams := domains.NewGetDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	getDomainParams.ID = data.Get("domain_id").(string)
	domainResult, err := apiClient.Domains.GetDomain(getDomainParams)
	if err != nil {
		return diag.FromErr(err)
	}
	domain := domainResult.Payload

	data.SetId(domain.ID)
	_ = data.Set("name", domain.Name)
	_ = data.Set("status", domain.Status)
	_ = data.Set("type", domain.Type)
	_ = data.Set("sso_id", domain.SSOID)
	_ = data.Set("sso_name", domain.SSOName)
	_ = data.Set("is_management_sso_domain", domain.IsManagementSSODomain)
	if len(domain.VCENTERS) < 1 {
		return diag.FromErr(fmt.Errorf("no vCenter Server instance found for domain %q", data.Id()))
	}
	_ = data.Set("vcenter_id", domain.VCENTERS[0].ID)
	_ = data.Set("vcenter_fqdn", domain.VCENTERS[0].Fqdn)

	err = setClustersDataToDomainDataSource(domain.Clusters, ctx, data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	flattenedNsxClusterRef := make([]map[string]interface{}, 1)
	flattenedNsxClusterRef[0] = *network.FlattenNsxClusterRef(domain.NSXTCluster)
	_ = data.Set("nsx_cluster_ref", flattenedNsxClusterRef)
	return nil
}

func setClustersDataToDomainDataSource(domainClusterRefs []*models.ClusterReference, ctx context.Context, data *schema.ResourceData, apiClient *client.VcfClient) error {
	clusterIds := make([]string, len(domainClusterRefs))
	for i, clusterReference := range domainClusterRefs {
		clusterIds[i] = *clusterReference.ID
	}
	// Sort the id slice, to have a deterministic order in every run of the domain datasource read
	sort.Strings(clusterIds)

	flattenedClusters := make([]map[string]interface{}, len(domainClusterRefs))
	for i, clusterId := range clusterIds {
		getClusterParams := clusters.GetClusterParams{ID: clusterId}
		getClusterParams.WithContext(ctx).WithTimeout(constants.DefaultVcfApiCallTimeout)
		clusterResult, err := apiClient.Clusters.GetCluster(&getClusterParams)
		if err != nil {
			return err
		}
		clusterRef := clusterResult.Payload
		flattenedClusters[i] = *cluster.FlattenCluster(clusterRef)

	}
	_ = data.Set("cluster", flattenedClusters)

	return nil
}
