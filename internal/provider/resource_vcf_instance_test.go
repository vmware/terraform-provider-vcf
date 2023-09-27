/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

func TestAccResourceVcfSddcBasic(t *testing.T) {
	sddcName := "terraform_test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVcfSddcConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSddcResourceExists(sddcName),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "instance_id"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "creation_timestamp"),
					resource.TestCheckResourceAttrSet("vcf_instance.sddc_1", "status"),
				),
			},
		},
	})
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
						"password": "TestTest123!",
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
				"root_nsx_manager_password": "TestTest123!TestTest123!",
				"nsx_admin_password":        "TestTest123!TestTest123!",
				"nsx_audit_password":        "TestTest123!TestTest123!",
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
				"cluster_evc_mode":          "",
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
						"cidr":       "",
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
	assert.Equal(t, *sddcSpec.SDDCID, "sddcId-1001")
	assert.Equal(t, sddcSpec.DvSwitchVersion, "7.0.0")
	assert.Equal(t, sddcSpec.SkipEsxThumbprintValidation, true)
	assert.Equal(t, sddcSpec.CEIPEnabled, false)
	assert.Equal(t, *sddcSpec.TaskName, "NewStarWarsСЪКС")
	assert.Equal(t, *sddcSpec.SDDCManagerSpec.IPAddress, "10.0.0.4")
	assert.Equal(t, *sddcSpec.SDDCManagerSpec.Hostname, "sddc-manager")
	// TODO test all fields in the SDDC Spec
}

func testAccCheckSddcResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no sddc id is set")
		}
		sddcBringUpID := rs.Primary.Attributes["instance_id"]
		client := testAccProvider.Meta().(*api_client.CloudBuilderClient)
		response, err := getLastBringUp(context.Background(), client)
		if err != nil {
			return fmt.Errorf("error occurred while retrieving all sddcs")
		}
		foundSddc := response
		if foundSddc.ID != sddcBringUpID {
			return fmt.Errorf("error retrieving SDDC Bring-Up with id %s", sddcBringUpID)
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
		  password = "TestTest123!"
		}
		ip_address = "10.0.0.4"
		hostname = "sddc-manager"
		root_user_credentials {
		  username = "root"
		  password = "TestTest123!"
		}
	  }
	  ntp_servers = [
		"10.0.0.250"
	  ]
	  dns {
		domain = "vsphere.local"
		name_server = "10.0.0.250"
		secondary_name_server = "10.0.0.250"
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
		root_nsx_manager_password = "TestTest123!TestTest123!"
		nsx_admin_password = "TestTest123!TestTest123!"
		nsx_audit_password = "TestTest123!TestTest123!"
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
		admin_user_sso_password = "TestTest123!"
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
