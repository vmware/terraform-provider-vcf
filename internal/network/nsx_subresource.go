// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// NsxSchema this helper function extracts the NSX schema, which
// contains the parameters required to install and configure NSX in a workload domain.
func NsxSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the NSX Manager cluster",
			},
			"vip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Virtual IP (VIP) for the NSX Manager cluster",
				ValidateFunc: validationutils.ValidateIPv4AddressSchema,
			},
			"vip_fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified domain name of the NSX Manager cluster VIP",
				ValidateFunc: validation.NoZeroValues,
			},
			"license_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "NSX license to be used",
				ValidateFunc: validation.NoZeroValues,
			},
			"form_factor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Form factor for the NSX Manager appliance. One among: large, medium, small",
				ValidateFunc: validation.StringInSlice([]string{
					"large", "medium", "small",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"nsx_manager_admin_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "NSX Manager admin user password",
				ValidateFunc: validationutils.ValidatePassword,
			},
			"nsx_manager_audit_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "NSX Manager audit user password",
				ValidateFunc: validationutils.ValidatePassword,
			},
			"nsx_manager_node": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification details of the NSX Manager virtual machines. 3 of these are required for the first workload domain",
				Elem:        NsxManagerNodeSchema(),
			},
		},
	}
}

// TODO support IpPoolSpecs.

// TryConvertToNsxSpec is a convenience method that converts a map[string]interface{}
// // received from the Terraform SDK to an API struct, used in VCF API calls.
func TryConvertToNsxSpec(object map[string]interface{}) (*vcf.NsxTSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, object is nil")
	}
	vip := object["vip"].(string)
	if len(vip) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, vip is required")
	}
	vipFqdn := object["vip_fqdn"].(string)
	if len(vipFqdn) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, vip_fqdn is required")
	}
	if object["nsx_manager_node"] == nil {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, nsx_manager is required")
	}
	nsxManagerAdminPassword := object["nsx_manager_admin_password"].(string)
	if len(nsxManagerAdminPassword) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, nsx_manager_admin_password is required")
	}
	licenseKey := object["license_key"].(string)
	if len(licenseKey) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, license_key is required")
	}

	result := &vcf.NsxTSpec{}
	result.Vip = &vip
	result.VipFqdn = vipFqdn
	result.NsxManagerAdminPassword = &nsxManagerAdminPassword
	result.LicenseKey = &licenseKey

	if formFactor, ok := object["form_factor"]; ok && !validationutils.IsEmpty(formFactor) {
		result.FormFactor = resource_utils.ToStringPointer(formFactor)
	}

	if nsxManagerAuditPassword, ok := object["nsx_manager_audit_password"]; ok && !validationutils.IsEmpty(nsxManagerAuditPassword) {
		result.NsxManagerAuditPassword = resource_utils.ToStringPointer(nsxManagerAuditPassword)
	}
	nsxManagerList := object["nsx_manager_node"].([]interface{})
	if len(nsxManagerList) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, at least one entry for nsx_manager_node is required")
	}

	var nsxManagerSpecs []vcf.NsxManagerSpec
	for _, nsxManagerListEntry := range nsxManagerList {
		nsxManager := nsxManagerListEntry.(map[string]interface{})
		nsxManagerSpec, err := TryConvertToNsxManagerNodeSpec(nsxManager)
		if err != nil {
			return nil, err
		}
		nsxManagerSpecs = append(nsxManagerSpecs, nsxManagerSpec)
	}
	result.NsxManagerSpecs = nsxManagerSpecs

	return result, nil
}

func FlattenNsxClusterRef(ctx context.Context, nsxtClusterRef vcf.NsxTClusterReference,
	apiClient *vcf.ClientWithResponses) (*[]interface{}, error) {
	flattenedNsxCluster := make(map[string]interface{})
	flattenedNsxCluster["id"] = nsxtClusterRef.Id
	flattenedNsxCluster["vip"] = nsxtClusterRef.Vip
	flattenedNsxCluster["vip_fqdn"] = nsxtClusterRef.VipFqdn

	res, err := apiClient.GetNsxClusterWithResponse(ctx, *nsxtClusterRef.Id)
	if err != nil {
		return nil, err
	}
	nsxtCluster, vcfErr := api_client.GetResponseAs[vcf.NsxTCluster](res)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}
	if nsxtCluster.Nodes != nil {
		nsxtManagerNodes := *nsxtCluster.Nodes
		// Since backend API returns objects in random order sort nsxtManagerNodes list to ensure
		// import is reproducible
		sort.SliceStable(nsxtManagerNodes, func(i, j int) bool {
			return *(nsxtManagerNodes[i].Id) < *(nsxtManagerNodes[j].Id)
		})
		nsxtManagersNodesRaw := *new([]map[string]interface{})
		for _, nsxtManagerNode := range nsxtManagerNodes {
			nsxtManagersNodeRaw := make(map[string]interface{})
			nsxtManagersNodeRaw["name"] = nsxtManagerNode.Name
			nsxtManagersNodeRaw["ip_address"] = nsxtManagerNode.IpAddress
			nsxtManagersNodeRaw["fqdn"] = nsxtManagerNode.Fqdn
			nsxtManagersNodesRaw = append(nsxtManagersNodesRaw, nsxtManagersNodeRaw)
		}
		flattenedNsxCluster["nsx_manager_node"] = nsxtManagersNodesRaw
	}

	result := *new([]interface{})
	result = append(result, flattenedNsxCluster)

	return &result, nil
}
