/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"
	"strings"
	"time"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the domain",
			},
			"org_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Organization name of the workload domain",
			},
			"vcenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the vCenter virtual machine to be created with the domain",
			},
			"vcenter_datacenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "vCenter datacenter name",
			},
			"vcenter_root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "Password for the vCenter root shell user (8-20 characters)",
				ValidateFunc: validation.StringLenBetween(8, 20),
			},
			"vcenter_vm_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCenter VM size. One among: xlarge, large, medium, small, tiny",
				ValidateFunc: validation.StringInSlice([]string{
					"xlarge", "large", "medium", "small", "tiny",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"vcenter_storage_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCenter storage size. One among: lstorage, xlstorage",
				ValidateFunc: validation.StringInSlice([]string{
					"lstorage", "xlstorage",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"vcenter_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 address of the vCenter virtual machine",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_subnet_mask": {
				Type:         schema.TypeString,
				Required:     false,
				Description:  "Subnet mask",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IPv4 gateway the vCenter VM can use to connect to the outside world",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_dns_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "DNS name of the virtual machine, e.g., vc-1.domain1.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
			"clusters": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification representing the clusters to be added to the workload domain",
				MinItems:    1,
				Elem:        clusterSchema(),
			},
			"nsxt_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification details for NSX-T configuration",
				MaxItems:    1,
				Elem:        network.NsxTSchema(),
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainCreationSpec := createDomainCreationSpec(data)
	validateDomainSpec := domains.NewValidateDomainsOperationsParamsWithContext(ctx)

	validationResult, err := apiClient.Domains.ValidateDomainsOperations(validateDomainSpec)
	if err != nil {
		return nil
	}
	if validation_utils.HasValidationFailed(validationResult.Payload) {
		return validation_utils.ConvertValidationResultToDiag(validationResult.Payload)
	}

	domainCreationParams := domains.NewCreateDomainParamsWithContext(ctx)
	domainCreationParams.DomainCreationSpec = domainCreationSpec

	_, accepted, err := apiClient.Domains.CreateDomain(domainCreationParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := accepted.Payload.ID
	err = vcfClient.WaitForTaskComplete(taskId)
	if err != nil {
		return diag.FromErr(err)
	}
	domainId, err := vcfClient.GetResourceIdAssociatedWithTask(taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(domainId)

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func createDomainCreationSpec(data *schema.ResourceData) *models.DomainCreationSpec {
	result := new(models.DomainCreationSpec)

	return result
}
