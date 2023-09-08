/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/nsxt_clusters"
	"github.com/vmware/vcf-sdk-go/models"
	"sort"
	"strings"
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Form factor for the NSX Manager appliance. One among: large, medium, small",
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
func TryConvertToNsxSpec(object map[string]interface{}) (*models.NsxTSpec, error) {
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

	result := &models.NsxTSpec{}
	result.Vip = &vip
	result.VipFqdn = &vipFqdn
	result.NsxManagerAdminPassword = nsxManagerAdminPassword
	result.LicenseKey = licenseKey

	if formFactor, ok := object["form_factor"]; ok && !validationutils.IsEmpty(formFactor) {
		result.FormFactor = formFactor.(string)
	}

	if nsxManagerAuditPassword, ok := object["nsx_manager_audit_password"]; ok && !validationutils.IsEmpty(nsxManagerAuditPassword) {
		result.NsxManagerAuditPassword = nsxManagerAuditPassword.(string)
	}
	nsxManagerList := object["nsx_manager_node"].([]interface{})
	if len(nsxManagerList) == 0 {
		return nil, fmt.Errorf("cannot convert to NsxTSpec, at least one entry for nsx_manager_node is required")
	}

	var nsxManagerSpecs []*models.NsxManagerSpec
	for _, nsxManagerListEntry := range nsxManagerList {
		nsxManager := nsxManagerListEntry.(map[string]interface{})
		nsxManagerSpec, err := TryConvertToNsxManagerNodeSpec(nsxManager)
		if err != nil {
			return nil, err
		}
		nsxManagerSpecs = append(nsxManagerSpecs, &nsxManagerSpec)
	}
	result.NsxManagerSpecs = nsxManagerSpecs

	return result, nil
}

func FlattenNsxClusterRef(ctx context.Context, nsxtClusterRef *models.NsxTClusterReference,
	apiClient *client.VcfClient) (*[]interface{}, error) {
	flattenedNsxCluster := make(map[string]interface{})
	if nsxtClusterRef == nil {
		return new([]interface{}), nil
	}
	flattenedNsxCluster["id"] = nsxtClusterRef.ID
	flattenedNsxCluster["vip"] = nsxtClusterRef.Vip
	flattenedNsxCluster["vip_fqdn"] = nsxtClusterRef.VipFqdn

	getNsxTClusterParams := nsxt_clusters.NewGetNSXTClusterParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).WithID(nsxtClusterRef.ID)

	nsxtClusterResponse, err := apiClient.NSXTClusters.GetNSXTCluster(getNsxTClusterParams)
	if err != nil {
		return nil, err
	}
	nsxtCluster := nsxtClusterResponse.Payload
	nsxtManagerNodes := nsxtCluster.Nodes
	// Since backend API returns objects in random order sort nsxtManagerNodes list to ensure
	// import is reproducible
	sort.SliceStable(nsxtManagerNodes, func(i, j int) bool {
		return nsxtManagerNodes[i].ID < nsxtManagerNodes[j].ID
	})
	nsxtManagersNodesRaw := *new([]map[string]interface{})
	for _, nsxtManagerNode := range nsxtManagerNodes {
		nsxtManagersNodeRaw := make(map[string]interface{})
		nsxtManagersNodeRaw["name"] = nsxtManagerNode.Name
		nsxtManagersNodeRaw["ip_address"] = nsxtManagerNode.IPAddress
		nsxtManagersNodeRaw["fqdn"] = nsxtManagerNode.Fqdn
		nsxtManagersNodesRaw = append(nsxtManagersNodesRaw, nsxtManagersNodeRaw)
	}

	flattenedNsxCluster["nsx_manager_node"] = nsxtManagersNodesRaw
	result := *new([]interface{})
	result = append(result, flattenedNsxCluster)

	return &result, nil
}
