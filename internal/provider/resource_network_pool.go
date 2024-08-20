// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/network_pools"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

type IpPoolModel struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type NetworkModel struct {
	Gateway types.String `tfsdk:"gateway"`
	Mask    types.String `tfsdk:"mask"`
	Subnet  types.String `tfsdk:"subnet"`
	Type    types.String `tfsdk:"type"`
	Mtu     types.Int64  `tfsdk:"mtu"`
	VlanId  types.Int64  `tfsdk:"vlan_id"`
	IpPools types.List   `tfsdk:"ip_pools"`
}

type ResourceNetworkPoolModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Id       types.String   `tfsdk:"id"`
	Name     types.String   `tfsdk:"name"`
	Networks types.List     `tfsdk:"network"`
}

type ResourceNetworkPool struct {
	client *client.VcfClient
}

func (r *ResourceNetworkPool) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "vcf_network_pool"
}

func (r *ResourceNetworkPool) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api_client.SddcManagerClient).ApiClient
}

func (r *ResourceNetworkPool) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceNetworkPool) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier for the network pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the network pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"network": schema.ListNestedBlock{
				Description: "Represents a network in a network pool",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"gateway": schema.StringAttribute{
							Optional:    true,
							Description: "Gateway for the network",
						},
						"mask": schema.StringAttribute{
							Optional:    true,
							Description: "Subnet mask for the subnet of the network",
						},
						"subnet": schema.StringAttribute{
							Optional:    true,
							Description: "Subnet associated with the network",
						},
						"mtu": schema.Int64Attribute{
							Optional:    true,
							Description: "Gateway for the network",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "Network Type of the network",
						},
						"vlan_id": schema.Int64Attribute{
							Required:    true,
							Description: "VLAN ID associated with the network",
						},
					},
					Blocks: map[string]schema.Block{
						"ip_pools": schema.ListNestedBlock{
							Description: "List of IP pool ranges to use",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"start": schema.StringAttribute{
										Optional:    true,
										Description: "Start IP address of the IP pool",
									},
									"end": schema.StringAttribute{
										Optional:    true,
										Description: "End IP address of the IP pool",
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

func (r *ResourceNetworkPool) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var data ResourceNetworkPoolModel

	res.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Create(ctx, 30*time.Minute)
	res.Diagnostics.Append(diags...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	createParams := network_pools.NewCreateNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	networkPool := models.NetworkPool{
		Name: data.Name.ValueString(),
	}

	var networks []NetworkModel
	res.Diagnostics.Append(types.List.ElementsAs(data.Networks, ctx, &networks, false)...)

	if len(networks) > 0 {
		networkPool.Networks = make([]*models.Network, len(networks))

		for i, network := range networks {
			networkPool.Networks[i] = &models.Network{
				VlanID:  int32(network.VlanId.ValueInt64()),
				Gateway: network.Gateway.ValueString(),
				Mask:    network.Mask.ValueString(),
				Subnet:  network.Subnet.ValueString(),
				Mtu:     int32(network.Mtu.ValueInt64()),
				Type:    network.Type.ValueString(),
			}

			var ipPools []IpPoolModel
			res.Diagnostics.Append(types.List.ElementsAs(network.IpPools, ctx, &ipPools, false)...)
			networkPool.Networks[i].IPPools = make([]*models.IPPool, len(ipPools))
			for j, ipPool := range ipPools {
				networkPool.Networks[i].IPPools[j] = &models.IPPool{
					Start: ipPool.Start.ValueString(),
					End:   ipPool.End.ValueString(),
				}
			}
		}
	}

	createParams.NetworkPool = &networkPool

	_, created, err := r.client.NetworkPools.CreateNetworkPool(createParams)
	if err != nil {
		res.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to create network pool", err.Error()))
		return
	}

	log.Println("created = ", created)
	createdNetworkPool := created.Payload
	data.Id = types.StringValue(createdNetworkPool.ID)

	res.Diagnostics.Append(res.State.Set(ctx, &data)...)
}

func (r *ResourceNetworkPool) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var data ResourceNetworkPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	params := network_pools.NewGetNetworkPoolByIDParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)

	networkPoolPayload, err := r.client.NetworkPools.GetNetworkPoolByID(params)
	if err != nil {
		res.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to retrieve network pool", err.Error()))
		return
	}

	networkPool := networkPoolPayload.Payload
	data.Id = types.StringValue(networkPool.ID)
	data.Name = types.StringValue(networkPool.Name)
}

func (r *ResourceNetworkPool) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	/* ... */
}

func (r *ResourceNetworkPool) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var data ResourceNetworkPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	params := network_pools.NewDeleteNetworkPoolParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	params.ID = data.Id.ValueString()

	log.Println(params)
	_, err := r.client.NetworkPools.DeleteNetworkPool(params)
	if err != nil {
		log.Println("error = ", err)
		res.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to delete network pool", err.Error()))
		return
	}

	log.Printf("%s: Delete complete", data.Id)
}
