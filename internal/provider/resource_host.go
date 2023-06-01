/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/client/hosts"
	"github.com/vmware/vcf-sdk-go/models"

	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FQDN of the host",
			},
			"network_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the network pool to associate the host with",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Storage Type. One among: VSAN, VSAN_REMOTE, NFS, VMFS_FC, VVOL",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the host",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the host",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the host. Known after commissioning.",
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
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient
	log.Println(d)
	params := hosts.NewCommissionHostsParams()
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

	params.HostCommissionSpecs = []*models.HostCommissionSpec{&commissionSpec}

	_, accepted, err := apiClient.Hosts.CommissionHosts(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Printf("%s commissionSpec commission initiated. waiting for task id = %s", *commissionSpec.Fqdn, accepted.Payload.ID)

	err = vcfClient.WaitForTaskComplete(accepted.Payload.ID)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	// Task complete, save the fqdn as id (required to decommission the commissionSpec)
	d.SetId(*commissionSpec.Fqdn)

	return resourceHostRead(ctx, d, meta)
}

func resourceHostRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	// ID is the fqdn, but GET api needs uuid
	hostFqdn := d.Id()
	hostUuid := ""
	if hostIdVal, ok := d.GetOk("host_id"); ok {
		hostUuid = hostIdVal.(string)
	}

	if hostUuid == "" {
		// Get all hosts and match the fqdn
		ok, err := apiClient.Hosts.GetHosts(nil)
		if err != nil {
			log.Println("error = ", err)
			diag.FromErr(err)
		}

		// Check if the resource with the known hostFqdn exists
		for _, host := range ok.Payload.Elements {
			if host.Fqdn == hostFqdn {
				_ = d.Set("host_id", host.ID)
				// storage_type is not returned by the VCF API ?!?
				_ = d.Set("network_pool_id", host.Networkpool.ID)
				_ = d.Set("fqdn", host.Fqdn)
				_ = d.Set("status", host.Status)
				return nil
			}
		}

		// Did not find the resource, set ID to ""
		log.Println("Did not find host with hostFqdn ", hostFqdn)
		d.SetId("")
		return nil
	} else {
		// Get a single host using UUID
		params := hosts.NewGetHostParams()
		params.ID = hostUuid

		host, err := apiClient.Hosts.GetHost(params)
		if err != nil {
			log.Println("error = ", err)
			return diag.FromErr(err)
		}
		_ = d.Set("host_id", host.Payload.ID)
		_ = d.Set("network_pool_id", host.Payload.Networkpool.ID)
		_ = d.Set("fqdn", host.Payload.Fqdn)
		_ = d.Set("status", host.Payload.Status)
		return nil
	}
}

/**
 * Updating hosts is not supported in VCF API.
 */
func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceHostRead(ctx, d, meta)
}
func resourceHostDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	params := hosts.NewDecommissionHostsParams()
	decommissionSpec := models.HostDecommissionSpec{}
	id := d.Id()
	decommissionSpec.Fqdn = &id
	params.HostDecommissionSpecs = []*models.HostDecommissionSpec{&decommissionSpec}

	log.Println(params)

	// DecommissionHosts(params *DecommissionHostsParams, opts ...ClientOption) (*DecommissionHostsOK, *DecommissionHostsAccepted, error)
	_, accepted, err := apiClient.Hosts.DecommissionHosts(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Printf("%s: Decommission task initiated. Task id %s", d.Id(), accepted.Payload.ID)
	err = vcfClient.WaitForTaskComplete(accepted.Payload.ID)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	return nil
}
