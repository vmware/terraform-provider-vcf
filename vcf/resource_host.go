package vcf

import (
	"context"
	"encoding/json"
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
			"network_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network pool to associate the host with",
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
		},
	}
}

func resourceHostCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient
	log.Println(d)
	params := hosts.NewCommissionHostsParams()
	host := models.HostCommissionSpec{}

	if fqdn, ok := d.GetOk("fqdn"); ok {
		fqdnVal := fqdn.(string)
		host.Fqdn = &fqdnVal
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageTypeVal := storageType.(string)
		host.StorageType = &storageTypeVal
	}

	if username, ok := d.GetOk("username"); ok {
		usernameVal := username.(string)
		host.Username = &usernameVal
	}

	if password, ok := d.GetOk("password"); ok {
		passwordVal := password.(string)
		host.Password = &passwordVal
	}

	if networkPoolName, ok := d.GetOk("network_pool_name"); ok {
		networkPoolNameVal := networkPoolName.(string)

		// GetNetworkPools(params *GetNetworkPoolsParams, opts ...ClientOption) (*GetNetworkPoolsOK, error)
		ok, err := apiClient.NetworkPools.GETNetworkPools(nil)
		if err != nil {
			log.Println("error = ", err)
			return diag.FromErr(err)
		}

		found := false
		for _, networkPool := range ok.Payload.Elements {
			if networkPool.Name == networkPoolNameVal {
				host.NetworkPoolID = &networkPool.ID
				found = true
				break
			}
		}

		if !found {
			log.Println("Did not find network pool ", networkPoolNameVal)
			return diag.Errorf("Did not find network pool %s", networkPoolNameVal)
		}
	}

	params.HostCommissionSpecs = []*models.HostCommissionSpec{&host}

	_, accepted, err := apiClient.Hosts.CommissionHosts(params)
	if err != nil {
		log.Println("error = ", err)
		diag.FromErr(err)
	}

	log.Printf("%s host commission initiated. waiting for task id = %s", *host.Fqdn, accepted.Payload.ID)

	err = vcfClient.WaitForTaskComplete(accepted.Payload.ID)
	if err != nil {
		log.Println("error = ", err)
		diag.FromErr(err)
	}

	// Task complete, save the fqdn as id (required to decommission the host)
	d.SetId(*host.Fqdn)

	return nil
}

func resourceHostRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	// ID is the fqdn, but GET api needs uuid
	id := d.Id()
	hostId := ""
	if hostIdVal, ok := d.GetOk("host_id"); ok {
		hostId = hostIdVal.(string)
	}

	if hostId == "" {
		// Get all hosts and match the fqdn
		ok, err := apiClient.Hosts.GETHosts(nil)
		if err != nil {
			log.Println("error = ", err)
			diag.FromErr(err)
		}

		jsonp, _ := json.MarshalIndent(ok.Payload, " ", " ")
		log.Println(string(jsonp))

		// Check if the resource with the known id exists
		for _, host := range ok.Payload.Elements {
			if host.Fqdn == id {
				_ = d.Set("host_id", host.ID)
				return nil
			}
		}

		// Did not find the resource, set ID to ""
		log.Println("Did not find host with id ", id)
		d.SetId("")
		return nil
	} else {
		// Get a single host
		// GetHost(params *GetHostParams, opts ...ClientOption) (*GetHostOK, error)

		params := hosts.NewGETHostParams()
		params.ID = hostId

		_, err := apiClient.Hosts.GETHosts(nil)
		if err != nil {
			log.Println("error = ", err)
			diag.FromErr(err)
		}

		// Found the host
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
	host := models.HostDecommissionSpec{}
	id := d.Id()
	host.Fqdn = &id
	params.HostDecommissionSpecs = []*models.HostDecommissionSpec{&host}

	log.Println(params)

	// DecommissionHosts(params *DecommissionHostsParams, opts ...ClientOption) (*DecommissionHostsOK, *DecommissionHostsAccepted, error)
	_, accepted, err := apiClient.Hosts.DecommissionHosts(params)
	if err != nil {
		log.Println("error = ", err)
		diag.FromErr(err)
	}

	log.Printf("%s: Decommission task initiated. Task id %s", d.Id(), accepted.Payload.ID)
	err = vcfClient.WaitForTaskComplete(accepted.Payload.ID)
	if err != nil {
		log.Println("error = ", err)
		diag.FromErr(err)
	}

	// Task complete, clear the fqdn which is set as id
	d.SetId("")
	return nil
}
