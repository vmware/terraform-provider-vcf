// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/credentials"
	"strconv"
)

func DataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCredentialsRead,
		Description: "Datasource used to extract credentials for different resources that are part of the SDDC deployment based on name, ip, type, domain or account type",
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the resource",
			},
			"resource_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The IP Address of the resource",
				ValidateFunc: validation.IsIPAddress,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The type of the resource. One among ESXI, VCENTER, PSC, NSX_MANAGER, NSX_CONTROLLER, NSXT_EDGE, NSXT_MANAGER, VRLI, VROPS, VRA, WSA, VRSLCM, VXRAIL_MANAGER, NSX_ALB, BACKUP",
				ValidateFunc: validation.StringInSlice(credentials.AllResourceTypes(), true),
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain in which context we do the credentials read.",
			},
			"account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type(s) of the account.One among USER, SYSTEM, SERVICE",
				ValidateFunc: validation.StringInSlice(
					credentials.AllAccountTypes(), true),
			},
			"page": {
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				Description:  "The page of credentials that is returned as result.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"page_size": {
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				Description:  "the size of the credentials list . Default is 0 so the user will get all records in one page",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"credentials": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of credentials read from the API",
				Elem:        credentials.CredentialDataSubresource(),
			},
		},
	}
}

func dataCredentialsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient
	creds, err := credentials.ReadCredentials(ctx, data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	flatCredentials := credentials.FlattenCredentials(creds)
	_ = data.Set("credentials", flatCredentials)

	id, err := createCredentialID(data)
	if err != nil {
		return diag.Errorf("error during id generation %s", err)
	}

	data.SetId(id)

	return nil
}

func createCredentialID(data *schema.ResourceData) (string, error) {
	params := []string{
		data.Get("resource_name").(string),
		data.Get("resource_ip").(string),
		data.Get("resource_type").(string),
		data.Get("domain_name").(string),
		data.Get("account_type").(string),
		strconv.Itoa(data.Get("page").(int)),
		strconv.Itoa(data.Get("page_size").(int)),
	}

	return credentials.HashFields(params)
}
