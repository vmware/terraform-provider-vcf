// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package cluster

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/client/hosts"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func CreateClusterUpdateSpec(data *schema.ResourceData, markForDeletion bool) (*models.ClusterUpdateSpec, error) {
	result := new(models.ClusterUpdateSpec)
	if markForDeletion {
		result.MarkForDeletion = true
		return result, nil
	}
	if data.HasChange("name") {
		result.Name = data.Get("name").(string)
	}

	if data.HasChange("host") {
		oldHostsValue, newHostsValue := data.GetChange("host")
		return SetExpansionOrContractionSpec(result,
			oldHostsValue.([]interface{}), newHostsValue.([]interface{}))
	}

	if data.HasChange("vsan_stretch_configuration") {
		return SetStretchOrUnstretchSpec(result, data)
	}

	return result, nil
}

// SetExpansionOrContractionSpec sets ClusterExpansionSpec or ClusterContractionSpec to a provided
// ClusterUpdateSpec depending on weather hosts are being added or removed.
func SetExpansionOrContractionSpec(updateSpec *models.ClusterUpdateSpec,
	oldHostsList, newHostsList []interface{}) (*models.ClusterUpdateSpec, error) {

	if len(newHostsList) == len(oldHostsList) {
		return nil, fmt.Errorf("adding and removing hosts is not supported in a single configuration change. Apply each change separately")
	}

	addedHosts, removedHosts := resource_utils.CalculateAddedRemovedResources(newHostsList, oldHostsList)
	if len(removedHosts) == 0 {
		var hostSpecs []*models.HostSpec
		for _, addedHostRaw := range addedHosts {
			hostSpec, err := TryConvertToHostSpec(addedHostRaw)
			if err != nil {
				return nil, err
			}
			hostSpecs = append(hostSpecs, hostSpec)
		}
		clusterExpansionSpec := &models.ClusterExpansionSpec{
			HostSpecs: hostSpecs,
		}
		updateSpec.ClusterExpansionSpec = clusterExpansionSpec
		return updateSpec, nil
	} else {
		var hostRefs []*models.HostReference
		for _, removedHostRaw := range removedHosts {
			hostRef := &models.HostReference{
				ID: removedHostRaw["id"].(string),
			}
			hostRefs = append(hostRefs, hostRef)
		}
		clusterContractionSpec := &models.ClusterCompactionSpec{
			Hosts: hostRefs,
		}
		updateSpec.ClusterCompactionSpec = clusterContractionSpec
		return updateSpec, nil
	}
}

// SetStretchOrUnstretchSpec sets ClusterStretchSpec or ClusterUnstretchSpec to a provided
// ClusterUpdateSpec depending on weather a witness host is being added or removed.
func SetStretchOrUnstretchSpec(updateSpec *models.ClusterUpdateSpec, data *schema.ResourceData) (*models.ClusterUpdateSpec, error) {
	configOld, configNew := data.GetChange("vsan_stretch_configuration")

	if len(configOld.([]interface{})) == len(configNew.([]interface{})) {
		return nil, fmt.Errorf("updating the stretch configuration is not supported")
	}

	configRaw := configNew.([]interface{})

	if len(configRaw) > 0 {
		// stretch
		config := configRaw[0].(map[string]interface{})
		witnessHosts := config["witness_host"].([]interface{})
		witnessHost := witnessHosts[0].(map[string]interface{})

		ip := witnessHost["vsan_ip"].(string)
		cidr := witnessHost["vsan_cidr"].(string)
		fqdn := witnessHost["fqdn"].(string)

		witnessSpec := models.WitnessSpec{
			Fqdn:     &fqdn,
			VSANCidr: &cidr,
			VSANIP:   &ip,
		}

		// All new hosts are added to the secondary fault domain. All existing hosts in the cluster go into the primary domain.
		secondaryFdHosts := config["secondary_fd_host"].([]interface{})
		var hostSpecs []*models.HostSpec
		for _, addedHostRaw := range secondaryFdHosts {
			hostSpec, err := TryConvertToHostSpec(addedHostRaw.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			hostSpecs = append(hostSpecs, hostSpec)
		}

		// MultiAZ support is not yet implemented
		var secondaryAzOverlayVlanId int32 = 0

		stretchSpec := &models.ClusterStretchSpec{
			HostSpecs:                         hostSpecs,
			SecondaryAzOverlayVlanID:          secondaryAzOverlayVlanId,
			WitnessSpec:                       &witnessSpec,
			IsEdgeClusterConfiguredForMultiAZ: false,
		}
		updateSpec.ClusterStretchSpec = stretchSpec
	} else {
		// unstretch
		updateSpec.ClusterUnstretchSpec = EmptySpec{}
	}
	return updateSpec, nil
}

type EmptySpec struct{}

func ValidateClusterUpdateOperation(ctx context.Context, clusterId string,
	clusterUpdateSpec *models.ClusterUpdateSpec, apiClient *client.VcfClient) diag.Diagnostics {
	validateClusterSpec := clusters.NewValidateClusterUpdateSpecParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	validateClusterSpec.ClusterUpdateSpec = clusterUpdateSpec
	validateClusterSpec.ID = clusterId

	validateResponse, err := apiClient.Clusters.ValidateClusterUpdateSpec(validateClusterSpec)
	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}
	if validationUtils.HasValidationFailed(validateResponse.Payload) {
		return validationUtils.ConvertValidationResultToDiag(validateResponse.Payload)
	}
	return nil
}

func TryConvertResourceDataToClusterSpec(data *schema.ResourceData) (*models.ClusterSpec, error) {
	intermediaryMap := map[string]interface{}{}
	intermediaryMap["name"] = data.Get("name")
	intermediaryMap["cluster_image_id"] = data.Get("cluster_image_id")
	intermediaryMap["evc_mode"] = data.Get("evc_mode")
	intermediaryMap["high_availability_enabled"] = data.Get("high_availability_enabled")
	intermediaryMap["geneve_vlan_id"] = data.Get("geneve_vlan_id")
	intermediaryMap["ip_address_pool"] = data.Get("ip_address_pool")
	intermediaryMap["host"] = data.Get("host")
	intermediaryMap["vds"] = data.Get("vds")
	intermediaryMap["vsan_datastore"] = data.Get("vsan_datastore")
	intermediaryMap["vmfs_datastore"] = data.Get("vmfs_datastore")
	intermediaryMap["vsan_remote_datastore_cluster"] = data.Get("vsan_remote_datastore_cluster")
	intermediaryMap["nfs_datastores"] = data.Get("nfs_datastores")
	intermediaryMap["vvol_datastores"] = data.Get("vvol_datastores")
	intermediaryMap["vsan_stretch_configuration"] = data.Get("vsan_stretch_configuration")
	return TryConvertToClusterSpec(intermediaryMap)
}

// TODO implement support for VxRailDetails.

// TryConvertToClusterSpec is a convenience method that converts a map[string]interface{}
// received from the Terraform SDK to an API struct, used in VCF API calls.
func TryConvertToClusterSpec(object map[string]interface{}) (*models.ClusterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, name is required")
	}
	result := &models.ClusterSpec{}
	result.Name = &name
	if clusterImageId, ok := object["cluster_image_id"]; ok && !validationUtils.IsEmpty(clusterImageId) {
		result.ClusterImageID = clusterImageId.(string)
	}
	if evcMode, ok := object["evc_mode"]; ok && len(evcMode.(string)) > 0 {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &models.AdvancedOptions{}
		}
		result.AdvancedOptions.EvcMode = evcMode.(string)
	}
	if highAvailabilityEnabled, ok := object["high_availability_enabled"]; ok && !validationUtils.IsEmpty(highAvailabilityEnabled) {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &models.AdvancedOptions{}
		}
		result.AdvancedOptions.HighAvailability = &models.HighAvailability{
			Enabled: resource_utils.ToBoolPointer(highAvailabilityEnabled),
		}
	}

	result.NetworkSpec = &models.NetworkSpec{}
	result.NetworkSpec.NsxClusterSpec = &models.NsxClusterSpec{}
	result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec = &models.NsxTClusterSpec{}

	if geneveVlanId, ok := object["geneve_vlan_id"]; ok && !validationUtils.IsEmpty(geneveVlanId) {
		vlanValue := int32(geneveVlanId.(int))
		result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.GeneveVlanID = &vlanValue
	}

	if ipAddressPoolRaw, ok := object["ip_address_pool"]; ok && !validationUtils.IsEmpty(ipAddressPoolRaw) {
		ipAddressPoolList := ipAddressPoolRaw.([]interface{})
		if !validationUtils.IsEmpty(ipAddressPoolList[0]) {
			ipAddressPoolSpec, err := network.GetIpAddressPoolSpecFromSchema(ipAddressPoolList[0].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.IPAddressPoolSpec = ipAddressPoolSpec
		}
	}

	if hostsRaw, ok := object["host"]; ok {
		hostsList := hostsRaw.([]interface{})
		if len(hostsList) > 0 {
			result.HostSpecs = []*models.HostSpec{}
			for _, hostListEntry := range hostsList {
				hostSpec, err := TryConvertToHostSpec(hostListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.HostSpecs = append(result.HostSpecs, hostSpec)
			}
		} else {
			return nil, fmt.Errorf("cannot convert to ClusterSpec, hosts list is empty")
		}
	} else {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, hosts list is not set")
	}

	if vdsRaw, ok := object["vds"]; ok {
		vdsList := vdsRaw.([]interface{})
		if len(vdsList) > 0 {
			result.NetworkSpec.VdsSpecs = []*models.VdsSpec{}
			for _, vdsListEntry := range vdsList {
				vdsSpec, err := network.TryConvertToVdsSpec(vdsListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.NetworkSpec.VdsSpecs = append(result.NetworkSpec.VdsSpecs, vdsSpec)
			}
		} else {
			return nil, fmt.Errorf("cannot convert to ClusterSpec, vds list is empty")
		}
	} else {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, vds list is not set")
	}

	datastoreSpec, err := tryConvertToClusterDatastoreSpec(object, name)
	if err != nil {
		return nil, err
	} else {
		result.DatastoreSpec = datastoreSpec
	}

	if stretchConf, ok := object["vsan_stretch_configuration"]; ok && !validationUtils.IsEmpty(stretchConf) {
		return nil, fmt.Errorf("cannot create stretched cluster, create the cluster first and apply the strech configuration later")
	}

	return result, nil
}

func tryConvertToClusterDatastoreSpec(object map[string]interface{}, clusterName string) (*models.DatastoreSpec, error) {
	result := &models.DatastoreSpec{}
	atLeastOneTypeOfDatastoreConfigured := false
	if vsanDatastoreRaw, ok := object["vsan_datastore"]; ok && !validationUtils.IsEmpty(vsanDatastoreRaw) {
		if len(vsanDatastoreRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vsan_datastore config for cluster %q", clusterName)
		}
		vsanDatastoreListEntry := vsanDatastoreRaw.([]interface{})[0]
		vsanDatastoreSpec, err := datastores.TryConvertToVsanDatastoreSpec(vsanDatastoreListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VSANDatastoreSpec = vsanDatastoreSpec
	}
	if vmfsDatastoreRaw, ok := object["vmfs_datastore"]; ok && !validationUtils.IsEmpty(vmfsDatastoreRaw) {
		if len(vmfsDatastoreRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vmfs_datastore config for cluster %q", clusterName)
		}
		vmfsDatastoreListEntry := vmfsDatastoreRaw.([]interface{})[0]
		vmfsDatastoreSpec, err := datastores.TryConvertToVmfsDatastoreSpec(vmfsDatastoreListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VmfsDatastoreSpec = vmfsDatastoreSpec
	}
	if vsanRemoteDatastoreClusterRaw, ok := object["vsan_remote_datastore_cluster"]; ok && !validationUtils.IsEmpty(vsanRemoteDatastoreClusterRaw) {
		if len(vsanRemoteDatastoreClusterRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vsan_remote_datastore_cluster config for cluster %q", clusterName)
		}
		vsanRemoteDatastoreClusterListEntry := vsanRemoteDatastoreClusterRaw.([]interface{})[0]
		vsanRemoteDatastoreClusterSpec, err := datastores.TryConvertToVSANRemoteDatastoreClusterSpec(
			vsanRemoteDatastoreClusterListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VSANRemoteDatastoreClusterSpec = vsanRemoteDatastoreClusterSpec
	}
	if nfsDatastoresRaw, ok := object["nfs_datastores"]; ok && !validationUtils.IsEmpty(nfsDatastoresRaw) {
		nfsDatastoresList := nfsDatastoresRaw.([]interface{})
		if len(nfsDatastoresList) > 0 {
			result.NfsDatastoreSpecs = []*models.NfsDatastoreSpec{}
			for _, nfsDatastoreListEntry := range nfsDatastoresList {
				nfsDatastoreSpec, err := datastores.TryConvertToNfsDatastoreSpec(
					nfsDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.NfsDatastoreSpecs = append(result.NfsDatastoreSpecs, nfsDatastoreSpec)
			}
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if vvolDatastoresRaw, ok := object["vvol_datastores"]; ok && !validationUtils.IsEmpty(vvolDatastoresRaw) {
		vvolDatastoresList := vvolDatastoresRaw.([]interface{})
		if len(vvolDatastoresList) > 0 {
			result.VvolDatastoreSpecs = []*models.VvolDatastoreSpec{}
			for _, vvolDatastoreListEntry := range vvolDatastoresList {
				vvolDatastoreSpec, err := datastores.TryConvertToVvolDatastoreSpec(
					vvolDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.VvolDatastoreSpecs = append(result.VvolDatastoreSpecs, vvolDatastoreSpec)
			}
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if !atLeastOneTypeOfDatastoreConfigured {
		return nil, fmt.Errorf("at least one type of datastore configuration required for cluster %q", clusterName)
	}

	return result, nil
}

func FlattenCluster(ctx context.Context, clusterObj *models.Cluster, apiClient *client.VcfClient) (*map[string]interface{}, error) {
	result := make(map[string]interface{})
	if clusterObj == nil {
		return &result, nil
	}

	result["id"] = clusterObj.ID
	result["name"] = clusterObj.Name
	result["primary_datastore_name"] = clusterObj.PrimaryDatastoreName
	result["primary_datastore_type"] = clusterObj.PrimaryDatastoreType
	result["is_default"] = clusterObj.IsDefault
	result["is_stretched"] = clusterObj.IsStretched

	flattenedVdsSpecs := getFlattenedVdsSpecsForRefs(clusterObj.VdsSpecs)
	result["vds"] = flattenedVdsSpecs

	flattenedHostSpecs, err := getFlattenedHostSpecsForRefs(ctx, clusterObj.Hosts, apiClient)
	if err != nil {
		return nil, err
	}
	result["host"] = flattenedHostSpecs

	return &result, nil
}

func ImportCluster(ctx context.Context, data *schema.ResourceData, apiClient *client.VcfClient,
	clusterId string) ([]*schema.ResourceData, error) {
	getClusterParams := clusters.NewGetClusterParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	getClusterParams.ID = clusterId
	clusterResult, err := apiClient.Clusters.GetCluster(getClusterParams)
	if err != nil {
		return nil, err
	}
	clusterObj := clusterResult.Payload

	data.SetId(clusterObj.ID)
	_ = data.Set("name", clusterObj.Name)
	_ = data.Set("primary_datastore_name", clusterObj.PrimaryDatastoreName)
	_ = data.Set("primary_datastore_type", clusterObj.PrimaryDatastoreType)
	_ = data.Set("is_default", clusterObj.IsDefault)
	_ = data.Set("is_stretched", clusterObj.IsStretched)
	flattenedVdsSpecs := getFlattenedVdsSpecsForRefs(clusterObj.VdsSpecs)
	_ = data.Set("vds", flattenedVdsSpecs)

	flattenedHostSpecs, err := getFlattenedHostSpecsForRefs(ctx, clusterObj.Hosts, apiClient)
	if err != nil {
		return nil, err
	}
	_ = data.Set("host", flattenedHostSpecs)

	//get all domains and find our cluster to set the "domain_id" attribute, because
	// cluster API doesn't provide parent domain ID.
	getDomainsParams := domains.NewGetDomainsParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(ctx)
	domainsResult, err := apiClient.Domains.GetDomains(getDomainsParams)
	if err != nil {
		return nil, err
	}
	allDomains := domainsResult.Payload.Elements
	for _, domain := range allDomains {
		for _, clusterRef := range domain.Clusters {
			if *clusterRef.ID == clusterId {
				_ = data.Set("domain_id", domain.ID)
			}
		}
	}

	return []*schema.ResourceData{data}, nil
}

// getFlattenedHostSpecsForRefs The HostRef is supposed to have all the relevant information,
// but the backend returns everything as nil except the host ID which forces us to make a separate request
// to get some useful info about the hosts in the cluster.
func getFlattenedHostSpecsForRefs(ctx context.Context, hostRefs []*models.HostReference,
	apiClient *client.VcfClient) ([]map[string]interface{}, error) {
	flattenedHostSpecs := *new([]map[string]interface{})
	// Sort for reproducibility
	sort.SliceStable(hostRefs, func(i, j int) bool {
		return hostRefs[i].ID < hostRefs[j].ID
	})
	for _, hostRef := range hostRefs {
		getHostParams := hosts.NewGetHostParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout)
		getHostParams.ID = hostRef.ID
		getHostResult, err := apiClient.Hosts.GetHost(getHostParams)
		if err != nil {
			return nil, err
		}
		hostObj := getHostResult.Payload
		flattenedHostSpecs = append(flattenedHostSpecs, *FlattenHost(hostObj))
	}
	return flattenedHostSpecs, nil
}

func getFlattenedVdsSpecsForRefs(vdsSpecs []*models.VdsSpec) []map[string]interface{} {
	flattenedVdsSpecs := *new([]map[string]interface{})
	// Since backend API returns objects in random order sort VDSSpec list to ensure
	// import is reproducible
	sort.SliceStable(vdsSpecs, func(i, j int) bool {
		return *vdsSpecs[i].Name < *vdsSpecs[j].Name
	})
	for _, vdsSpec := range vdsSpecs {
		flattenedVdsSpecs = append(flattenedVdsSpecs, network.FlattenVdsSpec(vdsSpec))
	}
	return flattenedVdsSpecs
}
