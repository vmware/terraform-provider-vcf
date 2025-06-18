// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package cluster

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func CreateClusterUpdateSpec(data *schema.ResourceData, markForDeletion bool) (*vcf.ClusterUpdateSpec, error) {
	result := new(vcf.ClusterUpdateSpec)
	if markForDeletion {
		result.MarkForDeletion = &markForDeletion
		return result, nil
	}
	if data.HasChange("name") {
		result.Name = utils.ToStringPointer(data.Get("name"))
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
func SetExpansionOrContractionSpec(updateSpec *vcf.ClusterUpdateSpec,
	oldHostsList, newHostsList []interface{}) (*vcf.ClusterUpdateSpec, error) {

	if len(newHostsList) == len(oldHostsList) {
		return nil, fmt.Errorf("adding and removing hosts is not supported in a single configuration change. Apply each change separately")
	}

	addedHosts, removedHosts := utils.CalculateAddedRemovedResources(newHostsList, oldHostsList)
	if len(removedHosts) == 0 {
		var hostSpecs []vcf.HostSpec
		for _, addedHostRaw := range addedHosts {
			hostSpec, err := TryConvertToHostSpec(addedHostRaw)
			if err != nil {
				return nil, err
			}
			hostSpecs = append(hostSpecs, *hostSpec)
		}
		clusterExpansionSpec := &vcf.ClusterExpansionSpec{
			DeployWithoutLicenseKeys: utils.ToPointer[bool](true),
			HostSpecs:                hostSpecs,
		}
		updateSpec.ClusterExpansionSpec = clusterExpansionSpec
		return updateSpec, nil
	} else {
		var hostRefs []vcf.HostReference
		for _, removedHostRaw := range removedHosts {
			hostRef := vcf.HostReference{
				Id: utils.ToStringPointer(removedHostRaw["id"]),
			}
			hostRefs = append(hostRefs, hostRef)
		}
		clusterContractionSpec := &vcf.ClusterCompactionSpec{
			Hosts: hostRefs,
		}
		updateSpec.ClusterCompactionSpec = clusterContractionSpec
		return updateSpec, nil
	}
}

// SetStretchOrUnstretchSpec sets ClusterStretchSpec or ClusterUnstretchSpec to a provided
// ClusterUpdateSpec depending on weather a witness host is being added or removed.
func SetStretchOrUnstretchSpec(updateSpec *vcf.ClusterUpdateSpec, data *schema.ResourceData) (*vcf.ClusterUpdateSpec, error) {
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

		witnessSpec := vcf.WitnessSpec{
			Fqdn:     fqdn,
			VsanCidr: cidr,
			VsanIp:   ip,
		}

		// All new hosts are added to the secondary fault domain. All existing hosts in the cluster go into the primary domain.
		secondaryFdHosts := config["secondary_fd_host"].([]interface{})
		var hostSpecs []vcf.HostSpec
		for _, addedHostRaw := range secondaryFdHosts {
			hostSpec, err := TryConvertToHostSpec(addedHostRaw.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			hostSpecs = append(hostSpecs, *hostSpec)
		}

		// MultiAZ support is not yet implemented
		var secondaryAzOverlayVlanId int32 = 0

		stretchSpec := &vcf.ClusterStretchSpec{
			DeployWithoutLicenseKeys:          utils.ToPointer[bool](true),
			HostSpecs:                         hostSpecs,
			SecondaryAzOverlayVlanId:          &secondaryAzOverlayVlanId,
			WitnessSpec:                       witnessSpec,
			IsEdgeClusterConfiguredForMultiAZ: utils.ToBoolPointer(false),
		}
		updateSpec.ClusterStretchSpec = stretchSpec
	} else {
		// unstretch
		updateSpec.ClusterUnstretchSpec = &vcf.ClusterUnstretchSpec{}
	}
	return updateSpec, nil
}

func ValidateClusterUpdateOperation(ctx context.Context, clusterId string,
	clusterUpdateSpec vcf.ClusterUpdateSpec, apiClient *vcf.ClientWithResponses) diag.Diagnostics {
	validateResponse, err := apiClient.ValidateClusterUpdateSpecWithResponse(ctx, clusterId, nil, clusterUpdateSpec)
	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}
	validationResult, vcfErr := api_client.GetResponseAs[vcf.Validation](validateResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if validationUtils.HasValidationFailed(validationResult) {
		return validationUtils.ConvertValidationResultToDiag(validationResult)
	}
	return nil
}

func TryConvertResourceDataToClusterSpec(data *schema.ResourceData) (*vcf.ClusterSpec, error) {
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

// TryConvertToClusterSpec is a convenience method that converts a map[string]interface{}
// received from the Terraform SDK to an API struct, used in VCF API calls.
func TryConvertToClusterSpec(object map[string]interface{}) (*vcf.ClusterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, name is required")
	}
	result := &vcf.ClusterSpec{}
	result.Name = &name
	if clusterImageId, ok := object["cluster_image_id"]; ok && !validationUtils.IsEmpty(clusterImageId) {
		result.ClusterImageId = utils.ToStringPointer(clusterImageId)
	}
	if evcMode, ok := object["evc_mode"]; ok && len(evcMode.(string)) > 0 {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &vcf.AdvancedOptions{}
		}
		result.AdvancedOptions.EvcMode = utils.ToStringPointer(evcMode)
	}
	if highAvailabilityEnabled, ok := object["high_availability_enabled"]; ok && !validationUtils.IsEmpty(highAvailabilityEnabled) {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &vcf.AdvancedOptions{}
		}
		result.AdvancedOptions.HighAvailability = &vcf.HighAvailability{
			Enabled: highAvailabilityEnabled.(bool),
		}
	}

	result.NetworkSpec = vcf.NetworkSpec{}
	result.NetworkSpec.NsxClusterSpec = &vcf.NsxClusterSpec{}
	result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec = &vcf.NsxTClusterSpec{}

	if geneveVlanId, ok := object["geneve_vlan_id"]; ok && !validationUtils.IsEmpty(geneveVlanId) {
		vlanValue := int32(geneveVlanId.(int))
		result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.GeneveVlanId = &vlanValue
	}

	if ipAddressPoolRaw, ok := object["ip_address_pool"]; ok && !validationUtils.IsEmpty(ipAddressPoolRaw) {
		ipAddressPoolList := ipAddressPoolRaw.([]interface{})
		if !validationUtils.IsEmpty(ipAddressPoolList[0]) {
			ipAddressPoolSpec, err := network.GetIpAddressPoolSpecFromSchema(ipAddressPoolList[0].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.IpAddressPoolSpec = ipAddressPoolSpec
		}
	}

	if hostsRaw, ok := object["host"]; ok {
		hostsList := hostsRaw.([]interface{})
		if len(hostsList) > 0 {
			result.HostSpecs = []vcf.HostSpec{}
			for _, hostListEntry := range hostsList {
				hostSpec, err := TryConvertToHostSpec(hostListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.HostSpecs = append(result.HostSpecs, *hostSpec)
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
			vdsSpecs := []vcf.VdsSpec{}
			for _, vdsListEntry := range vdsList {
				vdsSpec, err := network.TryConvertToVdsSpec(vdsListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				vdsSpecs = append(vdsSpecs, *vdsSpec)
			}
			result.NetworkSpec.VdsSpecs = &vdsSpecs
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
		result.DatastoreSpec = *datastoreSpec
	}

	if stretchConf, ok := object["vsan_stretch_configuration"]; ok && !validationUtils.IsEmpty(stretchConf) {
		return nil, fmt.Errorf("cannot create stretched cluster, create the cluster first and apply the strech configuration later")
	}

	return result, nil
}

func tryConvertToClusterDatastoreSpec(object map[string]interface{}, clusterName string) (*vcf.DatastoreSpec, error) {
	result := &vcf.DatastoreSpec{}
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
		result.VsanDatastoreSpec = vsanDatastoreSpec
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
		result.VsanRemoteDatastoreClusterSpec = vsanRemoteDatastoreClusterSpec
	}
	if nfsDatastoresRaw, ok := object["nfs_datastores"]; ok && !validationUtils.IsEmpty(nfsDatastoresRaw) {
		nfsDatastoresList := nfsDatastoresRaw.([]interface{})
		if len(nfsDatastoresList) > 0 {
			specs := []vcf.NfsDatastoreSpec{}
			for _, nfsDatastoreListEntry := range nfsDatastoresList {
				nfsDatastoreSpec, err := datastores.TryConvertToNfsDatastoreSpec(
					nfsDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				specs = append(specs, *nfsDatastoreSpec)
			}
			result.NfsDatastoreSpecs = &specs
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if vvolDatastoresRaw, ok := object["vvol_datastores"]; ok && !validationUtils.IsEmpty(vvolDatastoresRaw) {
		vvolDatastoresList := vvolDatastoresRaw.([]interface{})
		if len(vvolDatastoresList) > 0 {
			specs := []vcf.VvolDatastoreSpec{}
			for _, vvolDatastoreListEntry := range vvolDatastoresList {
				vvolDatastoreSpec, err := datastores.TryConvertToVvolDatastoreSpec(
					vvolDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				specs = append(specs, *vvolDatastoreSpec)
			}
			result.VvolDatastoreSpecs = &specs
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if !atLeastOneTypeOfDatastoreConfigured {
		return nil, fmt.Errorf("at least one type of datastore configuration required for cluster %q", clusterName)
	}

	return result, nil
}

func FlattenCluster(ctx context.Context, clusterObj *vcf.Cluster, apiClient *vcf.ClientWithResponses) (*map[string]interface{}, error) {
	result := make(map[string]interface{})
	if clusterObj == nil {
		return &result, nil
	}

	result["id"] = clusterObj.Id
	result["name"] = clusterObj.Name
	result["primary_datastore_name"] = clusterObj.PrimaryDatastoreName
	result["primary_datastore_type"] = clusterObj.PrimaryDatastoreType
	result["is_default"] = clusterObj.IsDefault
	result["is_stretched"] = clusterObj.IsStretched

	if clusterObj.VdsSpecs != nil {
		flattenedVdsSpecs := getFlattenedVdsSpecsForRefs(*clusterObj.VdsSpecs)
		result["vds"] = flattenedVdsSpecs
	}

	if clusterObj.Hosts != nil {
		flattenedHostSpecs, err := getFlattenedHostSpecsForRefs(ctx, *clusterObj.Hosts, apiClient)
		if err != nil {
			return nil, err
		}
		result["host"] = flattenedHostSpecs
	}

	return &result, nil
}

func ImportCluster(ctx context.Context, data *schema.ResourceData, apiClient *vcf.ClientWithResponses,
	clusterId string) ([]*schema.ResourceData, error) {
	clusterRes, err := apiClient.GetClusterWithResponse(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	clusterObj, vcfErr := api_client.GetResponseAs[vcf.Cluster](clusterRes)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}

	data.SetId(*clusterObj.Id)
	_ = data.Set("name", clusterObj.Name)
	_ = data.Set("primary_datastore_name", clusterObj.PrimaryDatastoreName)
	_ = data.Set("primary_datastore_type", clusterObj.PrimaryDatastoreType)
	_ = data.Set("is_default", clusterObj.IsDefault)
	_ = data.Set("is_stretched", clusterObj.IsStretched)
	if clusterObj.VdsSpecs != nil {
		flattenedVdsSpecs := getFlattenedVdsSpecsForRefs(*clusterObj.VdsSpecs)
		_ = data.Set("vds", flattenedVdsSpecs)
	}

	if clusterObj.Hosts != nil {
		flattenedHostSpecs, err := getFlattenedHostSpecsForRefs(ctx, *clusterObj.Hosts, apiClient)
		if err != nil {
			return nil, err
		}
		_ = data.Set("host", flattenedHostSpecs)
	}

	// get all domains and find our cluster to set the "domain_id" attribute, because
	// cluster API doesn't provide parent domain ID.
	getDomainsParams := &vcf.GetDomainsParams{}
	domainsRes, err := apiClient.GetDomainsWithResponse(ctx, getDomainsParams)
	if err != nil {
		return nil, err
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfDomain](domainsRes)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}
	allDomains := *page.Elements
	for _, domain := range allDomains {
		for _, clusterRef := range *domain.Clusters {
			if clusterRef.Id == clusterId {
				_ = data.Set("domain_id", domain.Id)
				_ = data.Set("domain_name", domain.Name)
			}
		}
	}

	return []*schema.ResourceData{data}, nil
}

// getFlattenedHostSpecsForRefs The HostRef is supposed to have all the relevant information,
// but the backend returns everything as nil except the host ID which forces us to make a separate request
// to get some useful info about the hosts in the cluster.
func getFlattenedHostSpecsForRefs(ctx context.Context, hostRefs []vcf.HostReference,
	apiClient *vcf.ClientWithResponses) ([]map[string]interface{}, error) {
	flattenedHostSpecs := *new([]map[string]interface{})
	// Sort for reproducibility
	sort.SliceStable(hostRefs, func(i, j int) bool {
		return *hostRefs[i].Id < *hostRefs[j].Id
	})
	for _, hostRef := range hostRefs {
		res, err := apiClient.GetHostWithResponse(ctx, *hostRef.Id)
		if err != nil {
			return nil, err
		}
		hostObj, vcfErr := api_client.GetResponseAs[vcf.Host](res)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
			return nil, errors.New(*vcfErr.Message)
		}
		flattenedHostSpecs = append(flattenedHostSpecs, *FlattenHost(*hostObj))
	}
	return flattenedHostSpecs, nil
}

func getFlattenedVdsSpecsForRefs(vdsSpecs []vcf.VdsSpec) []map[string]interface{} {
	flattenedVdsSpecs := *new([]map[string]interface{})
	// Since backend API returns objects in random order sort VDSSpec list to ensure
	// import is reproducible
	sort.SliceStable(vdsSpecs, func(i, j int) bool {
		return vdsSpecs[i].Name < vdsSpecs[j].Name
	})
	for _, vdsSpec := range vdsSpecs {
		flattenedVdsSpecs = append(flattenedVdsSpecs, network.FlattenVdsSpec(vdsSpec))
	}
	return flattenedVdsSpecs
}
