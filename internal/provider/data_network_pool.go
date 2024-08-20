// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/network_pools"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func DataSourceNetworkPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkPoolRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network pool",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the network pool",
			},
			"network": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The network in the network pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The gateway of the network",
						},
						"mask": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet mask of the network",
						},
						"mtu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The MTU of the network",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet of the network",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of network",
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The VLAN ID of the network",
						},
						"ip_pools": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP pools associated with the network",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The starting IP address of the IP pool",
									},
									"end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ending IP address of the IP pool",
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

func dataSourceNetworkPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	name := d.Get("name").(string)
	networkPool, err := getNetworkPoolByName(ctx, apiClient, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkPool.ID)
	_ = d.Set("name", networkPool.Name)
	_ = d.Set("networks", flattenNetworks(networkPool.Networks))

	return nil
}

func getNetworkPoolByName(ctx context.Context, apiClient *client.VcfClient, name string) (*models.NetworkPool, error) {
	params := network_pools.NewGetNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)

	networkPoolsPayload, err := apiClient.NetworkPools.GetNetworkPool(params)
	if err != nil {
		return nil, err
	}

	if networkPoolsPayload.Payload == nil {
		return nil, errors.New("network pool not found")
	}

	for _, networkPool := range networkPoolsPayload.Payload.Elements {
		if networkPool.Name == name {
			return networkPool, nil
		}
	}

	return nil, errors.New("network pool not found")
}

func flattenNetworks(networks []*models.Network) []interface{} {
	if networks == nil {
		return []interface{}{}
	}

	var result []interface{}
	for _, network := range networks {
		n := map[string]interface{}{
			"gateway": network.Gateway,
			"mask":    network.Mask,
			"mtu":     network.Mtu,
			"subnet":  network.Subnet,
			"type":    network.Type,
			"vlan_id": network.VlanID,
		}

		if network.IPPools != nil {
			var ipPools []interface{}
			for _, ipPool := range network.IPPools {
				ipPoolMap := map[string]interface{}{
					"start": ipPool.Start,
					"end":   ipPool.End,
				}
				ipPools = append(ipPools, ipPoolMap)
			}
			n["ip_pools"] = ipPools
		}

		result = append(result, n)
	}

	return result
}
