// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
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
			"domain_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The identifier of the domain which contains the vcenter and source cluster",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The identifier of the source cluster within the vCenter server (e.g. domain-c1)",
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

	vcenterId, err := getVcenterId(data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := data.Get("cluster_id").(string)

	spec := vcf.PersonalityUploadSpec{
		Name:       &name,
		UploadMode: mode,
		UploadSpecReferredMode: &vcf.PersonalityUploadSpecReferred{
			ClusterId: clusterId,
			VCenterId: *vcenterId,
		},
	}

	uploadPersonalityTask, err := client.UploadPersonalityWithResponse(ctx, spec)

	if err != nil {
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](uploadPersonalityTask)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if err = api_client.NewTaskTracker(ctx, client, *task.Id).WaitForTask(); err != nil {
		return diag.FromErr(err)
	}

	personalitiesResp, err := client.GetPersonalitiesWithResponse(ctx, &vcf.GetPersonalitiesParams{
		PersonalityName: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	personalities, vcfErr := api_client.GetResponseAs[vcf.PageOfPersonality](personalitiesResp)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if personalities.Elements == nil || len(*personalities.Elements) == 0 {
		return diag.Errorf("Personality %s not found", name)
	} else {
		elements := *personalities.Elements
		data.SetId(*elements[0].PersonalityId)
	}

	return nil
}

func resourceClusterPersonalityRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	// Just check if the personality exists. There are no computed attributes.
	if _, err := client.GetPersonalityWithResponse(ctx, data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceClusterPersonalityDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	id := data.Id()
	if _, err := client.DeletePersonalityWithResponse(ctx, &vcf.DeletePersonalityParams{
		PersonalityId: &id,
	}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getVcenterId(data *schema.ResourceData, meta interface{}) (*string, error) {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	domainId := data.Get("domain_id").(string)

	if vcs, err := client.GetVcentersWithResponse(context.Background(), nil); err != nil {
		return nil, err
	} else {
		page, vcfErr := api_client.GetResponseAs[vcf.PageOfVcenter](vcs)
		if vcfErr != nil {
			api_client.LogError(vcfErr)
			return nil, errors.New(*vcfErr.Message)
		}
		for _, vc := range *page.Elements {
			if vc.Domain.Id == domainId {
				return vc.Id, nil
			}
		}
	}

	return nil, fmt.Errorf("vcenter for domain %s not found", domainId)
}
