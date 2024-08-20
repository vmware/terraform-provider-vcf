package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/credentials"
)

func TestAccResourceAutorotatePolicy_resourceId(t *testing.T) {
	rotateDays := acctest.RandIntRange(credentials.AutoRotateDaysMin, credentials.AutorotateDaysMax)
	timeAfter := time.Now().AddDate(0, 0, rotateDays)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{{
			Config: testAccAutorotatePolicyResourceIdConfig(rotateDays),
			Check: resource.TestCheckResourceAttrWith("vcf_credentials_auto_rotate_policy.vc_0_autorotate", "auto_rotate_next_schedule", func(value string) error {
				nextAutorotate, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return err
				}
				nextAutorotateDate := nextAutorotate.Format(time.DateOnly)
				timeAfterDate := timeAfter.Format(time.DateOnly)
				if nextAutorotateDate != timeAfterDate {
					return fmt.Errorf("%s different from %s", nextAutorotate, timeAfterDate)
				}

				return nil
			}),
		}},
	})
}

func testAccAutorotatePolicyResourceIdConfig(rotateDays int) string {
	return fmt.Sprintf(`
		data "vcf_credentials" "sddc_creds" {
			resource_type = "VCENTER"
		}

		resource "vcf_credentials_auto_rotate_policy" "vc_0_autorotate" {
			resource_id = data.vcf_credentials.sddc_creds.credentials[0].resource[0].id
			resource_type = data.vcf_credentials.sddc_creds.credentials[0].resource[0].type
			user_name = data.vcf_credentials.sddc_creds.credentials[0].user_name
			enable_auto_rotation = true
			auto_rotate_days = %v
		}
	`, rotateDays)
}
