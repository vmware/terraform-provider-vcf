// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/terraform-provider-vcf/internal/vcenter"
)

func CreateDomainCreationSpec(data *schema.ResourceData) (*vcf.DomainCreationSpec, error) {
	result := &vcf.DomainCreationSpec{}
	domainName := data.Get("name").(string)
	result.DomainName = &domainName

	if orgName, ok := data.GetOk("org_name"); ok {
		result.OrgName = utils.ToStringPointer(orgName)
	}

	vcenterSpec, err := generateVcenterSpecFromResourceData(data)
	if err == nil {
		result.VcenterSpec = *vcenterSpec
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
		result.ComputeSpec = *computeSpec
	} else {
		return nil, err
	}

	return result, nil
}

func ReadAndSetClustersDataToDomainResource(domainClusterRefs []vcf.ClusterReference,
	ctx context.Context, data *schema.ResourceData, apiClient *vcf.ClientWithResponses) error {
	clusterIdsInTheCurrentDomain := make(map[string]bool, len(domainClusterRefs))
	for _, clusterReference := range domainClusterRefs {
		clusterIdsInTheCurrentDomain[clusterReference.Id] = true
	}

	clustersResult, err := apiClient.GetClustersWithResponse(ctx, nil)
	if err != nil {
		return err
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfCluster](clustersResult)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return errors.New(*vcfErr.Message)
	}
	domainClusterData := data.Get("cluster")
	domainClusterDataList := domainClusterData.([]interface{})
	allClusters := page.Elements
	for _, domainClusterRaw := range domainClusterDataList {
		domainCluster := domainClusterRaw.(map[string]interface{})
		if allClusters != nil {
			for _, clusterObj := range *allClusters {
				_, ok := clusterIdsInTheCurrentDomain[*clusterObj.Id]
				// go over clusters that are in the domain, skip the rest
				if !ok {
					continue
				}
				if domainCluster["name"] == *clusterObj.Name {
					domainCluster["id"] = *clusterObj.Id
					domainCluster["primary_datastore_name"] = *clusterObj.PrimaryDatastoreName
					domainCluster["primary_datastore_type"] = *clusterObj.PrimaryDatastoreType
					domainCluster["is_default"] = *clusterObj.IsDefault
					domainCluster["is_stretched"] = *clusterObj.IsStretched
				}
			}
		}
	}
	_ = data.Set("cluster", domainClusterData)

	return nil
}

func SetBasicDomainAttributes(ctx context.Context, domainId string, data *schema.ResourceData,
	apiClient *vcf.ClientWithResponses) (*vcf.Domain, error) {
	domainRes, err := apiClient.GetDomainWithResponse(ctx, domainId)
	if err != nil {
		return nil, err
	}
	domain, vcfErr := api_client.GetResponseAs[vcf.Domain](domainRes)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}

	data.SetId(*domain.Id)
	_ = data.Set("name", *domain.Name)
	_ = data.Set("status", *domain.Status)
	_ = data.Set("type", *domain.Type)
	_ = data.Set("sso_id", *domain.SsoId)
	_ = data.Set("sso_name", *domain.SsoName)
	_ = data.Set("is_management_sso_domain", *domain.IsManagementSsoDomain)
	if domain.Vcenters == nil || len(*domain.Vcenters) < 1 {
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
	if domain.Vcenters != nil {
		vcenterConfig["id"] = (*domain.Vcenters)[0].Id
		vcenterConfig["fqdn"] = (*domain.Vcenters)[0].Fqdn
	}
	_ = data.Set("vcenter_configuration", vcenterConfigRaw)

	return domain, nil
}

func CreateDomainUpdateSpec(data *schema.ResourceData, markForDeletion bool) vcf.DomainUpdateSpec {
	result := vcf.DomainUpdateSpec{}
	if markForDeletion {
		result.MarkForDeletion = utils.ToBoolPointer(true)
		return result
	}
	if data.HasChange("name") {
		result.Name = utils.ToStringPointer(data.Get("name"))
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

func ImportDomain(ctx context.Context, data *schema.ResourceData, apiClient *vcf.ClientWithResponses,
	domainId string, allowManagementDomain bool) ([]*schema.ResourceData, error) {
	domainObj, err := SetBasicDomainAttributes(ctx, domainId, data, apiClient)
	if err != nil {
		return nil, err
	}
	if !allowManagementDomain && data.Get("type").(string) == "MANAGEMENT" {
		return nil, fmt.Errorf("domain %s cannot be imported as it is management domain", domainId)
	}

	if domainObj.Clusters != nil {
		err = setClustersDataToDomainDataSource(*domainObj.Clusters, ctx, data, apiClient)
		if err != nil {
			return nil, err
		}
	}

	if domainObj.NsxtCluster != nil {
		flattenedNsxClusterRef, err := network.FlattenNsxClusterRef(ctx, *domainObj.NsxtCluster, apiClient)
		if err != nil {
			return nil, err
		}
		_ = data.Set("nsx_configuration", *flattenedNsxClusterRef)
	}

	return []*schema.ResourceData{data}, nil
}

func setClustersDataToDomainDataSource(domainClusterRefs []vcf.ClusterReference, ctx context.Context,
	data *schema.ResourceData, apiClient *vcf.ClientWithResponses) error {
	clusterIds := make([]string, len(domainClusterRefs))
	for i, clusterReference := range domainClusterRefs {
		clusterIds[i] = clusterReference.Id
	}
	// Sort the id slice, to have a deterministic order in every run of the domain datasource read
	sort.Strings(clusterIds)

	flattenedClusters := make([]map[string]interface{}, len(domainClusterRefs))
	for i, clusterId := range clusterIds {
		res, err := apiClient.GetClusterWithResponse(ctx, clusterId)
		if err != nil {
			return err
		}
		clusterRef, vcfErr := api_client.GetResponseAs[vcf.Cluster](res)
		if vcfErr != nil {
			api_client.LogError(vcfErr)
			return errors.New(*vcfErr.Message)
		}
		flattenedCluster, err := cluster.FlattenCluster(ctx, clusterRef, apiClient)
		if err != nil {
			return err
		}
		flattenedClusters[i] = *flattenedCluster

	}
	_ = data.Set("cluster", flattenedClusters)

	return nil
}

func generateNsxSpecFromResourceData(data *schema.ResourceData) (*vcf.NsxTSpec, error) {
	if nsxConfigRaw, ok := data.GetOk("nsx_configuration"); ok && len(nsxConfigRaw.([]interface{})) > 0 {
		nsxConfigList := nsxConfigRaw.([]interface{})
		nsxConfigListEntry := nsxConfigList[0].(map[string]interface{})
		nsxSpec, err := network.TryConvertToNsxSpec(nsxConfigListEntry)
		return nsxSpec, err
	}
	return nil, nil
}

func generateVcenterSpecFromResourceData(data *schema.ResourceData) (*vcf.VcenterSpec, error) {
	if vcenterConfigRaw, ok := data.GetOk("vcenter_configuration"); ok && len(vcenterConfigRaw.([]interface{})) > 0 {
		vcenterConfigList := vcenterConfigRaw.([]interface{})
		vcenterConfigListEntry := vcenterConfigList[0].(map[string]interface{})
		vcenterSpec, err := vcenter.TryConvertToVcenterSpec(vcenterConfigListEntry)
		return vcenterSpec, err
	}
	return nil, nil
}

func generateComputeSpecFromResourceData(data *schema.ResourceData) (*vcf.ComputeSpec, error) {
	if clusterConfigRaw, ok := data.GetOk("cluster"); ok && !validationUtils.IsEmpty(clusterConfigRaw) {
		clusterConfigList := clusterConfigRaw.([]interface{})
		result := &vcf.ComputeSpec{}
		var clusterSpecs []vcf.ClusterSpec
		for _, clusterConfigListEntry := range clusterConfigList {
			clusterSpec, err := cluster.TryConvertToClusterSpec(clusterConfigListEntry.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			clusterSpecs = append(clusterSpecs, *clusterSpec)
		}
		result.ClusterSpecs = clusterSpecs
		return result, nil
	}
	return nil, fmt.Errorf("no cluster configuration")
}
