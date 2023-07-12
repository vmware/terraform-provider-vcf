/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package vcenter

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
	"strings"
)

// VCSubresourceSchema this helper function extracts the vcenter schema, which
// contains the parameters required to configure Vcenter in a workload domain.
func VCSubresourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vCenter",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "FQDN of the vCenter",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the vCenter virtual machine to be created with the domain",
			},
			"datacenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "vCenter datacenter name",
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ForceNew:     true,
				Description:  "Password for the vCenter root shell user (8-20 characters)",
				ValidateFunc: validation_utils.ValidatePassword,
			},
			"vm_size": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VCenter VM size. One among: xlarge, large, medium, small, tiny",
				ValidateFunc: validation.StringInSlice([]string{
					"xlarge", "large", "medium", "small", "tiny",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"storage_size": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VCenter storage size. One among: lstorage, xlstorage",
				ValidateFunc: validation.StringInSlice([]string{
					"lstorage", "xlstorage",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IPv4 address of the vCenter virtual machine",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"subnet_mask": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Subnet mask",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IPv4 gateway the vCenter VM can use to connect to the outside world",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"dns_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "DNS name of the virtual machine, e.g., vc-1.domain1.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func TryConvertToVcenterSpec(object map[string]interface{}) (*models.VcenterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, name is required")
	}
	datacenterName := object["datacenter_name"].(string)
	if len(datacenterName) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, datacenter_name is required")
	}
	rootPassword := object["root_password"].(string)
	if len(rootPassword) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, root_password is required")
	}
	ipAddress := object["ip_address"].(string)
	if len(ipAddress) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, ip_address is required")
	}
	subnetMask := object["subnet_mask"].(string)
	if len(subnetMask) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, subnet_mask is required")
	}
	gateway := object["gateway"].(string)
	if len(gateway) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, gateway is required")
	}
	dnsName := object["dns_name"].(string)
	if len(dnsName) == 0 {
		return nil, fmt.Errorf("cannot conver to VcenterSpec, dns_name is required")
	}
	vcenterStorageSize := object["vcenter_storage_size"].(string)
	vcenterVmSize := object["vcenter_vm_size"].(string)
	networkDetailsSpec := new(models.NetworkDetailsSpec)
	networkDetailsSpec.IPAddress = &ipAddress
	networkDetailsSpec.SubnetMask = subnetMask
	networkDetailsSpec.DNSName = dnsName
	networkDetailsSpec.Gateway = gateway

	return &models.VcenterSpec{
		DatacenterName:     &datacenterName,
		Name:               &name,
		RootPassword:       &rootPassword,
		StorageSize:        vcenterStorageSize,
		VMSize:             vcenterVmSize,
		NetworkDetailsSpec: networkDetailsSpec,
	}, nil
}
