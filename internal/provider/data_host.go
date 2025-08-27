// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/vcf-sdk-go/vcf"
)

func DataSourceHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the ESXi host.",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified domain name of the ESXi host.",
			},
			"hardware": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The hardware information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hardware model of the ESXi host.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hardware vendor of the ESXi host.",
						},
						"hybrid": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the ESXi host is hybrid.",
						},
					},
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the ESXi running on the host.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ESXi host.",
			},
			"domain": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The workload domain information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the workload domain.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the workload domain.",
						},
					},
				},
			},
			"cluster": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cluster information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the cluster.",
						},
					},
				},
			},
			"network_pool": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The network pool associated with the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network pool.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the network pool.",
						},
					},
				},
			},
			"cpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CPU information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cores": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"cpu_cores": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about each of the CPU cores.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency_mhz": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The frequency of the CPU core in MHz.",
									},
									"manufacturer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The manufacturer of the CPU.",
									},
									"model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The model of the CPU.",
									},
								},
							},
						},
						"frequency_mhz": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total CPU frequency in MHz.",
						},
						"used_frequency_mhz": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Used CPU frequency in MHz.",
						},
					},
				},
			},
			"memory": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The memory information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_capacity_mb": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The total memory capacity in MB.",
						},
						"used_capacity_mb": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The used memory capacity in MB.",
						},
					},
				},
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The storage type of the ESXi host.",
			},
			"storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The storage information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_capacity_mb": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The total storage capacity in MB.",
						},
						"used_capacity_mb": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The used storage capacity in MB.",
						},
						"disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The disks information of the ESXi host.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity_mb": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The capacity of the disk in MB.",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the disk.",
									},
									"manufacturer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The manufacturer of the disk.",
									},
									"model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The model of the disk.",
									},
								},
							},
						},
					},
				},
			},
			"physical_nics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The physical NICs information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device name of the NIC.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address of the NIC.",
						},
						"speed": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The speed of the NIC.",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unit of the NIC speed.",
						},
					},
				},
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IP addresses information of the ESXi host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the IP address.",
						},
					},
				},
			},
		},
	}
}

func dataSourceHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	fqdn := d.Get("fqdn").(string)
	host, err := getHostByFqdn(ctx, apiClient, fqdn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*host.Id)

	// Fully qualified domain name information.
	_ = d.Set("fqdn", host.Fqdn)

	// Hardware information.
	hardware := []map[string]interface{}{
		{
			"model":  host.HardwareModel,
			"vendor": host.HardwareVendor,
			"hybrid": host.Hybrid,
		},
	}
	_ = d.Set("hardware", hardware)

	// ESXi version information.
	_ = d.Set("version", host.EsxiVersion)

	// Status information.
	_ = d.Set("status", host.Status)

	// Domain information.
	domain := []map[string]interface{}{
		{
			"id":   host.Domain.Id,
			"name": host.Domain.Name,
		},
	}
	_ = d.Set("domain", domain)

	// Cluster information.
	cluster := []map[string]interface{}{
		{
			"id": host.Cluster.Id,
		},
	}
	_ = d.Set("cluster", cluster)

	// Network pool information.
	networkPool := []map[string]interface{}{
		{
			"id":   host.Networkpool.Id,
			"name": host.Networkpool.Name,
		},
	}
	_ = d.Set("network_pool", networkPool)

	// CPU information.
	var cpuCores []map[string]interface{}
	for _, core := range *host.Cpu.CpuCores {
		cpuCore := map[string]interface{}{
			"frequency_mhz": core.FrequencyMHz,
			"manufacturer":  core.Manufacturer,
			"model":         core.Model,
		}
		cpuCores = append(cpuCores, cpuCore)
	}

	cpu := map[string]interface{}{
		"cores":              host.Cpu.Cores,
		"cpu_cores":          cpuCores,
		"frequency_mhz":      host.Cpu.FrequencyMHz,
		"used_frequency_mhz": host.Cpu.UsedFrequencyMHz,
	}

	if err := d.Set("cpu", []interface{}{cpu}); err != nil {
		return diag.FromErr(err)
	}

	// Memory information.
	memory := []map[string]interface{}{
		{
			"total_capacity_mb": host.Memory.TotalCapacityMB,
			"used_capacity_mb":  host.Memory.UsedCapacityMB,
		},
	}
	_ = d.Set("memory", memory)

	// Compatible storage type information.
	_ = d.Set("compatible_storage_type", host.CompatibleStorageType)

	// Storage information.
	var disks []map[string]interface{}
	for _, disk := range *host.Storage.Disks {
		diskInfo := map[string]interface{}{
			"capacity_mb":  disk.CapacityMB,
			"disk_type":    disk.DiskType,
			"manufacturer": disk.Manufacturer,
			"model":        disk.Model,
		}
		disks = append(disks, diskInfo)
	}

	storage := []map[string]interface{}{
		{
			"total_capacity_mb": host.Storage.TotalCapacityMB,
			"used_capacity_mb":  host.Storage.UsedCapacityMB,
			"disks":             disks,
		},
	}
	_ = d.Set("storage", storage)

	// Physical NICs information.
	var physicalNics []map[string]interface{}
	for _, nic := range *host.PhysicalNics {
		physicalNic := map[string]interface{}{
			"device_name": nic.DeviceName,
			"mac_address": nic.MacAddress,
			"speed":       nic.Speed,
			"unit":        nic.Unit,
		}
		physicalNics = append(physicalNics, physicalNic)
	}
	_ = d.Set("physical_nics", physicalNics)

	// IP addresses information.
	var ipAddresses []map[string]interface{}
	for _, ip := range *host.IpAddresses {
		ipAddress := map[string]interface{}{
			"ip_address": ip.IpAddress,
			"type":       ip.Type,
		}
		ipAddresses = append(ipAddresses, ipAddress)
	}
	_ = d.Set("ip_addresses", ipAddresses)

	return nil
}

func getHostByFqdn(ctx context.Context, apiClient *vcf.ClientWithResponses, fqdn string) (*vcf.Host, error) {
	hostsResponse, err := apiClient.GetHostsWithResponse(ctx, &vcf.GetHostsParams{})
	if err != nil {
		return nil, err
	}

	resp, vcfErr := api_client.GetResponseAs[vcf.PageOfHost](hostsResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}

	for _, hostElement := range *resp.Elements {
		if *hostElement.Fqdn == fqdn {
			return &hostElement, nil
		}
	}

	return nil, fmt.Errorf("host FQDN '%s' not found", fqdn)
}
