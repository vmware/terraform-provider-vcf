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
	"github.com/vmware/terraform-provider-vcf/internal/domain"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"github.com/vmware/terraform-provider-vcf/internal/vcenter"
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
			"nsx_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Represents NSX Manager cluster references associated with the domain",
				Elem:        network.NsxSchema(),
			},
			"vcenter_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specification describing vCenter Server instance settings",
				Elem:        vcenter.VCSubresourceSchema(),
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
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient
	domainId := data.Get("domain_id").(string)

	_, err := domain.ImportDomain(ctx, data, apiClient, domainId, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
