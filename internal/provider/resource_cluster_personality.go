// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/vcf-sdk-go/client/personalities"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

const (
	uploadModeReferred = "REFERRED"
)

func ResourceClusterPersonality() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterPersonalityCreate,
		ReadContext:   resourceClusterPersonalityRead,
		DeleteContext: resourceClusterPersonalityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name for the personality",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			// Maybe use Domain name & cluster name as inputs and work out the identifiers
			"vcenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The identifier of the vCenter Server where the source cluster resides",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The identifier for the source cluster",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
		},
	}
}

func resourceClusterPersonalityCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	mode := uploadModeReferred
	name := data.Get("name").(string)
	clusterId := data.Get("cluster_id").(string)
	vcenterId := data.Get("vcenter_id").(string)
	spec := models.PersonalityUploadSpec{
		Name:       name,
		UploadMode: &mode,
		UploadSpecReferredMode: &models.PersonalityUploadSpecReferred{
			ClusterID: &clusterId,
			VCenterID: &vcenterId,
		},
	}

	if _, task, err := client.Personalities.UploadPersonality(&personalities.UploadPersonalityParams{PersonalityUploadSpec: &spec}); err != nil {
		return diag.FromErr(err)
	} else if err := meta.(*api_client.SddcManagerClient).WaitForTaskComplete(ctx, task.Payload.ID, false); err != nil {
		return diag.FromErr(err)
	}

	if personalitiesResp, err := client.Personalities.GetPersonalities(&personalities.GetPersonalitiesParams{
		PersonalityName: &name,
	}); err != nil {
		return diag.FromErr(err)
	} else if len(personalitiesResp.Payload.Elements) == 0 {
		return diag.Errorf("Personality %s not found", name)
	} else {
		data.SetId(*personalitiesResp.Payload.Elements[0].PersonalityID)
	}

	return nil
}

func resourceClusterPersonalityRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	// Just check if the personality exists. There are no computed attributes.
	if _, err := client.Personalities.GetPersonality(&personalities.GetPersonalityParams{
		PersonalityID: data.Id(),
	}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceClusterPersonalityDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	id := data.Id()
	if _, err := client.Personalities.DeletePersonality(&personalities.DeletePersonalityParams{
		PersonalityID: &id,
	}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
