// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/hosts"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

type HostModel struct {
	Fqdn                  types.String       `tfsdk:"fqdn"`
	Id                    types.String       `tfsdk:"id"`
	Status                types.String       `tfsdk:"status"`
	CompatibleStorageType types.String       `tfsdk:"compatible_storage_type"`
	Domain                []DomainModel      `tfsdk:"domain"`
	Cluster               []ClusterModel     `tfsdk:"cluster"`
	NetworkPool           []NetworkPoolModel `tfsdk:"network_pool"`
	Version               types.String       `tfsdk:"version"`
	Hardware              []HardwareModel    `tfsdk:"hardware"`
	IpAddresses           []IpAddressModel   `tfsdk:"ip_addresses"`
	Cpu                   CpuModel           `tfsdk:"cpu"`
	Memory                MemoryModel        `tfsdk:"memory"`
	Storage               []StorageModel     `tfsdk:"storage"`
	PhysicalNics          []PhysicalNicModel `tfsdk:"physical_nics"`
}

type DomainModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ClusterModel struct {
	Id types.String `tfsdk:"id"`
}

type NetworkPoolModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type HardwareModel struct {
	Model  types.String `tfsdk:"model"`
	Vendor types.String `tfsdk:"vendor"`
	Hybrid types.Bool   `tfsdk:"hybrid"`
}

type IpAddressModel struct {
	IpAddress types.String `tfsdk:"ip_address"`
	Type      types.String `tfsdk:"type"`
}

type CpuModel struct {
	Cores            types.Int64    `tfsdk:"cores"`
	CpuCoreDetails   []CpuCoreModel `tfsdk:"cpu_core_details"`
	FrequencyMHz     types.Float64  `tfsdk:"frequency_mhz"`
	UsedFrequencyMHz types.Float64  `tfsdk:"used_frequency_mhz"`
}

type CpuCoreModel struct {
	FrequencyMHz types.Float64 `tfsdk:"frequency_mhz"`
	Manufacturer types.String  `tfsdk:"manufacturer"`
	Model        types.String  `tfsdk:"model"`
}

type MemoryModel struct {
	TotalCapacityMB types.Float64 `tfsdk:"total_capacity_mb"`
	UsedCapacityMB  types.Float64 `tfsdk:"used_capacity_mb"`
}

type StorageModel struct {
	TotalCapacityMB types.Float64 `tfsdk:"total_capacity_mb"`
	UsedCapacityMB  types.Float64 `tfsdk:"used_capacity_mb"`
	Disks           []DiskModel   `tfsdk:"disks"`
}

type DiskModel struct {
	CapacityMB   types.Float64 `tfsdk:"capacity_mb"`
	DiskType     types.String  `tfsdk:"disk_type"`
	Manufacturer types.String  `tfsdk:"manufacturer"`
	Model        types.String  `tfsdk:"model"`
}

type PhysicalNicModel struct {
	DeviceName types.String `tfsdk:"device_name"`
	MacAddress types.String `tfsdk:"mac_address"`
	Speed      types.Int64  `tfsdk:"speed"`
	Unit       types.String `tfsdk:"unit"`
}

type DataSourceHost struct {
	client *client.VcfClient
}

func (d *DataSourceHost) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vcf_host"
}

func (d *DataSourceHost) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api_client.SddcManagerClient).ApiClient
}

func (d *DataSourceHost) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"fqdn": schema.StringAttribute{
				Required:    true,
				Description: "The fully qualified domain name of the ESXi host.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the ESXi host.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the ESXi host.",
			},
			"domain": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The domain details of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"cluster": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The clusters of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"network_pool": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The network pool of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the ESXi host.",
			},
			"compatible_storage_type": schema.StringAttribute{
				Computed:    true,
				Description: "The compatible storage type of the ESXi host.",
			},
			"hardware": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The hardware details of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"model": schema.StringAttribute{
							Computed: true,
						},
						"vendor": schema.StringAttribute{
							Computed: true,
						},
						"hybrid": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
			"ip_addresses": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The IP addresses of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"cpu": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The CPU details of the ESXi host.",
				Attributes: map[string]schema.Attribute{
					"cores": schema.Int64Attribute{
						Computed: true,
					},
					"cpu_cores": schema.ListNestedAttribute{
						Computed:    true,
						Description: "The CPU core details of the ESXi host.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"frequency_mhz": schema.Float64Attribute{
									Computed: true,
								},
								"manufacturer": schema.StringAttribute{
									Computed: true,
								},
								"model": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"frequency_mhz": schema.Float64Attribute{
						Computed: true,
					},
					"used_frequency_mhz": schema.Float64Attribute{
						Computed: true,
					},
				},
			},
			"memory": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The memory details of the ESXi host.",
				Attributes: map[string]schema.Attribute{
					"total_capacity_mb": schema.Float64Attribute{
						Computed: true,
					},
					"used_capacity_mb": schema.Float64Attribute{
						Computed: true,
					},
				},
			},
			"storage": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The storage details of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"total_capacity_mb": schema.Float64Attribute{
							Computed: true,
						},
						"used_capacity_mb": schema.Float64Attribute{
							Computed: true,
						},
						"disks": schema.ListNestedAttribute{
							Computed:    true,
							Description: "The disks of the ESXi host.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"capacity_mb": schema.Float64Attribute{
										Computed: true,
									},
									"disk_type": schema.StringAttribute{
										Computed: true,
									},
									"manufacturer": schema.StringAttribute{
										Computed: true,
									},
									"model": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"physical_nics": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The physical NICs of the ESXi host.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_name": schema.StringAttribute{
							Computed: true,
						},
						"mac_address": schema.StringAttribute{
							Computed: true,
						},
						"speed": schema.Int64Attribute{
							Computed: true,
						},
						"unit": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceHost) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data []HostModel

	params := hosts.NewGetHostsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)

	hostPayload, err := d.client.Hosts.GetHosts(params)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to retrieve hosts", err.Error()))
		return
	}

	for _, element := range hostPayload.Payload.Elements {
		// Debugging logs to print the values of id and status
		tflog.Debug(ctx, "Mapping host", map[string]interface{}{
			"fqdn":   element.Fqdn,
			"id":     element.ID,
			"status": element.Status,
		})

		host := mapHostElementToModel(element)
		data = append(data, host)
	}

	// Debugging log to trace the mapped data
	tflog.Debug(ctx, "Mapped data", map[string]interface{}{
		"hosts": data,
	})

	resp.State.Set(ctx, &data)
}

func mapHostElementToModel(element *models.Host) HostModel {
	// Inline dereferencing logic
	domainID := derefString(element.Domain.ID)
	clusterID := derefString(element.Cluster.ID)
	networkPoolID := derefString(element.Networkpool.ID)

	return HostModel{
		Fqdn:   types.StringValue(element.Fqdn),
		Id:     types.StringValue(element.ID),
		Status: types.StringValue(element.Status),
		Domain: []DomainModel{
			{
				Id:   types.StringValue(domainID),
				Name: types.StringValue(element.Domain.Name),
			},
		},
		Cluster: []ClusterModel{
			{
				Id: types.StringValue(clusterID),
			},
		},
		NetworkPool: []NetworkPoolModel{
			{
				Id:   types.StringValue(networkPoolID),
				Name: types.StringValue(element.Networkpool.Name),
			},
		},
		Version:               types.StringValue(element.EsxiVersion),
		CompatibleStorageType: types.StringValue(element.CompatibleStorageType),
		Hardware: []HardwareModel{
			{
				Model:  types.StringValue(element.HardwareModel),
				Vendor: types.StringValue(element.HardwareVendor),
				Hybrid: types.BoolValue(element.Hybrid),
			},
		},
		IpAddresses: []IpAddressModel{
			{
				IpAddress: types.StringValue(element.IPAddresses[0].IPAddress),
				Type:      types.StringValue(element.IPAddresses[0].Type),
			},
		},
		Cpu: CpuModel{
			Cores: types.Int64Value(int64(element.CPU.Cores)),
			CpuCoreDetails: []CpuCoreModel{
				{
					FrequencyMHz: types.Float64Value(element.CPU.CPUCores[0].FrequencyMHz),
					Manufacturer: types.StringValue(element.CPU.CPUCores[0].Manufacturer),
					Model:        types.StringValue(element.CPU.CPUCores[0].Model),
				},
			},
			FrequencyMHz:     types.Float64Value(element.CPU.FrequencyMHz),
			UsedFrequencyMHz: types.Float64Value(element.CPU.UsedFrequencyMHz),
		},
		Memory: MemoryModel{
			TotalCapacityMB: types.Float64Value(element.Memory.TotalCapacityMB),
			UsedCapacityMB:  types.Float64Value(element.Memory.UsedCapacityMB),
		},
		Storage: []StorageModel{
			{
				TotalCapacityMB: types.Float64Value(element.Storage.TotalCapacityMB),
				UsedCapacityMB:  types.Float64Value(element.Storage.UsedCapacityMB),
				Disks: []DiskModel{
					{
						CapacityMB:   types.Float64Value(element.Storage.Disks[0].CapacityMB),
						DiskType:     types.StringValue(element.Storage.Disks[0].DiskType),
						Manufacturer: types.StringValue(element.Storage.Disks[0].Manufacturer),
						Model:        types.StringValue(element.Storage.Disks[0].Model),
					},
				},
			},
		},
		PhysicalNics: []PhysicalNicModel{
			{
				DeviceName: types.StringValue(element.PhysicalNics[0].DeviceName),
				MacAddress: types.StringValue(element.PhysicalNics[0].MacAddress),
				Speed:      types.Int64Value(element.PhysicalNics[0].Speed),
				Unit:       types.StringValue(element.PhysicalNics[0].Unit),
			},
		},
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
