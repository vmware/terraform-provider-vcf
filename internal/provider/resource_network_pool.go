/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/vcf-sdk-go/client/network_pools"
	"github.com/vmware/vcf-sdk-go/models"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNetworkPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPoolCreate,
		ReadContext:   resourceNetworkPoolRead,
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
				ForceNew:    true, // Updating network pools is partially supported in VCF API.
				Description: "The name of the network pool",
			},
			"network": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true, // Updating network pools is partially supported in VCF API.
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

func resourceNetworkPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	createParams := network_pools.NewCreateNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	networkPool := models.NetworkPool{}

	if name, ok := d.GetOk("name"); ok {
		networkPool.Name = name.(string)
	}

	if len(d.Get("network").([]interface{})) > 0 {
		networks := d.Get("network").([]interface{})
		networkPool.Networks = make([]*models.Network, len(networks))

		for i, network := range networks {
			networkMap := network.(map[string]interface{})
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

	createParams.NetworkPool = &networkPool

	_, created, err := apiClient.NetworkPools.CreateNetworkPool(createParams)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("created = ", created)
	createdNetworkPool := created.Payload
	d.SetId(createdNetworkPool.ID)

	return nil
}

func resourceNetworkPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	params := network_pools.NewGetNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	params.ID = d.Id()

	networkPoolPayload, err := apiClient.NetworkPools.GetNetworkPool(params)
	if err != nil {
		return diag.FromErr(err)
	}
	networkPool := networkPoolPayload.Payload
	d.SetId(networkPool.ID)
	_ = d.Set("name", networkPool.Name)

	return nil
}

func resourceNetworkPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	params := network_pools.NewDeleteNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
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
