// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package vcenter

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
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
				Description: "ID of the vCenter Server instance",
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified domain name of the vCenter Server instance",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the vCenter Server Appliance virtual machine to be created for the workload domain",
			},
			"datacenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "vSphere datacenter name",
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "root password for the vCenter Server Appliance (8-20 characters)",
				ValidateFunc: validationUtils.ValidatePassword,
			},
			"vm_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vCenter Server instance size. One among: xlarge, large, medium, small, tiny",
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
				Description: "vCenter Server storage size. One among: lstorage, xlstorage",
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
				Description:  "IPv4 address of the vCenter virtual machine",
				ValidateFunc: validationUtils.ValidateIPv4AddressSchema,
			},
			"subnet_mask": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 subnet mask of the vCenter Server instance",
				ValidateFunc: validationUtils.ValidateIPv4AddressSchema,
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 gateway of the vCenter Server instance",
				ValidateFunc: validationUtils.ValidateIPv4AddressSchema,
			},
		},
	}
}

func TryConvertToVcenterSpec(object map[string]interface{}) (*models.VcenterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, name is required")
	}
	datacenterName := object["datacenter_name"].(string)
	if len(datacenterName) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, datacenter_name is required")
	}
	rootPassword := object["root_password"].(string)
	if len(rootPassword) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, root_password is required")
	}
	ipAddress := object["ip_address"].(string)
	if len(ipAddress) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, ip_address is required")
	}
	subnetMask := object["subnet_mask"].(string)
	if len(subnetMask) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, subnet_mask is required")
	}
	gateway := object["gateway"].(string)
	if len(gateway) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, gateway is required")
	}
	fqdn := object["fqdn"].(string)
	if len(fqdn) == 0 {
		return nil, fmt.Errorf("cannot convert to VcenterSpec, fqdn is required")
	}
	vcenterStorageSize, ok := object["storage_size"].(string)
	if !ok {
		vcenterStorageSize = ""
	}
	vcenterVmSize, ok := object["vm_size"].(string)
	if !ok {
		vcenterVmSize = ""
	}
	networkDetailsSpec := new(models.NetworkDetailsSpec)
	networkDetailsSpec.IPAddress = &ipAddress
	networkDetailsSpec.SubnetMask = subnetMask
	networkDetailsSpec.DNSName = fqdn
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
