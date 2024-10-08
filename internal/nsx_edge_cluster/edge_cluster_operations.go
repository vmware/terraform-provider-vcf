// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsx_edge_cluster

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
)

const (
	clusterTypeNsxT = "NSX-T"
)

func GetNsxEdgeClusterCreationSpec(data *schema.ResourceData, client *client.VcfClient) (*models.EdgeClusterCreationSpec, error) {
	// No other types are supported yet
	clusterType := clusterTypeNsxT
	adminPassword := data.Get("admin_password").(string)
	auditPassword := data.Get("audit_password").(string)
	rootPassword := data.Get("root_password").(string)
	name := data.Get("name").(string)
	profileType := data.Get("profile_type").(string)
	profileSpec := getClusterProfileSpec(data)
	routingType := data.Get("routing_type").(string)
	formFactor := data.Get("form_factor").(string)
	highAvailability := data.Get("high_availability").(string)
	tier0Name := data.Get("tier0_name").(string)
	tier1Name := data.Get("tier1_name").(string)
	mtu := int32(data.Get("mtu").(int))
	asn := int64(data.Get("asn").(int))
	tier1Unhosted := data.Get("tier1_unhosted").(bool)
	skipTepRoutabilityCheck := data.Get("skip_tep_routability_check").(bool)

	transitSubnets := resource_utils.ToStringSlice(data.Get("transit_subnets").([]interface{}))
	internalTransitSubnets := resource_utils.ToStringSlice(data.Get("internal_transit_subnets").([]interface{}))

	nodes := data.Get("edge_node").([]interface{})
	nodeSpecs := make([]*models.NsxTEdgeNodeSpec, 0, len(nodes))

	for _, node := range nodes {
		node := node.(map[string]interface{})
		nodeSpec, err := getNodeSpec(node, client)
		if err != nil {
			return nil, err
		}
		nodeSpecs = append(nodeSpecs, nodeSpec)
	}

	spec := &models.EdgeClusterCreationSpec{
		EdgeAdminPassword:             &adminPassword,
		EdgeAuditPassword:             &auditPassword,
		EdgeClusterName:               &name,
		EdgeClusterProfileSpec:        profileSpec,
		EdgeClusterProfileType:        &profileType,
		EdgeClusterType:               &clusterType,
		EdgeFormFactor:                &formFactor,
		EdgeNodeSpecs:                 nodeSpecs,
		EdgeRootPassword:              &rootPassword,
		InternalTransitSubnets:        internalTransitSubnets,
		Mtu:                           &mtu,
		Asn:                           asn,
		Tier0Name:                     tier0Name,
		Tier0RoutingType:              routingType,
		Tier0ServicesHighAvailability: highAvailability,
		Tier1Name:                     tier1Name,
		Tier1Unhosted:                 tier1Unhosted,
		TransitSubnets:                transitSubnets,
		SkipTepRoutabilityCheck:       skipTepRoutabilityCheck,
	}

	return spec, nil
}

func GetNsxEdgeClusterShrinkageSpec(currentNodes []*models.EdgeNodeReference, newNodes []interface{}) *models.EdgeClusterShrinkageSpec {
	ids := make([]string, 0)

	for _, currentNode := range currentNodes {
		found := false
		for _, newNode := range newNodes {
			name := newNode.(map[string]interface{})["name"].(string)
			if name == *currentNode.HostName {
				found = true
			}
		}
		if !found {
			ids = append(ids, *currentNode.ID)
		}
	}

	return &models.EdgeClusterShrinkageSpec{
		EdgeNodeIds: ids,
	}
}

func GetNsxEdgeClusterExpansionSpec(currentNodes []*models.EdgeNodeReference,
	newNodesRaw []interface{}, client *client.VcfClient) (*models.EdgeClusterExpansionSpec, error) {
	newNodes := getNewNodes(currentNodes, newNodesRaw)
	nodeSpecs := make([]*models.NsxTEdgeNodeSpec, 0, len(newNodes))
	spec := &models.EdgeClusterExpansionSpec{}

	for _, newNode := range newNodes {
		node := newNode.(map[string]interface{})

		adminPassword := node["admin_password"].(string)
		auditPassword := node["audit_password"].(string)
		rootPassword := node["root_password"].(string)

		spec.EdgeNodeAdminPassword = &adminPassword
		spec.EdgeNodeAuditPassword = &auditPassword
		spec.EdgeNodeRootPassword = &rootPassword

		nodeSpec, err := getNodeSpec(node, client)
		if err != nil {
			return nil, err
		}
		nodeSpecs = append(nodeSpecs, nodeSpec)
	}

	spec.EdgeNodeSpecs = nodeSpecs
	return spec, nil
}

func getNewNodes(currentNodes []*models.EdgeNodeReference, newNodesRaw []interface{}) []interface{} {
	result := make([]interface{}, 0)

	for _, newNode := range newNodesRaw {
		found := false
		name := newNode.(map[string]interface{})["name"].(string)

		for _, m := range currentNodes {
			if name == *m.HostName {
				found = true
			}
		}

		if !found {
			result = append(result, newNode)
		}
	}

	return result
}

func getNodeSpec(node map[string]interface{}, client *client.VcfClient) (*models.NsxTEdgeNodeSpec, error) {
	name := node["name"].(string)
	tep1IP := node["tep1_ip"].(string)
	tep2IP := node["tep2_ip"].(string)
	tepGateway := node["tep_gateway"].(string)
	tepVlan := int32(node["tep_vlan"].(int))

	managementIP := node["management_ip"].(string)
	managementGateway := node["management_gateway"].(string)

	firstVdsUplink := node["first_nsx_vds_uplink"].(string)
	secondVdsUplink := node["second_nsx_vds_uplink"].(string)

	interRackCluster := node["inter_rack_cluster"].(bool)

	var clusterId string
	if computeClusterId, contains := node["compute_cluster_id"]; contains {
		clusterId = computeClusterId.(string)
	}

	if computeClusterName, contains := node["compute_cluster_name"]; contains {
		if clusterId != "" {
			return nil, errors.New("you cannot set compute_cluster_id and compute_cluster_name at the same time")
		}
		cluster, err := getComputeCluster(computeClusterName.(string), client)

		if err != nil {
			return nil, err
		}

		clusterId = cluster.ID
	}

	nodeSpec := &models.NsxTEdgeNodeSpec{
		ClusterID:          &clusterId,
		EdgeNodeName:       &name,
		EdgeTep1IP:         tep1IP,
		EdgeTep2IP:         tep2IP,
		EdgeTepGateway:     tepGateway,
		EdgeTepVlan:        &tepVlan,
		ManagementGateway:  &managementGateway,
		ManagementIP:       &managementIP,
		FirstNsxVdsUplink:  firstVdsUplink,
		SecondNsxVdsUplink: secondVdsUplink,
		InterRackCluster:   &interRackCluster,
		UplinkNetwork:      getUplinkNetworkSpecs(node),
	}

	mgmtNetworkRaw := node["management_network"].([]interface{})
	if len(mgmtNetworkRaw) > 0 {
		mgmtNetworkData := mgmtNetworkRaw[0].(map[string]interface{})
		name := mgmtNetworkData["portgroup_name"].(string)
		vlan := int32(mgmtNetworkData["vlan_id"].(int))

		nodeSpec.VMManagementPortgroupName = &name
		nodeSpec.VMManagementPortgroupVlan = &vlan
	}

	return nodeSpec, nil
}

func getUplinkNetworkSpecs(node map[string]interface{}) []*models.NsxTEdgeUplinkNetwork {
	uplinks := node["uplink"].([]interface{})
	specs := make([]*models.NsxTEdgeUplinkNetwork, 0, len(uplinks))

	for _, uplink := range uplinks {
		ip := uplink.(map[string]interface{})["interface_ip"].(string)
		vlan := int32(uplink.(map[string]interface{})["vlan"].(int))
		bgpPeersRaw := uplink.(map[string]interface{})["bgp_peer"].([]interface{})
		spec := &models.NsxTEdgeUplinkNetwork{
			UplinkInterfaceIP: &ip,
			UplinkVlan:        &vlan,
			BgpPeers:          getBgpPeerSpecs(bgpPeersRaw),
		}

		specs = append(specs, spec)
	}
	return specs
}

func getBgpPeerSpecs(bgpPeersRaw []interface{}) []*models.BgpPeerSpec {
	peers := make([]*models.BgpPeerSpec, 0, len(bgpPeersRaw))

	for _, peer := range bgpPeersRaw {
		ip := peer.(map[string]interface{})["ip"].(string)
		password := peer.(map[string]interface{})["password"].(string)
		asn := int64(peer.(map[string]interface{})["asn"].(int))
		peer := &models.BgpPeerSpec{
			Asn:      &asn,
			IP:       &ip,
			Password: &password,
		}

		peers = append(peers, peer)
	}

	return peers
}

func getClusterProfileSpec(data *schema.ResourceData) *models.NsxTEdgeClusterProfileSpec {
	var profileSpec *models.NsxTEdgeClusterProfileSpec = nil
	profileType := data.Get("profile_type").(string)
	if profileType == "CUSTOM" {
		profileRaw := data.Get("profile").([]interface{})

		if len(profileRaw) > 0 {
			// there can be only one profile spec
			profile := profileRaw[0].(map[string]interface{})
			name := profile["name"].(string)
			allowedHop := int64(profile["bfd_allowed_hop"].(int))
			declareDeadMultiple := int64(profile["bfd_declare_dead_multiple"].(int))
			probeInterval := int64(profile["bfd_probe_interval"].(int))
			standbyRelocationThreshold := int64(profile["standby_relocation_threshold"].(int))

			profileSpec = &models.NsxTEdgeClusterProfileSpec{
				BfdAllowedHop:              &allowedHop,
				BfdDeclareDeadMultiple:     &declareDeadMultiple,
				BfdProbeInterval:           &probeInterval,
				EdgeClusterProfileName:     &name,
				StandbyRelocationThreshold: &standbyRelocationThreshold,
			}
		}
	}

	return profileSpec
}

func getComputeCluster(name string, client *client.VcfClient) (*models.Cluster, error) {
	params := clusters.NewGetClustersParams().WithTimeout(constants.DefaultVcfApiCallTimeout)

	ok, err := client.Clusters.GetClusters(params)

	if err != nil {
		return nil, err
	}

	computeClusters := ok.Payload.Elements

	if len(computeClusters) > 0 {
		for _, cluster := range computeClusters {
			if cluster.Name == name {
				return cluster, nil
			}
		}
	}

	return nil, fmt.Errorf("cluster %s not found", name)
}
