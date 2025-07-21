// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
)

func TestAccResourceVcfInstanceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVcfSddcConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceResourceExists(),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "instance_id"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "creation_timestamp"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "status"),
				),
			},
		},
	})
}

func testAccCheckInstanceResourceExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "vcf_instance" {
				continue
			}

			instanceId := rs.Primary.Attributes["id"]
			client := testAccProvider.Meta().(*api_client.InstallerClient)
			response, err := getLastBringUp(context.Background(), client)
			if err != nil {
				return fmt.Errorf("error occurred while retrieving all sddcs")
			}
			foundInstance := response
			if *foundInstance.Id != instanceId {
				return fmt.Errorf("error retrieving SDDC Bring-Up with id %s", instanceId)
			}
			return nil

		}
		return nil
	}
}

func testAccCheckVcfSddcConfigBasic() string {
	return fmt.Sprintf(
		`resource "vcf_instance" "sddc_1" {
	  instance_id = "sddcId-1001"
	  skip_esx_thumbprint_validation = true
	  management_pool_name = "bringup-networkpool"
	  ceip_enabled = false
	  version = "5.2.0"
	  sddc_manager {
		hostname = "sddc-manager"
		ssh_password = "MnogoSl0jn@P@rol@!"
		root_user_password = "MnogoSl0jn@P@rol@!"
		local_user_password = "MnogoSl0jn@P@rol@!"
	  }
	  ntp_servers = [
		"10.0.0.250"
	  ]
	  dns {
		domain = "vrack.vsphere.local"
		name_server = "10.0.0.250"
	  }
	  network {
		subnet = "10.0.0.0/22"
		vlan_id = "0"
		mtu = "1500"
		network_type = "MANAGEMENT"
		gateway = "10.0.0.250"
		active_uplinks = [
			"uplink1",
			"uplink2"
		]
	  }
	  network {
		active_uplinks = [
			"uplink1",
			"uplink2"
		]
		subnet = "10.0.4.0/24"
		include_ip_address_ranges {
		  start_ip_address = "10.0.4.7"
		  end_ip_address = "10.0.4.48"
		}
		include_ip_address_ranges {
		  start_ip_address = "10.0.4.3"
		  end_ip_address = "10.0.4.6"
		}
		include_ip_address = [
		  "10.0.4.50",
		  "10.0.4.49"]
		vlan_id = "0"
		mtu = "8940"
		network_type = "VSAN"
		gateway = "10.0.4.253"
	  }
	  network {
		active_uplinks = [
			"uplink1",
			"uplink2"
		]
		subnet = "10.0.8.0/24"
		include_ip_address_ranges {
		  start_ip_address = "10.0.8.3"
		  end_ip_address = "10.0.8.50"
		}
		vlan_id = "0"
		mtu = "8940"
		network_type = "VMOTION"
		gateway = "10.0.8.253"
	  }
	  nsx {
		nsx_manager_size = "medium"
		nsx_manager {
		  hostname = "nsx-mgmt-1"
		}
		root_nsx_manager_password = "MnogoSl0jn@P@rol@!"
		nsx_admin_password = "MnogoSl0jn@P@rol@!"
		nsx_audit_password = "MnogoSl0jn@P@rol@!"
		vip_fqdn = "vip-nsx-mgmt"
		transport_vlan_id = 0
	  }
	  vsan {
		datastore_name = "sfo01-m01-vsan"
		failures_to_tolerate = 1
	  }
	  dvs {
		mtu = 8940
		nioc {
		  traffic_type = "VSAN"
		  value = "HIGH"
		}
		nioc {
		  traffic_type = "VMOTION"
		  value = "LOW"
		}
		nioc {
		  traffic_type = "VDP"
		  value = "LOW"
		}
		nioc {
		  traffic_type = "VIRTUALMACHINE"
		  value = "HIGH"
		}
		nioc {
		  traffic_type = "MANAGEMENT"
		  value = "NORMAL"
		}
		nioc {
		  traffic_type = "NFS"
		  value = "LOW"
		}
		nioc {
		  traffic_type = "HBR"
		  value = "LOW"
		}
		nioc {
		  traffic_type = "FAULTTOLERANCE"
		  value = "LOW"
		}
		nioc {
		  traffic_type = "ISCSI"
		  value = "LOW"
		}
		dvs_name = "SDDC-Dswitch-Private"
		vmnic_mapping {
			vmnic = "vmnic0"
			uplink = "uplink1"
		}
		vmnic_mapping {
			vmnic = "vmnic1"
			uplink = "uplink2"
		}
		networks = [
		  "MANAGEMENT",
		  "VSAN",
		  "VMOTION"
		]
		nsx_teamings {
		  policy = "LOADBALANCE_SRCID"
		  active_uplinks = ["uplink1", "uplink2"]
		}
		nsxt_switch_config {
		  host_switch_operational_mode = "ENS_INTERRUPT"
		  transport_zones {
			name = "nsx-vlan-transportzone"
			transport_type = "VLAN"
		  }
		  transport_zones {
			name = "overlay-tz-sfo-m01-nsx01"
			transport_type = "OVERLAY"
		  }
		}
	  }
	  cluster {
		datacenter_name = "SDDC-Datacenter"
		cluster_name = "SDDC-Cluster1"
		cluster_evc_mode = ""
		resource_pool {
		  name = "Mgmt-ResourcePool"
		  type = "management"
		}
		resource_pool {
		  name = "Network-ResourcePool"
		  type = "network"
		}
		resource_pool {
		  name = "Compute-ResourcePool"
		  type = "compute"
		}
		resource_pool {
		  name = "User-RP"
		  type = "compute"
		}
	  }
	  vcenter {
		vcenter_hostname = "vcenter-1"
		root_vcenter_password = "MnogoSl0jn@P@rol@!"
		vm_size = "tiny"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		hostname = "esxi-1"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		hostname = "esxi-2"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		hostname = "esxi-3"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		hostname = "esxi-4"
	  }
	}`,
		os.Getenv(constants.VcfTestHost1Pass),
		os.Getenv(constants.VcfTestHost2Pass),
		os.Getenv(constants.VcfTestHost3Pass),
		os.Getenv(constants.VcfTestHost4Pass))
}

func TestVcfInstanceSchemaParse(t *testing.T) {
	input := map[string]interface{}{
		"instance_id":                    "sddcId-1001",
		"skip_esx_thumbprint_validation": true,
		"ceip_enabled":                   false,
		"version":                        "5.2.0",
		"sddc_manager": []interface{}{
			map[string]interface{}{
				"hostname":            "sddc-manager",
				"root_user_password":  "MnogoSl0jn@P@rol@!",
				"local_user_password": "MnogoSl0jn@P@rol@!",
				"ssh_password":        "MnogoSl0jn@P@rol@!",
			},
		},
		"ntp_servers": []interface{}{"10.0.0.250"},
		"dns": []interface{}{
			map[string]interface{}{
				"domain":                "vsphere.local",
				"name_server":           "10.0.0.250",
				"secondary_name_server": "10.0.0.250",
			},
		},
		"network": []interface{}{
			map[string]interface{}{
				"vlan_id":      0,
				"mtu":          "8940",
				"network_type": "VSAN",
				"gateway":      "10.0.4.253",
				"include_ip_address_ranges": []interface{}{
					map[string]interface{}{
						"start_ip_address": "10.0.4.7",
						"end_ip_address":   "10.0.4.48",
					},
					map[string]interface{}{
						"start_ip_address": "10.0.4.3",
						"end_ip_address":   "10.0.4.6",
					},
				},
				"include_ip_address": []interface{}{
					"10.0.4.50",
					"10.0.4.49",
				},
			},
		},
		"nsx": []interface{}{
			map[string]interface{}{
				"nsx_manager_size": "medium",
				"nsx_manager": []interface{}{
					map[string]interface{}{
						"hostname": "nsx-mgmt-1",
					},
				},
				"root_nsx_manager_password": "MnogoSl0jn@P@rol@!",
				"nsx_admin_password":        "MnogoSl0jn@P@rol@!",
				"nsx_audit_password":        "MnogoSl0jn@P@rol@!",
				"vip_fqdn":                  "vip-nsx-mgmt",
				"transport_vlan_id":         0,
				"overlay_transport_zone": []interface{}{
					map[string]interface{}{
						"zone_name":    "overlay-tz",
						"network_name": "net-overlay",
					},
				},
			},
		},
		"vsan": []interface{}{
			map[string]interface{}{
				"datastore_name":       "sfo01-m01-vsan",
				"failures_to_tolerate": 1,
			},
		},
		"dvs": []interface{}{
			map[string]interface{}{
				"mtu":      "8940",
				"dvs_name": "SDDC-Dswitch-Private",
				"nioc": []interface{}{
					map[string]interface{}{
						"traffic_type": "VDP",
						"value":        "LOW",
					},
					map[string]interface{}{
						"traffic_type": "VMOTION",
						"value":        "LOW",
					},
					map[string]interface{}{
						"traffic_type": "VSAN",
						"value":        "HIGH",
					},
				},
				"vmnic_mapping": []interface{}{
					map[string]interface{}{
						"vmnic":  "vmnic0",
						"uplink": "uplink1",
					},
					map[string]interface{}{
						"vmnic":  "vmnic1",
						"uplink": "uplink2",
					},
				},
				"networks": []interface{}{
					"MANAGEMENT",
					"VSAN",
					"VMOTION",
				},
				"nsx_teaming": []interface{}{
					map[string]interface{}{
						"policy":         "LOADBALANCE_SRCID",
						"active_uplinks": []interface{}{"uplink1", "uplink2"},
					},
				},
				"nsxt_switch_config": []interface{}{
					map[string]interface{}{
						"host_switch_operational_mode": "ENS_INTERRUPT",
						"transport_zones": []interface{}{
							map[string]interface{}{
								"name":           "nsx-vlan-transportzone",
								"transport_type": "VLAN",
							},
							map[string]interface{}{
								"name":           "overlay-tz-sfo-m01-nsx01",
								"transport_type": "OVERLAY",
							},
						},
					},
				},
			},
		},
		"cluster": []interface{}{
			map[string]interface{}{
				"cluster_name":              "SDDC-Cluster1",
				"host_failures_to_tolerate": 2,
				"resource_pool": []interface{}{
					map[string]interface{}{
						"name": "Mgmt-ResourcePool",
						"type": "management",
					},
					map[string]interface{}{
						"name":                          "Compute-ResourcePool",
						"type":                          "compute",
						"cpu_reservation_expandable":    false,
						"cpu_reservation_mhz":           1000,
						"cpu_reservation_percentage":    10,
						"cpu_shares_level":              "normal",
						"cpu_shares_value":              10,
						"memory_reservation_expandable": false,
						"memory_reservation_mb":         "1000",
						"memory_shares_level":           "normal",
						"memory_shares_value":           10,
					},
				},
			},
		},
		"vcenter": []interface{}{
			map[string]interface{}{
				"vcenter_hostname":      "vcenter-1",
				"root_vcenter_password": "TestTest1!",
				"vm_size":               "tiny",
			},
		},
		"host": []interface{}{
			map[string]interface{}{
				"credentials": []interface{}{
					map[string]interface{}{
						"username": "root",
						"password": "MnogoSl0jn@P@rol@!",
					},
				},
				"ip_address_private": []interface{}{
					map[string]interface{}{
						"subnet":     "255.255.252.0",
						"ip_address": "10.0.0.100",
						"gateway":    "10.0.0.250",
					},
				},
				"hostname":    "esxi-1",
				"vswitch":     "vSwitch0",
				"association": "SDDC-Datacenter",
			},
		},
		"automation": []interface{}{
			map[string]interface{}{
				"hostname":              "automation-1",
				"internal_cluster_cidr": "240.0.0.0/15",
				"ip_pool":               []interface{}{"10.0.0.81", "10.0.0.91"},
				"admin_user_password":   "MnogoSl0jn@P@rol@!",
				"node_prefix":           "automation-node",
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"admin_user_password": "MnogoSl0jn@P@rol@!",
				"appliance_size":      "medium",
				"load_balancer_fqdn":  "load-balancer-fqdn",
				"node": []interface{}{
					map[string]interface{}{
						"hostname":           "operations-1",
						"type":               "master",
						"root_user_password": "MnogoSl0jn@P@rol@!",
					},
				},
			},
		},
		"operations_fleet_management": []interface{}{
			map[string]interface{}{
				"hostname":            "operations-1",
				"admin_user_password": "MnogoSl0jn@P@rol@!",
				"root_user_password":  "MnogoSl0jn@P@rol@!",
			},
		},
		"operations_collector": []interface{}{
			map[string]interface{}{
				"hostname":           "operations-1",
				"appliance_size":     "medium",
				"root_user_password": "MnogoSl0jn@P@rol@!",
			},
		},
	}
	var testResourceData = schema.TestResourceDataRaw(t, resourceVcfInstanceSchema(), input)
	sddcSpec := buildSddcSpec(testResourceData)

	// assert.Equal determines pointer equality based on the referenced values and not by the actual memory addresses
	assert.Equal(t, "sddcId-1001", sddcSpec.SddcId)
	assert.Equal(t, utils.ToPointer[bool](true), sddcSpec.SkipEsxThumbprintValidation)
	assert.Equal(t, utils.ToPointer[bool](nil), sddcSpec.CeipEnabled)
	assert.Equal(t, "sddc-manager", sddcSpec.SddcManagerSpec.Hostname)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.SddcManagerSpec.RootPassword)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.SddcManagerSpec.LocalUserPassword)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.SddcManagerSpec.SshPassword)
	assert.Equal(t, utils.ToPointer[[]string]([]string{"10.0.0.250"}), sddcSpec.NtpServers)
	assert.Equal(t, "vsphere.local", sddcSpec.DnsSpec.Subdomain)
	assert.Equal(t, utils.ToPointer[[]string]([]string{"10.0.0.250", "10.0.0.250"}), sddcSpec.DnsSpec.Nameservers)
	assert.Equal(t, int32(0), sddcSpec.NetworkSpecs[0].VlanId)
	assert.Equal(t, utils.ToPointer[int32](int32(8940)), sddcSpec.NetworkSpecs[0].Mtu)
	assert.Equal(t, "VSAN", sddcSpec.NetworkSpecs[0].NetworkType)
	assert.Equal(t, utils.ToPointer[string]("10.0.4.253"), sddcSpec.NetworkSpecs[0].Gateway)
	assert.Equal(t, utils.ToPointer[[]string]([]string{"10.0.4.50", "10.0.4.49"}), sddcSpec.NetworkSpecs[0].IncludeIpAddress)
	assert.Equal(t, "10.0.4.7", (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[0].StartIpAddress)
	assert.Equal(t, "10.0.4.48", (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[0].EndIpAddress)
	assert.Equal(t, "10.0.4.3", (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[1].StartIpAddress)
	assert.Equal(t, "10.0.4.6", (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[1].EndIpAddress)
	assert.Equal(t, "medium", *sddcSpec.NsxtSpec.NsxtManagerSize)
	assert.Equal(t, utils.ToPointer[string]("nsx-mgmt-1"), sddcSpec.NsxtSpec.NsxtManagers[0].Hostname)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.NsxtSpec.RootNsxtManagerPassword)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.NsxtSpec.NsxtAdminPassword)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.NsxtSpec.NsxtAuditPassword)
	assert.Equal(t, "vip-nsx-mgmt", sddcSpec.NsxtSpec.VipFqdn)
	assert.Equal(t, utils.ToPointer[int32](int32(0)), sddcSpec.NsxtSpec.TransportVlanId)
	assert.Equal(t, utils.ToPointer[string]("sfo01-m01-vsan"), sddcSpec.DatastoreSpec.VsanSpec.DatastoreName)
	assert.Equal(t, utils.ToPointer[int32](int32(8940)), (*sddcSpec.DvsSpecs)[0].Mtu)
	assert.Equal(t, utils.ToPointer[string]("SDDC-Dswitch-Private"), (*sddcSpec.DvsSpecs)[0].DvsName)
	assert.Equal(t, utils.ToPointer[[]string]([]string{"MANAGEMENT", "VSAN", "VMOTION"}), (*sddcSpec.DvsSpecs)[0].Networks)
	assert.Equal(t, utils.ToPointer[string]("SDDC-Cluster1"), sddcSpec.ClusterSpec.ClusterName)
	assert.Equal(t, utils.ToPointer[string](""), sddcSpec.ClusterSpec.ClusterEvcMode)
	assert.Equal(t, utils.ToPointer[string]("Mgmt-ResourcePool"), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[0].Name)
	assert.Equal(t, utils.ToPointer[installer.ResourcePoolSpecType](installer.ResourcePoolSpecType("management")), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[0].Type)
	assert.Equal(t, "Compute-ResourcePool", *(*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].Name)
	assert.Equal(t, utils.ToPointer[installer.ResourcePoolSpecType](installer.ResourcePoolSpecType("compute")), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].Type)
	assert.Equal(t, utils.ToPointer[bool](false), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationExpandable)
	assert.Equal(t, utils.ToPointer[int64](int64(1000)), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationMhz)
	assert.Equal(t, utils.ToPointer[int32](int32(10)), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationPercentage)
	assert.Equal(t, utils.ToPointer[installer.ResourcePoolSpecCpuSharesLevel](installer.ResourcePoolSpecCpuSharesLevel("normal")), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuSharesLevel)
	assert.Equal(t, utils.ToPointer[int32](int32(10)), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuSharesValue)
	assert.Equal(t, false, *(*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemoryReservationExpandable)
	assert.Equal(t, utils.ToPointer[int64](int64(1000)), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemoryReservationMb)
	assert.Equal(t, utils.ToPointer[installer.ResourcePoolSpecMemorySharesLevel](installer.ResourcePoolSpecMemorySharesLevel("normal")), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemorySharesLevel)
	assert.Equal(t, utils.ToPointer[int32](int32(10)), (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemorySharesValue)
	assert.Equal(t, "vcenter-1", sddcSpec.VcenterSpec.VcenterHostname)
	assert.Equal(t, "TestTest1!", sddcSpec.VcenterSpec.RootVcenterPassword)
	assert.Equal(t, utils.ToPointer[string]("tiny"), sddcSpec.VcenterSpec.VmSize)
	assert.Equal(t, "esxi-1", (*sddcSpec.HostSpecs)[0].Hostname)
	assert.Equal(t, "MnogoSl0jn@P@rol@!", (*sddcSpec.HostSpecs)[0].Credentials.Password)
	assert.Equal(t, utils.ToPointer[string]("root"), (*sddcSpec.HostSpecs)[0].Credentials.Username)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), (*sddcSpec.VcfAutomationSpec).AdminUserPassword)
	assert.Equal(t, "automation-1", (*sddcSpec.VcfAutomationSpec).Hostname)
	assert.Equal(t, utils.ToPointer[string]("automation-node"), (*sddcSpec.VcfAutomationSpec).NodePrefix)
	assert.Equal(t, "10.0.0.81", (*sddcSpec.VcfAutomationSpec.IpPool)[0])
	assert.Equal(t, "10.0.0.91", (*sddcSpec.VcfAutomationSpec.IpPool)[1])
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), (*sddcSpec.VcfOperationsSpec).AdminUserPassword)
	assert.Equal(t, utils.ToPointer[string]("medium"), (*sddcSpec.VcfOperationsSpec).ApplianceSize)
	assert.Equal(t, utils.ToPointer[string]("load-balancer-fqdn"), (*sddcSpec.VcfOperationsSpec).LoadBalancerFqdn)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), sddcSpec.VcfOperationsSpec.Nodes[0].RootUserPassword)
	assert.Equal(t, utils.ToPointer[string]("master"), sddcSpec.VcfOperationsSpec.Nodes[0].Type)
	assert.Equal(t, "operations-1", sddcSpec.VcfOperationsSpec.Nodes[0].Hostname)
	assert.Equal(t, "operations-1", sddcSpec.VcfOperationsCollectorSpec.Hostname)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), (*sddcSpec.VcfOperationsCollectorSpec).RootUserPassword)
	assert.Equal(t, utils.ToPointer[string]("medium"), (*sddcSpec.VcfOperationsCollectorSpec).ApplianceSize)
	assert.Equal(t, "operations-1", sddcSpec.VcfOperationsFleetManagementSpec.Hostname)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), (*sddcSpec.VcfOperationsFleetManagementSpec).RootUserPassword)
	assert.Equal(t, utils.ToPointer[string]("MnogoSl0jn@P@rol@!"), (*sddcSpec.VcfOperationsFleetManagementSpec).AdminUserPassword)
	assert.Equal(t, utils.ToStringPointer("5.2.0"), sddcSpec.Version)
	assert.Equal(t, utils.ToPointer[int32](int32(1)), sddcSpec.DatastoreSpec.VsanSpec.FailuresToTolerate)
	assert.Equal(t, "LOADBALANCE_SRCID", (*(*sddcSpec.DvsSpecs)[0].NsxTeamings)[0].Policy)
	assert.Equal(t, utils.ToStringPointer("ENS_INTERRUPT"), (*sddcSpec.DvsSpecs)[0].NsxtSwitchConfig.HostSwitchOperationalMode)
}
