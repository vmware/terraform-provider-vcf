// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Updating users is not supported in VCF API.
				Description: "The name of the user",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain of the user",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the user. One of: USER, GROUP, SERVICE",
				ValidateFunc: validation.StringInSlice([]string{
					"USER", "GROUP", "SERVICE"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the role to assign to the user",
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API Key for a service user",
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient
	user := vcf.User{}

	if name, ok := d.GetOk("name"); ok {
		user.Name = name.(string)
	}

	if domain, ok := d.GetOk("domain"); ok {
		user.Domain = utils.ToStringPointer(domain)
	}

	if roleType, ok := d.GetOk("type"); ok {
		user.Type = roleType.(string)
	}

	if roleName, ok := d.GetOk("role_name"); ok {
		roleNameVal := roleName.(string)

		roleResult, err := client.GetRolesWithResponse(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		if roleResult.StatusCode() != 200 {
			vcfError := api_client.GetError(roleResult.Body)
			api_client.LogError(vcfError)
			return diag.FromErr(errors.New(*vcfError.Message))
		}

		roleFound := false
		for _, role := range *roleResult.JSON200.Elements {
			if *role.Name == roleNameVal {
				user.Role = vcf.RoleReference{Id: *role.Id}
				roleFound = true
				break
			}
		}

		if !roleFound {
			return diag.Errorf("role not found: %s", roleNameVal)
		}
	}

	created, err := client.AddUsersWithResponse(ctx, []vcf.User{user})
	if err != nil {
		return diag.FromErr(err)
	}
	if created.StatusCode() != 201 {
		vcfError := api_client.GetError(created.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	createdUser := (*created.JSON201.Elements)[0]
	d.SetId(*createdUser.Id)
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	id := d.Id()

	ok, err := client.GetUsersWithResponse(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if the resource with the known id exists
	for _, user := range *ok.JSON200.Elements {
		if *user.Id == id {
			_ = d.Set("api_key", user.ApiKey)
			_ = d.Set("creation_timestamp", user.CreationTimestamp)
			return nil
		}
	}

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	res, err := client.RemoveUserWithResponse(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if res.StatusCode() >= 400 {
		vcfError := api_client.GetError(res.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	log.Printf("%s: Delete complete", d.Id())
	d.SetId("")
	return nil
}
