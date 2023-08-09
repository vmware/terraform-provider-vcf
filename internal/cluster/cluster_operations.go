/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package cluster

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/models"
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

	// TODO support vSAN stretch/unstretch operations by adding a "witness" attribute to vcf_cluster and checking for change on it.
	if data.HasChange("host") {
		oldHostsValue, newHostsValue := data.GetChange("host")
		resultUpdated, err := SetExpansionOrContractionSpec(result,
			oldHostsValue.([]interface{}), newHostsValue.([]interface{}))
		if err != nil {
			return nil, err
		}
		return resultUpdated, nil
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

func ValidateClusterUpdateOperation(ctx context.Context, clusterId string,
	clusterUpdateSpec *models.ClusterUpdateSpec, apiClient *client.VcfClient) diag.Diagnostics {
	validateClusterSpec := clusters.NewValidateClusterOperationsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	validateClusterSpec.ClusterUpdateSpec = clusterUpdateSpec
	validateClusterSpec.ID = clusterId

	validateResponse, err := apiClient.Clusters.ValidateClusterOperations(validateClusterSpec)
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
	intermediaryMap["clusterImageId"] = data.Get("clusterImageId")
	intermediaryMap["evc_mode"] = data.Get("evc_mode")
	intermediaryMap["high_availability_enabled"] = data.Get("high_availability_enabled")
	intermediaryMap["geneve_vlan_id"] = data.Get("geneve_vlan_id")
	intermediaryMap["host"] = data.Get("host")
	intermediaryMap["vds"] = data.Get("vds")
	intermediaryMap["vsan_datastore"] = data.Get("vsan_datastore")
	intermediaryMap["vmfs_datastore"] = data.Get("vmfs_datastore")
	intermediaryMap["vsan_remote_datastore_cluster"] = data.Get("vsan_remote_datastore_cluster")
	intermediaryMap["nfs_datastores"] = data.Get("nfs_datastores")
	intermediaryMap["vvol_datastores"] = data.Get("vvol_datastores")
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
		result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.GeneveVlanID = int32(geneveVlanId.(int))
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

func FlattenCluster(clusterObj *models.Cluster) *map[string]interface{} {
	result := make(map[string]interface{})
	if clusterObj == nil {
		return &result
	}

	result["id"] = clusterObj.ID
	result["name"] = clusterObj.Name
	result["primary_datastore_name"] = clusterObj.PrimaryDatastoreName
	result["primary_datastore_type"] = clusterObj.PrimaryDatastoreType
	result["is_default"] = clusterObj.IsDefault
	result["is_stretched"] = clusterObj.IsStretched

	// TODO typically the VCF 4.5.1 returns only the IDs for the hosts inside models.Cluster
	// consider getting the fqdn, ip and az name with an additional GET request
	flattenedHosts := make([]map[string]interface{}, len(clusterObj.Hosts))
	for j, host := range clusterObj.Hosts {
		flattenedHosts[j] = *FlattenHost(host)
	}
	result["host"] = flattenedHosts

	return &result
}
