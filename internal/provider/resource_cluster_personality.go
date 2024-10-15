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

	task, err := client.UploadPersonalityWithResponse(ctx, spec)

	if err != nil {
		return diag.FromErr(err)
	}
	if task.StatusCode() != 202 {
		vcfError := api_client.GetError(task.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	if err := meta.(*api_client.SddcManagerClient).WaitForTaskComplete(ctx, *task.JSON202.Id, false); err != nil {
		return diag.FromErr(err)
	}

	personalitiesResp, err := client.GetPersonalitiesWithResponse(ctx, &vcf.GetPersonalitiesParams{
		PersonalityName: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	} else if personalitiesResp.StatusCode() != 200 {
		vcfError := api_client.GetError(personalitiesResp.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	if personalitiesResp.JSON200.Elements != nil && len(*personalitiesResp.JSON200.Elements) == 0 {
		return diag.Errorf("Personality %s not found", name)
	} else {
		personalities := *personalitiesResp.JSON200.Elements
		data.SetId(*personalities[0].PersonalityId)
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
		for _, vc := range *vcs.JSON200.Elements {
			if vc.Domain.Id == domainId {
				return vc.Id, nil
			}
		}
	}

	return nil, fmt.Errorf("vcenter for domain %s not found", domainId)
}
