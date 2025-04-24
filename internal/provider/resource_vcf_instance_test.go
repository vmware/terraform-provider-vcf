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
)

func TestAccResourceVcfSddcBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVcfSddcConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSddcResourceExists(),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "instance_id"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "creation_timestamp"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "status"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "sddc_manager_fqdn"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "sddc_manager_id"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "sddc_manager_version"),
				),
			},
		},
	})
}

func testAccCheckSddcResourceExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "vcf_instance" {
				continue
			}

			instanceId := rs.Primary.Attributes["id"]
			client := testAccProvider.Meta().(*api_client.CloudBuilderClient)
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
	  dv_switch_version = "7.0.0"
	  skip_esx_thumbprint_validation = true
	  management_pool_name = "bringup-networkpool"
	  ceip_enabled = false
	  esx_license = %q
	  task_name = "workflowconfig/workflowspec-ems.json"
	  sddc_manager {
		second_user_credentials {
		  username = "vcf"
		  password = "MnogoSl0jn@P@rol@!"
		}
		ip_address = "10.0.0.4"
		hostname = "sddc-manager"
		root_user_credentials {
		  username = "root"
		  password = "MnogoSl0jn@P@rol@!"
		}
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
	  }
	  network {
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
		  ip = "10.0.0.31"
		}
		root_nsx_manager_password = "MnogoSl0jn@P@rol@!"
		nsx_admin_password = "MnogoSl0jn@P@rol@!"
		nsx_audit_password = "MnogoSl0jn@P@rol@!"
		overlay_transport_zone {
		  zone_name = "overlay-tz"
		  network_name = "net-overlay"
		}
		vip = "10.0.0.30"
		vip_fqdn = "vip-nsx-mgmt"
		license = %q
		transport_vlan_id = 0
	  }
	  vsan {
		license = %q
		datastore_name = "sfo01-m01-vsan"
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
		vmnics = [
		  "vmnic0",
		  "vmnic1"
		]
		networks = [
		  "MANAGEMENT",
		  "VSAN",
		  "VMOTION"
		]
	  }
	  cluster {
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
	  psc {
		psc_sso_domain = "vsphere.local"
		admin_user_sso_password = "MnogoSl0jn@P@rol@!"
	  }
	  vcenter {
		vcenter_ip = "10.0.0.6"
		vcenter_hostname = "vcenter-1"
		license = %q
		root_vcenter_password = "TestTest1!"
		vm_size = "tiny"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		ip_address_private {
		  subnet = "255.255.252.0"
		  cidr = ""
		  ip_address = "10.0.0.100"
		  gateway = "10.0.0.250"
		}
		hostname = "esxi-1"
		vswitch = "vSwitch0"
		association = "SDDC-Datacenter"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		ip_address_private {
		  subnet = "255.255.252.0"
		  cidr = ""
		  ip_address = "10.0.0.101"
		  gateway = "10.0.0.250"
		}
		hostname = "esxi-2"
		vswitch = "vSwitch0"
		association = "SDDC-Datacenter"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		ip_address_private {
		  subnet = "255.255.255.0"
		  cidr = ""
		  ip_address = "10.0.0.102"
		  gateway = "10.0.0.250"
		}
		hostname = "esxi-3"
		vswitch = "vSwitch0"
		association = "SDDC-Datacenter"
	  }
	  host {
		credentials {
		  username = "root"
		  password = %q
		}
		ip_address_private {
		  subnet = "255.255.255.0"
		  cidr = ""
		  ip_address = "10.0.0.103"
		  gateway = "10.0.0.250"
		}
		hostname = "esxi-4"
		vswitch = "vSwitch0"
		association = "SDDC-Datacenter"
	  }
	}`,
		os.Getenv(constants.VcfTestEsxiLicenseKey),
		os.Getenv(constants.VcfTestNsxLicenseKey),
		os.Getenv(constants.VcfTestVsanLicenseKey),
		os.Getenv(constants.VcfTestVcenterLicenseKey),
		os.Getenv(constants.VcfTestHost1Pass),
		os.Getenv(constants.VcfTestHost2Pass),
		os.Getenv(constants.VcfTestHost3Pass),
		os.Getenv(constants.VcfTestHost4Pass))
}

func TestVcfInstanceSchemaParse(t *testing.T) {
	input := map[string]interface{}{
		"instance_id":                    "sddcId-1001",
		"dv_switch_version":              "7.0.0",
		"skip_esx_thumbprint_validation": true,
		"ceip_enabled":                   false,
		"task_name":                      "NewStarWarsСЪКС",
		"sddc_manager": []interface{}{
			map[string]interface{}{
				"ip_address": "10.0.0.4",
				"hostname":   "sddc-manager",
				"root_user_credentials": []interface{}{
					map[string]interface{}{
						"username": "root",
						"password": "MnogoSl0jn@P@rol@!",
					},
				},
				"second_user_credentials": []interface{}{
					map[string]interface{}{
						"username": "vcf",
						"password": "MnogoSl0jn@P@rol@!",
					},
				},
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
						"ip":       "10.0.0.31",
					},
				},
				"root_nsx_manager_password": "MnogoSl0jn@P@rol@!",
				"nsx_admin_password":        "MnogoSl0jn@P@rol@!",
				"nsx_audit_password":        "MnogoSl0jn@P@rol@!",
				"vip":                       "10.0.0.30",
				"vip_fqdn":                  "vip-nsx-mgmt",
				"license":                   "XXX",
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
				"license":        "XXX",
				"datastore_name": "sfo01-m01-vsan",
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
				"vmnics": []interface{}{
					"vmnic0",
					"vmnic1",
				},
				"networks": []interface{}{
					"MANAGEMENT",
					"VSAN",
					"VMOTION",
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
		"psc": []interface{}{
			map[string]interface{}{
				"psc_sso_domain":          "vsphere.local",
				"admin_user_sso_password": "TestTest123!",
			},
		},
		"vcenter": []interface{}{
			map[string]interface{}{
				"vcenter_ip":            "10.0.0.6",
				"vcenter_hostname":      "vcenter-1",
				"license":               "XXX",
				"root_vcenter_password": "TestTest1!",
				"vm_size":               "tiny",
			},
		},
		"host": []interface{}{
			map[string]interface{}{
				"credentials": []interface{}{
					map[string]interface{}{
						"username": "root",
						"password": "TestTest123!",
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
	}
	var testResourceData = schema.TestResourceDataRaw(t, resourceVcfInstanceSchema(), input)
	sddcSpec := buildSddcSpec(testResourceData)
	assert.Equal(t, sddcSpec.SddcId, "sddcId-1001")
	assert.Equal(t, sddcSpec.DvSwitchVersion, "7.0.0")
	assert.Equal(t, sddcSpec.SkipEsxThumbprintValidation, true)
	assert.Equal(t, sddcSpec.CeipEnabled, false)
	assert.Equal(t, *sddcSpec.TaskName, "NewStarWarsСЪКС")
	assert.Equal(t, *sddcSpec.SddcManagerSpec.IpAddress, "10.0.0.4")
	assert.Equal(t, sddcSpec.SddcManagerSpec.Hostname, "sddc-manager")
	assert.Equal(t, *sddcSpec.SddcManagerSpec.RootUserCredentials.Username, "root")
	assert.Equal(t, sddcSpec.SddcManagerSpec.RootUserCredentials.Password, "MnogoSl0jn@P@rol@!")
	assert.Equal(t, *sddcSpec.SddcManagerSpec.SecondUserCredentials.Username, "vcf")
	assert.Equal(t, sddcSpec.SddcManagerSpec.SecondUserCredentials.Password, "MnogoSl0jn@P@rol@!")
	assert.Equal(t, sddcSpec.NtpServers, []string{"10.0.0.250"})
	assert.Equal(t, *sddcSpec.DnsSpec.Domain, "vsphere.local")
	assert.Equal(t, *sddcSpec.DnsSpec.Domain, "vsphere.local")
	assert.Equal(t, sddcSpec.DnsSpec.Nameserver, "10.0.0.250")
	assert.Equal(t, sddcSpec.DnsSpec.SecondaryNameserver, "10.0.0.250")
	assert.Equal(t, sddcSpec.NetworkSpecs[0].VlanId, "0")
	assert.Equal(t, sddcSpec.NetworkSpecs[0].Mtu, "8940")
	assert.Equal(t, sddcSpec.NetworkSpecs[0].NetworkType, "VSAN")
	assert.Equal(t, sddcSpec.NetworkSpecs[0].Gateway, "10.0.4.253")
	assert.Equal(t, sddcSpec.NetworkSpecs[0].IncludeIpAddress, []string{"10.0.4.50", "10.0.4.49"})
	assert.Equal(t, (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[0].StartIpAddress, "10.0.4.7")
	assert.Equal(t, (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[0].EndIpAddress, "10.0.4.48")
	assert.Equal(t, (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[1].StartIpAddress, "10.0.4.3")
	assert.Equal(t, (*sddcSpec.NetworkSpecs[0].IncludeIpAddressRanges)[1].EndIpAddress, "10.0.4.6")
	assert.Equal(t, *sddcSpec.NsxtSpec.NsxtManagerSize, "medium")
	assert.Equal(t, sddcSpec.NsxtSpec.NsxtManagers[0].Hostname, "nsx-mgmt-1")
	assert.Equal(t, sddcSpec.NsxtSpec.NsxtManagers[0].Ip, "10.0.0.31")
	assert.Equal(t, *sddcSpec.NsxtSpec.RootNsxtManagerPassword, "MnogoSl0jn@P@rol@!")
	assert.Equal(t, sddcSpec.NsxtSpec.NsxtAdminPassword, "MnogoSl0jn@P@rol@!")
	assert.Equal(t, sddcSpec.NsxtSpec.NsxtAuditPassword, "MnogoSl0jn@P@rol@!")
	assert.Equal(t, *sddcSpec.NsxtSpec.Vip, "10.0.0.30")
	assert.Equal(t, sddcSpec.NsxtSpec.VipFqdn, "vip-nsx-mgmt")
	assert.Equal(t, sddcSpec.NsxtSpec.NsxtLicense, "XXX")
	assert.Equal(t, sddcSpec.NsxtSpec.TransportVlanId, int32(0))
	assert.Equal(t, sddcSpec.NsxtSpec.OverLayTransportZone.ZoneName, "overlay-tz")
	assert.Equal(t, sddcSpec.NsxtSpec.OverLayTransportZone.NetworkName, "net-overlay")
	assert.Equal(t, sddcSpec.VsanSpec.LicenseFile, "XXX")
	assert.Equal(t, *sddcSpec.VsanSpec.DatastoreName, "sfo01-m01-vsan")
	assert.Equal(t, (*sddcSpec.DvsSpecs)[0].Mtu, int32(8940))
	assert.Equal(t, (*sddcSpec.DvsSpecs)[0].DvsName, "SDDC-Dswitch-Private")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[0].TrafficType, "VDP")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[0].Value, "LOW")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[1].TrafficType, "VMOTION")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[1].Value, "LOW")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[2].TrafficType, "VSAN")
	assert.Equal(t, (*(*sddcSpec.DvsSpecs)[0].NiocSpecs)[2].Value, "HIGH")
	assert.Equal(t, (*sddcSpec.DvsSpecs)[0].Vmnics, []string{"vmnic0", "vmnic1"})
	assert.Equal(t, (*sddcSpec.DvsSpecs)[0].Networks, []string{"MANAGEMENT", "VSAN", "VMOTION"})
	assert.Equal(t, *sddcSpec.ClusterSpec.ClusterName, "SDDC-Cluster1")
	assert.Equal(t, sddcSpec.ClusterSpec.ClusterEvcMode, "")
	assert.Equal(t, sddcSpec.ClusterSpec.HostFailuresToTolerate, utils.ToInt32Pointer(2))
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[0].Name, "Mgmt-ResourcePool")
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[0].Type, "management")
	assert.Equal(t, *(*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].Name, "Compute-ResourcePool")
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].Type, "compute")
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationExpandable, false)
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationMhz, int64(1000))
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuReservationPercentage, utils.ToInt32Pointer(10))
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuSharesLevel, "normal")
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].CpuSharesValue, int32(10))
	assert.Equal(t, *(*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemoryReservationExpandable, false)
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemoryReservationMb, int64(1000))
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemorySharesLevel, "normal")
	assert.Equal(t, (*sddcSpec.ClusterSpec.ResourcePoolSpecs)[1].MemorySharesValue, int32(10))
	assert.Equal(t, (*sddcSpec.PscSpecs)[0].AdminUserSsoPassword, "TestTest123!")
	assert.Equal(t, (*sddcSpec.PscSpecs)[0].PscSsoSpec.SsoDomain, "vsphere.local")
	assert.Equal(t, sddcSpec.VcenterSpec.VcenterIp, "10.0.0.6")
	assert.Equal(t, sddcSpec.VcenterSpec.VcenterHostname, "vcenter-1")
	assert.Equal(t, sddcSpec.VcenterSpec.LicenseFile, "XXX")
	assert.Equal(t, sddcSpec.VcenterSpec.RootVcenterPassword, "TestTest1!")
	assert.Equal(t, sddcSpec.VcenterSpec.VmSize, "tiny")
	assert.Equal(t, *sddcSpec.HostSpecs[0].Credentials.Username, "root")
	assert.Equal(t, sddcSpec.HostSpecs[0].Credentials.Password, "TestTest123!")
	assert.Equal(t, sddcSpec.HostSpecs[0].Hostname, "esxi-1")
	assert.Equal(t, *sddcSpec.HostSpecs[0].VSwitch, "vSwitch0")
	assert.Equal(t, *sddcSpec.HostSpecs[0].Association, "SDDC-Datacenter")
	assert.Equal(t, sddcSpec.HostSpecs[0].IpAddressPrivate.IpAddress, "10.0.0.100")
	assert.Equal(t, sddcSpec.HostSpecs[0].IpAddressPrivate.Subnet, "255.255.252.0")
	assert.Equal(t, sddcSpec.HostSpecs[0].IpAddressPrivate.Cidr, "")
	assert.Equal(t, sddcSpec.HostSpecs[0].IpAddressPrivate.Gateway, "10.0.0.250")
}
