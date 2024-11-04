// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
)

func ResourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fully qualified domain name of ESXi host",
			},
			"network_pool_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Name of the network pool to associate the ESXi host with",
				ConflictsWith: []string{"network_pool_id"},
			},
			"network_pool_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "ID of the network pool to associate the ESXi host with",
				ConflictsWith: []string{"network_pool_name"},
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Storage Type. One among: VSAN, VSAN_ESA, VSAN_REMOTE, NFS, VMFS_FC, VVOL",
				ValidateFunc: validation.StringInSlice([]string{"VSAN", "VSAN_ESA", "VSAN_REMOTE", "NFS", "VMFS_FC", "VVOL"}, false),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to authenticate to the ESXi host",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password to authenticate to the ESXi host",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignable status of the host.",
			},
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	commissionSpec := vcf.HostCommissionSpec{}

	if fqdn, ok := d.GetOk("fqdn"); ok {
		commissionSpec.Fqdn = fqdn.(string)
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		commissionSpec.StorageType = storageType.(string)
	}

	if username, ok := d.GetOk("username"); ok {
		commissionSpec.Username = username.(string)
	}

	if password, ok := d.GetOk("password"); ok {
		commissionSpec.Password = password.(string)
	}

	if networkPoolId, ok := d.GetOk("network_pool_id"); ok {
		commissionSpec.NetworkPoolId = networkPoolId.(string)
	}

	if networkPoolName, ok := d.GetOk("network_pool_name"); ok {
		if commissionSpec.NetworkPoolId != "" {
			return diag.FromErr(errors.New("you cannot set network_pool_id and network_pool_name at the same time"))
		}

		networkPool, err := getNetworkPool(networkPoolName.(string), apiClient, ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		commissionSpec.NetworkPoolId = *networkPool.Id
	}

	accepted, err := apiClient.CommissionHostsWithResponse(ctx, []vcf.HostCommissionSpec{commissionSpec})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](accepted.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	tflog.Info(ctx, fmt.Sprintf("%s commissionSpec commission initiated. waiting for task id = %s",
		commissionSpec.Fqdn, *task.Id))

	err = vcfClient.WaitForTaskComplete(ctx, *task.Id, false)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	hostId, err := vcfClient.GetResourceIdAssociatedWithTask(ctx, *task.Id, "Esxi")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hostId)

	return resourceHostRead(ctx, d, meta)
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	hostResponse, err := apiClient.GetHostWithResponse(ctx, d.Id())
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	host, vcfErr := api_client.GetResponseAs[vcf.Host](hostResponse.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	_ = d.Set("network_pool_id", host.Networkpool.Id)
	_ = d.Set("network_pool_name", host.Networkpool.Name)
	_ = d.Set("fqdn", host.Fqdn)
	_ = d.Set("status", host.Status)
	params := &vcf.GetCredentialsParams{
		ResourceName: host.Fqdn,
	}
	getCredentialsResponse, err := apiClient.GetCredentialsWithResponse(ctx, params)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfCredential](getCredentialsResponse.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}
	for _, credential := range *page.Elements {
		// we're interested in the SSH credentials, not service account
		if *credential.AccountType != "USER" || *credential.CredentialType != "SSH" {
			continue
		}
		if credential.Resource.ResourceId != *host.Id {
			return diag.FromErr(fmt.Errorf("hostId doesn't match host FQDN when requesting credentials"))
		}
		_ = d.Set("username", *credential.Username)
		_ = d.Set("password", credential.Password)
	}

	return nil
}

// There is no update method for commissioned hosts.
func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDomainRead(ctx, d, meta)
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	decommissionSpec := vcf.HostDecommissionSpec{}
	decommissionSpec.Fqdn = d.Get("fqdn").(string)

	accepted, err := apiClient.DecommissionHostsWithResponse(ctx, []vcf.HostDecommissionSpec{decommissionSpec})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](accepted.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	log.Printf("%s %s: Decommission task initiated. Task id %s",
		d.Get("fqdn").(string), d.Id(), *task.Id)
	err = vcfClient.WaitForTaskComplete(ctx, *task.Id, false)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	return nil
}

func getNetworkPool(name string, client *vcf.ClientWithResponses, ctx context.Context) (*vcf.NetworkPool, error) {
	ok, err := client.GetNetworkPoolWithResponse(ctx)

	if err != nil {
		return nil, err
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfNetworkPool](ok.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}

	if page != nil {
		for _, pool := range *page.Elements {
			if pool.Name == name {
				return &pool, nil
			}
		}
	}

	return nil, fmt.Errorf("network pool %s not found", name)
}
