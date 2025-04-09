// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

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
				ValidateFunc: validation.NoZeroValues,
				Optional:     true,
				Description:  "The ID of the workload domain.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the workload domain.",
			},
			"cluster": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cluster references associated with the workload domain",
				Elem:        clusterSubresourceSchema(),
			},
			"nsx_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The NSX Manager cluster references associated with the workload domain.",
				Elem:        network.NsxSchema(),
			},
			"vcenter_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The vCenter Server instance references associated with the workload domain.",
				Elem:        vcenter.VCSubresourceSchema(),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the workload domain.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of workload domain.",
			},
			"sso_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the SSO domain associated with the workload domain.",
			},
			"sso_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the SSO domain associated with the workload domain.",
			},
			"is_management_sso_domain": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the workload domain is joined to the management domain's SSO domain.",
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	domainId := data.Get("domain_id").(string)
	domainName := data.Get("name").(string)

	if domainId == "" {
		if domainName == "" {
			return diag.Errorf("either 'domain_id' or 'name' must be provided")
		}

		domainInfo, err := getDomainByName(ctx, apiClient, domainName)
		if err != nil {
			return diag.FromErr(err)
		}
		domainId = *domainInfo.Id
	}

	_, err := domain.ImportDomain(ctx, data, apiClient, domainId, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getDomainByName(ctx context.Context, apiClient *vcf.ClientWithResponses, name string) (*vcf.Domain, error) {
	domainsResponse, err := apiClient.GetDomainsWithResponse(ctx, &vcf.GetDomainsParams{})
	if err != nil {
		return nil, err
	}

	resp, vcfErr := api_client.GetResponseAs[vcf.PageOfDomain](domainsResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}

	for _, domainElement := range *resp.Elements {
		if *domainElement.Name == name {
			return &domainElement, nil
		}
	}

	return nil, fmt.Errorf("domain name '%s' not found", name)
}
