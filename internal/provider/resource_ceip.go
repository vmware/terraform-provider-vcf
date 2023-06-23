/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/client/ceip"
	"github.com/vmware/vcf-sdk-go/models"
	"strings"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DISABLED_STATE = "DISABLED"
	ENABLED_STATE  = "ENABLED"

	ENABLE_API_PARAM  = "ENABLE"
	DISABLE_API_PARAM = "DISABLE"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User provided CEIP operation. One among: ENABLED, DISABLED",
				ValidateFunc: validation.StringInSlice([]string{ENABLED_STATE, DISABLED_STATE}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
		},
	}
}

func resourceCeipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCeipUpdate(ctx, d, meta)
}

func resourceCeipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*SddcManagerClient).ApiClient

	ceipResult, err := apiClient.CEIP.GetCEIPStatus(nil)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	d.SetId(ceipResult.Payload.InstanceID)
	return nil
}

func resourceCeipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	params := ceip.NewUpdateCEIPStatusParams()
	updateSpec := models.CEIPUpdateSpec{}

	if status, ok := d.GetOk("status"); ok {
		statusVal := status.(string)
		// the VCF PATCH API requires the params "ENABLE/DISABLE" while the resource states are "ENABLED/DISABLED"
		var enableApiParam string
		if statusVal == ENABLED_STATE {
			enableApiParam = ENABLE_API_PARAM
		} else if statusVal == DISABLED_STATE {
			enableApiParam = DISABLE_API_PARAM
		}
		updateSpec.Status = &enableApiParam
	}

	params.CEIPUpdateSpec = &updateSpec
	_, ceipAccepted, err := apiClient.CEIP.UpdateCEIPStatus(params)
	if err != nil {
		tflog.Error(ctx, err.Error())
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
func resourceCeipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	params := ceip.NewUpdateCEIPStatusParams()
	updateSpec := models.CEIPUpdateSpec{}
	statusVal := DISABLE_API_PARAM
	updateSpec.Status = &statusVal
	params.CEIPUpdateSpec = &updateSpec

	_, ceipAccepted, err := apiClient.CEIP.UpdateCEIPStatus(params)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	if vcfClient.WaitForTask(ceipAccepted.Payload.ID) != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
