// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/terraform-provider-vcf/internal/vcenter"
)

func CreateDomainCreationSpec(data *schema.ResourceData) (*models.DomainCreationSpec, error) {
	result := new(models.DomainCreationSpec)
	domainName := data.Get("name").(string)
	result.DomainName = &domainName

	if orgName, ok := data.GetOk("org_name"); ok {
		result.OrgName = orgName.(string)
	}

	vcenterSpec, err := generateVcenterSpecFromResourceData(data)
	if err == nil {
		result.VcenterSpec = vcenterSpec
	} else {
		return nil, err
	}

	nsxSpec, err := generateNsxSpecFromResourceData(data)
	if err == nil {
		result.NsxTSpec = nsxSpec
	} else {
		return nil, err
	}

	computeSpec, err := generateComputeSpecFromResourceData(data)
	if err == nil {
		result.ComputeSpec = computeSpec
	} else {
		return nil, err
	}

	return result, nil
}

func ReadAndSetClustersDataToDomainResource(domainClusterRefs []*models.ClusterReference,
	ctx context.Context, data *schema.ResourceData, apiClient *client.VcfClient) error {
	clusterIdsInTheCurrentDomain := make(map[string]bool, len(domainClusterRefs))
	for _, clusterReference := range domainClusterRefs {
		clusterIdsInTheCurrentDomain[*clusterReference.ID] = true
	}

	getClustersParams := clusters.NewGetClustersParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)

	clustersResult, err := apiClient.Clusters.GetClusters(getClustersParams)
	if err != nil {
		return err
	}
	domainClusterData := data.Get("cluster")
	domainClusterDataList := domainClusterData.([]interface{})
	allClusters := clustersResult.Payload.Elements
	for _, domainClusterRaw := range domainClusterDataList {
		domainCluster := domainClusterRaw.(map[string]interface{})
		for _, clusterObj := range allClusters {
			_, ok := clusterIdsInTheCurrentDomain[clusterObj.ID]
			// go over clusters that are in the domain, skip the rest
			if !ok {
				continue
			}
			if domainCluster["name"] == clusterObj.Name {
				domainCluster["id"] = clusterObj.ID
				domainCluster["primary_datastore_name"] = clusterObj.PrimaryDatastoreName
				domainCluster["primary_datastore_type"] = clusterObj.PrimaryDatastoreType
				domainCluster["is_default"] = clusterObj.IsDefault
				domainCluster["is_stretched"] = clusterObj.IsStretched
			}
		}
	}
	_ = data.Set("cluster", domainClusterData)

	return nil
}

func SetBasicDomainAttributes(ctx context.Context, domainId string, data *schema.ResourceData,
	apiClient *client.VcfClient) (*models.Domain, error) {
	getDomainParams := domains.NewGetDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	getDomainParams.ID = domainId
	domainResult, err := apiClient.Domains.GetDomain(getDomainParams)
	if err != nil {
		return nil, err
	}
	domain := domainResult.Payload

	data.SetId(domain.ID)
	_ = data.Set("name", domain.Name)
	_ = data.Set("status", domain.Status)
	_ = data.Set("type", domain.Type)
	_ = data.Set("sso_id", domain.SSOID)
	_ = data.Set("sso_name", domain.SSOName)
	_ = data.Set("is_management_sso_domain", domain.IsManagementSSODomain)
	if len(domain.VCENTERS) < 1 {
		return nil, fmt.Errorf("no vCenter Server instance found for domain %q", domainId)
	}
	vcenterConfigAttribute, vcenterConfigExists := data.GetOk("vcenter_configuration")
	var vcenterConfigRaw []interface{}
	if vcenterConfigExists {
		vcenterConfigRaw = vcenterConfigAttribute.([]interface{})
	} else {
		vcenterConfigRaw = *new([]interface{})
		vcenterConfigRaw = append(vcenterConfigRaw, make(map[string]interface{}))
	}
	vcenterConfig := vcenterConfigRaw[0].(map[string]interface{})
	vcenterConfig["id"] = domain.VCENTERS[0].ID
	vcenterConfig["fqdn"] = domain.VCENTERS[0].Fqdn
	_ = data.Set("vcenter_configuration", vcenterConfigRaw)

	return domain, nil
}

func CreateDomainUpdateSpec(data *schema.ResourceData, markForDeletion bool) *models.DomainUpdateSpec {
	result := new(models.DomainUpdateSpec)
	if markForDeletion {
		result.MarkForDeletion = true
		return result
	}
	if data.HasChange("name") {
		result.Name = data.Get("name").(string)
	}

	// TODO implement support for IPPoolSpecs in NsxTSpec
	// by placing the added cluster spec in the DomainUpdateSpec
	//nsxtSpec, err := generateNsxSpecFromResourceData(data)
	//if err == nil {
	//	result.NsxTSpec = nsxtSpec
	//} else {
	//	return nil, err
	//}

	return result
}

func ImportDomain(ctx context.Context, data *schema.ResourceData, apiClient *client.VcfClient,
	domainId string, allowManagementDomain bool) ([]*schema.ResourceData, error) {
	domainObj, err := SetBasicDomainAttributes(ctx, domainId, data, apiClient)
	if err != nil {
		return nil, err
	}
	if !allowManagementDomain && data.Get("type").(string) == "MANAGEMENT" {
		return nil, fmt.Errorf("domain %s cannot be imported as it is management domain", domainId)
	}

	err = setClustersDataToDomainDataSource(domainObj.Clusters, ctx, data, apiClient)
	if err != nil {
		return nil, err
	}
	flattenedNsxClusterRef, err := network.FlattenNsxClusterRef(ctx, domainObj.NSXTCluster, apiClient)
	if err != nil {
		return nil, err
	}
	_ = data.Set("nsx_configuration", *flattenedNsxClusterRef)
	return []*schema.ResourceData{data}, nil
}

func setClustersDataToDomainDataSource(domainClusterRefs []*models.ClusterReference, ctx context.Context, data *schema.ResourceData, apiClient *client.VcfClient) error {
	clusterIds := make([]string, len(domainClusterRefs))
	for i, clusterReference := range domainClusterRefs {
		clusterIds[i] = *clusterReference.ID
	}
	// Sort the id slice, to have a deterministic order in every run of the domain datasource read
	sort.Strings(clusterIds)

	flattenedClusters := make([]map[string]interface{}, len(domainClusterRefs))
	for i, clusterId := range clusterIds {
		getClusterParams := clusters.GetClusterParams{ID: clusterId}
		getClusterParams.WithContext(ctx).WithTimeout(constants.DefaultVcfApiCallTimeout)
		clusterResult, err := apiClient.Clusters.GetCluster(&getClusterParams)
		if err != nil {
			return err
		}
		clusterRef := clusterResult.Payload
		flattenedCluster, err := cluster.FlattenCluster(ctx, clusterRef, apiClient)
		if err != nil {
			return err
		}
		flattenedClusters[i] = *flattenedCluster

	}
	_ = data.Set("cluster", flattenedClusters)

	return nil
}

func generateNsxSpecFromResourceData(data *schema.ResourceData) (*models.NsxTSpec, error) {
	if nsxConfigRaw, ok := data.GetOk("nsx_configuration"); ok && len(nsxConfigRaw.([]interface{})) > 0 {
		nsxConfigList := nsxConfigRaw.([]interface{})
		nsxConfigListEntry := nsxConfigList[0].(map[string]interface{})
		nsxSpec, err := network.TryConvertToNsxSpec(nsxConfigListEntry)
		return nsxSpec, err
	}
	return nil, nil
}

func generateVcenterSpecFromResourceData(data *schema.ResourceData) (*models.VcenterSpec, error) {
	if vcenterConfigRaw, ok := data.GetOk("vcenter_configuration"); ok && len(vcenterConfigRaw.([]interface{})) > 0 {
		vcenterConfigList := vcenterConfigRaw.([]interface{})
		vcenterConfigListEntry := vcenterConfigList[0].(map[string]interface{})
		vcenterSpec, err := vcenter.TryConvertToVcenterSpec(vcenterConfigListEntry)
		return vcenterSpec, err
	}
	return nil, nil
}

func generateComputeSpecFromResourceData(data *schema.ResourceData) (*models.ComputeSpec, error) {
	if clusterConfigRaw, ok := data.GetOk("cluster"); ok && !validationUtils.IsEmpty(clusterConfigRaw) {
		clusterConfigList := clusterConfigRaw.([]interface{})
		result := new(models.ComputeSpec)
		var clusterSpecs []*models.ClusterSpec
		for _, clusterConfigListEntry := range clusterConfigList {
			clusterSpec, err := cluster.TryConvertToClusterSpec(clusterConfigListEntry.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			clusterSpecs = append(clusterSpecs, clusterSpec)
		}
		result.ClusterSpecs = clusterSpecs
		return result, nil
	}
	return nil, fmt.Errorf("no cluster configuration")
}
