/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcf

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/client/ceip"
	"github.com/vmware/vcf-sdk-go/models"

	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCeip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCeipCreate,
		ReadContext:   resourceCeipRead,
		UpdateContext: resourceCeipUpdate,
		DeleteContext: resourceCeipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User provided CEIP operation. One among: ENABLE, DISABLE",
			},
		},
	}
}

func resourceCeipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCeipUpdate(ctx, d, meta)
}

func resourceCeipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	ok, err := apiClient.CEIP.GETCEIPStatus(nil)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	jsonp, _ := json.MarshalIndent(ok.Payload, " ", " ")
	log.Println(string(jsonp))

	d.SetId(ok.Payload.InstanceID)
	return nil
}

func resourceCeipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	params := ceip.NewUpdateCEIPStatusParams()
	updateSpec := models.CEIPUpdateSpec{}

	if status, ok := d.GetOk("status"); ok {
		statusVal := status.(string)
		updateSpec.Status = &statusVal
	}

	params.CEIPUpdateSpec = &updateSpec
	_, ceipAccepted, err := apiClient.CEIP.UpdateCEIPStatus(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	if vcfClient.WaitForTask(ceipAccepted.Payload.ID) != nil {
		return diag.FromErr(err)
	}

	return resourceCeipRead(ctx, d, meta)
}

/**
 * Mapping deletion of ceip resource to disabling ceip.
 */
func resourceCeipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	params := ceip.NewUpdateCEIPStatusParams()
	updateSpec := models.CEIPUpdateSpec{}
	statusVal := "DISABLE"
	updateSpec.Status = &statusVal
	params.CEIPUpdateSpec = &updateSpec

	_, ceipAccepted, err := apiClient.CEIP.UpdateCEIPStatus(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	if vcfClient.WaitForTask(ceipAccepted.Payload.ID) != nil {
		return diag.FromErr(err)
	}

	log.Printf("%s: Delete complete", d.Id())
	d.SetId("")
	return nil
}
