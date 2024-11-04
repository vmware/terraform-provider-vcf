// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
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
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
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
	client *vcf.ClientWithResponses
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

	networkPool := vcf.NetworkPool{
		Name: data.Name.ValueString(),
	}

	var networks []NetworkModel
	res.Diagnostics.Append(types.List.ElementsAs(data.Networks, ctx, &networks, false)...)

	if len(networks) > 0 {
		networkPool.Networks = make([]vcf.Network, len(networks))

		for i, network := range networks {
			networkPool.Networks[i] = vcf.Network{
				VlanId:  int32(network.VlanId.ValueInt64()),
				Gateway: network.Gateway.ValueString(),
				Mask:    network.Mask.ValueString(),
				Subnet:  network.Subnet.ValueString(),
				Mtu:     int32(network.Mtu.ValueInt64()),
				Type:    network.Type.ValueString(),
			}

			var ipPoolModels []IpPoolModel
			res.Diagnostics.Append(types.List.ElementsAs(network.IpPools, ctx, &ipPoolModels, false)...)
			ipPoolsSlice := make([]vcf.IpPool, len(ipPoolModels))
			networkPool.Networks[i].IpPools = &ipPoolsSlice
			for j, ipPoolModel := range ipPoolModels {
				ipPools := networkPool.Networks[i].IpPools
				if ipPools != nil {
					(*ipPools)[j].Start = ipPoolModel.Start.ValueString()
					(*ipPools)[j].End = ipPoolModel.End.ValueString()
				}
			}
		}
	}

	created, err := r.client.CreateNetworkPoolWithResponse(ctx, networkPool)
	if err != nil {
		res.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to create network pool", err.Error()))
		return
	}
	pool, vcfErr := api_client.GetResponseAs[vcf.NetworkPool](created.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return
	}

	log.Println("created = ", pool.Name)
	data.Id = types.StringValue(*pool.Id)

	res.Diagnostics.Append(res.State.Set(ctx, &data)...)
}

func (r *ResourceNetworkPool) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var data ResourceNetworkPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	networkPoolPayload, _ := r.client.GetNetworkPoolByIDWithResponse(ctx, data.Id.ValueString())
	pool, vcfErr := api_client.GetResponseAs[vcf.NetworkPool](networkPoolPayload.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return
	}

	data.Id = types.StringValue(*pool.Id)
	data.Name = types.StringValue(pool.Name)
}

func (r *ResourceNetworkPool) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	/* ... */
}

func (r *ResourceNetworkPool) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var data ResourceNetworkPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	networkPoolPayload, _ := r.client.DeleteNetworkPoolWithResponse(ctx, data.Id.ValueString(), nil)
	_, vcfErr := api_client.GetResponseAs[vcf.NetworkPool](networkPoolPayload.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return
	}

	log.Printf("%s: Delete complete", data.Id)
}
