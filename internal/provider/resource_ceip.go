// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
)

const (
	DisabledState = "DISABLED"
	EnabledState  = "ENABLED"

	EnableApiParam  = "ENABLE"
	DisableApiParam = "DISABLE"
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
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User provided CEIP operation. One among: ENABLED, DISABLED",
				ValidateFunc: validation.StringInSlice([]string{EnabledState, DisabledState}, true),
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
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	ceipResult, err := apiClient.GetCeipStatusWithResponse(ctx)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	d.SetId(*ceipResult.JSON200.InstanceId)
	return nil
}

func resourceCeipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	var enableApiParam string
	if status, ok := d.GetOk("status"); ok {
		statusVal := status.(string)
		// the VCF PATCH API requires the params "ENABLE/DISABLE" while the resource states are "ENABLED/DISABLED"
		if statusVal == EnabledState {
			enableApiParam = EnableApiParam
		} else if statusVal == DisabledState {
			enableApiParam = DisableApiParam
		}
	}

	res, err := apiClient.SetCeipStatusWithResponse(ctx, vcf.CeipUpdateSpec{Status: enableApiParam})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	if res.StatusCode() != 202 {
		vcfError := api_client.GetError(res.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	if vcfClient.WaitForTask(ctx, *res.JSON202.Id) != nil {
		return diag.FromErr(err)
	}

	return resourceCeipRead(ctx, d, meta)
}

/**
 * Mapping deletion of ceip resource to disabling ceip.
 */
func resourceCeipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	updateSpec := vcf.CeipUpdateSpec{}
	statusVal := DisableApiParam
	updateSpec.Status = statusVal

	ceipAccepted, err := apiClient.SetCeipStatusWithResponse(ctx, updateSpec)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	if vcfClient.WaitForTask(ctx, *ceipAccepted.JSON202.Id) != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
