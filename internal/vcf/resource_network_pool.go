/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcf

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/client/network_pools"
	"github.com/vmware/vcf-sdk-go/models"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNetworkPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPoolCreate,
		ReadContext:   resourceNetworkPoolRead,
		UpdateContext: resourceNetworkPoolUpdate,
		DeleteContext: resourceNetworkPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network pool",
			},
			"network": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Represents a network in a network pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway": {
							Type:        schema.TypeString,
							Description: "Gateway for the network",
							Optional:    true,
						},
						"mask": {
							Type:        schema.TypeString,
							Description: "Subnet mask for the subnet of the network",
							Optional:    true,
						},
						"mtu": {
							Type:        schema.TypeInt,
							Description: "Gateway for the network",
							Optional:    true,
						},
						"subnet": {
							Type:        schema.TypeString,
							Description: "Subnet associated with the network",
							Optional:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Network Type of the network",
							Optional:    true,
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Description: "VLAN ID associated with the network",
							Optional:    true,
						},
						"ip_pools": {
							Type:        schema.TypeList,
							Description: "List of IP pool ranges to use",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:        schema.TypeString,
										Description: "Start IP address of the IP pool",
										Optional:    true,
									},
									"end": {
										Type:        schema.TypeString,
										Description: "End IP address of the IP pool",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNetworkPoolCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	params := network_pools.NewCreateNetworkPoolParams()
	networkPool := models.NetworkPool{}

	if name, ok := d.GetOk("name"); ok {
		networkPool.Name = name.(string)
	}

	if len(d.Get("network").([]interface{})) > 0 {
		networks := d.Get("network").([]interface{})
		networkPool.Networks = make([]*models.Network, len(networks))

		for i, network := range networks {
			networkMap := network.(map[string]interface{})
			log.Println(spew.Sdump(networkMap))
			networkPool.Networks[i] = &models.Network{
				Gateway: networkMap["gateway"].(string),
				Mask:    networkMap["mask"].(string),
				Mtu:     int32(networkMap["mtu"].(int)),
				Subnet:  networkMap["subnet"].(string),
				Type:    networkMap["type"].(string),
				VlanID:  int32(networkMap["vlan_id"].(int)),
			}

			ipPools := networkMap["ip_pools"].([]interface{})
			networkPool.Networks[i].IPPools = make([]*models.IPPool, len(ipPools))
			for j, ipPool := range ipPools {
				ipPoolMap := ipPool.(map[string]interface{})

				networkPool.Networks[i].IPPools[j] = &models.IPPool{
					Start: ipPoolMap["start"].(string),
					End:   ipPoolMap["end"].(string),
				}
			}
		}
	}

	params.NetworkPool = &networkPool
	log.Println(spew.Sdump(params.NetworkPool))

	_, created, err := apiClient.NetworkPools.CreateNetworkPool(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Println("created = ", created)
	createdNetworkPool := created.Payload
	d.SetId(createdNetworkPool.ID)

	return nil
}

func resourceNetworkPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	params := network_pools.NewGETNetworkPoolParams()
	params.ID = d.Id()

	_, err := apiClient.NetworkPools.GETNetworkPool(params)
	if err != nil {
		log.Println("error = ", err)
		// TODO : Look for not found error
		/*
			log.Println("Did not find network pool with id ", id)
			d.SetId("")
			return nil
		*/
		return diag.FromErr(err)
	}

	// jsonp, _ := json.MarshalIndent(ok.Payload, " ", " ")
	// log.Println(string(jsonp))
	return nil
}

/**
 * Updating network pools is partially supported in VCF API.
 * ipPools can be added or removed, but not yet implemented here in the provider.
 */
func resourceNetworkPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNetworkPoolRead(ctx, d, meta)
}
func resourceNetworkPoolDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	params := network_pools.NewDeleteNetworkPoolParams()
	params.ID = d.Id()

	log.Println(params)
	_, err := apiClient.NetworkPools.DeleteNetworkPool(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Printf("%s: Delete complete", d.Id())
	d.SetId("")
	return nil
}
