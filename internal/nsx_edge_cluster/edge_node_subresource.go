// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsx_edge_cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func EdgeNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the edge node",
				ValidateFunc: validation.NoZeroValues,
			},
			"compute_cluster_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the compute cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"compute_cluster_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The id of the compute cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"admin_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The administrator password for the edge node",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"audit_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The audit password for the edge node",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The root user password for the edge node",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"tep1_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The IP address (CIDR) of the first tunnel endpoint",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"tep2_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The IP address (CIDR) of the second tunnel endpoint",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"tep_gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The gateway for the tunnel endpoints",
				ValidateFunc: validationUtils.ValidateIPv4AddressSchema,
			},
			"tep_vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The VLAN ID for the tunnel endpoint",
				ValidateFunc: validation.IntBetween(0, 4095),
			},
			"inter_rack_cluster": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether or not this is an inter-rack cluster. True for L2 non-uniform and L3, false for L2 uniform",
			},
			"management_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The IP address (CIDR) for the management network",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"management_gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The gateway address for the management network",
				ValidateFunc: validationUtils.ValidateIPv4AddressSchema,
			},
			"management_network": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The management network which will be created for this node",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"portgroup_name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The name of the portgroup",
							ValidateFunc: validation.NoZeroValues,
						},
						"vlan_id": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "The VLAN ID for the portgroup",
							ValidateFunc: validation.IntBetween(0, 4095),
						},
					},
				},
			},
			"first_nsx_vds_uplink": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the first NSX-enabled VDS uplink",
				ValidateFunc: validation.NoZeroValues,
			},
			"second_nsx_vds_uplink": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the second NSX-enabled VDS uplink",
				ValidateFunc: validation.NoZeroValues,
			},
			"uplink": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifications of Tier-0 uplinks for the edge node",
				Elem:        UplinkNetworkSchema(),
			},
		},
	}
}

func UplinkNetworkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The IP address (CIDR) for the distributed switch uplink",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The VLAN ID for the distributed switch uplink",
				ValidateFunc: validation.IntBetween(0, 4095),
			},
			"bgp_peer": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of BGP Peer configurations",
				Elem:        BgpPeerSchema(),
			},
		},
	}
}

func BgpPeerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IP address",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password",
			},
			"asn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ASN",
				ValidateFunc: validationUtils.ValidASN,
			},
		},
	}
}
