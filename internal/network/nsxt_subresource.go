/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// NsxTSchema this helper function extracts the NSX-T schema, which
// contains the parameters required to install and configure NSX-T in a workload domain.
func NsxTSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Virtual IP address which would act as proxy/alias for NSX-T Managers",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vip_fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "FQDN for VIP so that common SSL certificates can be installed across all managers",
				ValidateFunc: validation.NoZeroValues,
			},
			"license_key": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "NSX license value",
				ValidateFunc: validation.NoZeroValues,
			},
			"form_factor": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "NSX manager form factor",
				ValidateFunc: validation.NoZeroValues,
			},
			"nsx_manager_admin_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "NSX manager admin password (basic auth and SSH)",
				ValidateFunc: validation_utils.ValidatePassword,
			},
			"nsx_manager_audit_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "NSX manager Audit password",
				ValidateFunc: validation_utils.ValidatePassword,
			},
			"nsx_manager": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification details of the NSX Manager virtual machines. 3 of these are required for the first workload domain",
				Elem:        NsxtManagerSchema(),
			},
		},
	}
}

func TryConvertToNsxtSpec(object map[string]interface{}) (*models.NsxTSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, object is nil")
	}
	vip := object["vip"].(string)
	if len(vip) == 0 {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, vip is required")
	}
	vipFqdn := object["vip_fqdn"].(string)
	if len(vipFqdn) == 0 {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, vip_fqdn is required")
	}
	if object["nsx_manager"] == nil {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, nsx_manager is required")
	}
	nsxManagerAdminPassword := object["nsx_manager_admin_password"].(string)
	if len(nsxManagerAdminPassword) == 0 {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, nsx_manager_admin_password is required")
	}
	licenseKey := object["license_key"].(string)
	if len(licenseKey) == 0 {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, license_key is required")
	}

	result := &models.NsxTSpec{}
	result.Vip = &vip
	result.VipFqdn = &vipFqdn
	result.NsxManagerAdminPassword = nsxManagerAdminPassword
	result.LicenseKey = licenseKey

	if formFactor, ok := object["form_factor"]; ok && !validation_utils.IsEmpty(formFactor) {
		result.FormFactor = formFactor.(string)
	}

	if nsxManagerAuditPassword, ok := object["nsx_manager_audit_password"]; ok && !validation_utils.IsEmpty(nsxManagerAuditPassword) {
		result.NsxManagerAuditPassword = nsxManagerAuditPassword.(string)
	}
	nsxManagerList := object["nsx_manager"].([]interface{})
	if len(nsxManagerList) == 0 {
		return nil, fmt.Errorf("cannot conver to NsxTSpec, at least one entry for nsx_manager is required")
	}

	var nsxManagerSpecs []*models.NsxManagerSpec
	for _, nsxManagerListEntry := range nsxManagerList {
		nsxManager := nsxManagerListEntry.(map[string]interface{})
		nsxManagerSpec, err := TryConvertToNsxManagerSpecs(nsxManager)
		if err != nil {
			return nil, err
		}
		nsxManagerSpecs = append(nsxManagerSpecs, &nsxManagerSpec)
	}
	result.NsxManagerSpecs = nsxManagerSpecs

	return result, nil
}
