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
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/credentials"
	"github.com/vmware/vcf-sdk-go/client/hosts"
	"github.com/vmware/vcf-sdk-go/client/network_pools"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
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
	params := hosts.NewCommissionHostsParamsWithTimeout(constants.DefaultVcfApiCallTimeout)
	commissionSpec := models.HostCommissionSpec{}

	if fqdn, ok := d.GetOk("fqdn"); ok {
		fqdnVal := fqdn.(string)
		commissionSpec.Fqdn = &fqdnVal
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageTypeVal := storageType.(string)
		commissionSpec.StorageType = &storageTypeVal
	}

	if username, ok := d.GetOk("username"); ok {
		usernameVal := username.(string)
		commissionSpec.Username = &usernameVal
	}

	if password, ok := d.GetOk("password"); ok {
		passwordVal := password.(string)
		commissionSpec.Password = &passwordVal
	}

	if networkPoolId, ok := d.GetOk("network_pool_id"); ok {
		networkPoolIdStr := networkPoolId.(string)
		commissionSpec.NetworkPoolID = &networkPoolIdStr
	}

	if networkPoolName, ok := d.GetOk("network_pool_name"); ok {
		if commissionSpec.NetworkPoolID != nil {
			return diag.FromErr(errors.New("you cannot set network_pool_id and network_pool_name at the same time"))
		}

		networkPool, err := getNetworkPool(networkPoolName.(string), apiClient)

		if err != nil {
			return diag.FromErr(err)
		}

		commissionSpec.NetworkPoolID = &networkPool.ID
	}

	params.HostCommissionSpecs = []*models.HostCommissionSpec{&commissionSpec}

	_, accepted, err := apiClient.Hosts.CommissionHosts(params)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	taskId := accepted.Payload.ID

	tflog.Info(ctx, fmt.Sprintf("%s commissionSpec commission initiated. waiting for task id = %s",
		*commissionSpec.Fqdn, taskId))

	err = vcfClient.WaitForTaskComplete(ctx, taskId, false)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	hostId, err := vcfClient.GetResourceIdAssociatedWithTask(ctx, taskId, "Esxi")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hostId)

	return resourceHostRead(ctx, d, meta)
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	hostId := d.Id()

	getHostParams := hosts.NewGetHostParams().WithTimeout(constants.DefaultVcfApiCallTimeout)
	getHostParams.ID = hostId

	hostResponse, err := apiClient.Hosts.GetHost(getHostParams)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	host := hostResponse.Payload

	_ = d.Set("network_pool_id", host.Networkpool.ID)
	_ = d.Set("network_pool_name", host.Networkpool.Name)
	_ = d.Set("fqdn", host.Fqdn)
	_ = d.Set("status", host.Status)

	getHostCredentialsParams := credentials.NewGetCredentialsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).WithResourceName(&host.Fqdn)
	getCredentialsResponse, err := apiClient.Credentials.GetCredentials(getHostCredentialsParams)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	for _, credential := range getCredentialsResponse.Payload.Elements {
		if credential == nil {
			continue
		}
		// we're interested in the SSH credentials, not service account
		if *credential.AccountType != "USER" || *credential.CredentialType != "SSH" {
			continue
		}
		if *credential.Resource.ResourceID != hostId {
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

	params := hosts.NewDecommissionHostsParamsWithTimeout(constants.DefaultVcfApiCallTimeout)
	decommissionSpec := models.HostDecommissionSpec{}
	decommissionSpec.Fqdn = resource_utils.ToStringPointer(d.Get("fqdn"))
	params.HostDecommissionSpecs = []*models.HostDecommissionSpec{&decommissionSpec}

	log.Println(params)

	_, accepted, err := apiClient.Hosts.DecommissionHosts(params)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	log.Printf("%s %s: Decommission task initiated. Task id %s",
		d.Get("fqdn").(string), d.Id(), accepted.Payload.ID)
	err = vcfClient.WaitForTaskComplete(ctx, accepted.Payload.ID, false)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	return nil
}

func getNetworkPool(name string, client *client.VcfClient) (*models.NetworkPool, error) {
	params := network_pools.NewGetNetworkPoolParams().WithTimeout(constants.DefaultVcfApiCallTimeout)

	ok, err := client.NetworkPools.GetNetworkPool(params)

	if err != nil {
		return nil, err
	}

	networkPools := ok.Payload.Elements

	if len(networkPools) > 0 {
		for _, pool := range networkPools {
			if pool.Name == name {
				return pool, nil
			}
		}
	}

	return nil, fmt.Errorf("network pool %s not found", name)
}
