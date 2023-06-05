/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/client/users"
	"github.com/vmware/vcf-sdk-go/models"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			Create: schema.DefaultTimeout(2 * time.Hour),
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
					"USER", "GROUP", "SERVICE"}, false),
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
	client := meta.(*SddcManagerClient).ApiClient
	log.Println(d)
	params := users.NewAddUsersParams()
	user := models.User{}

	if name, ok := d.GetOk("name"); ok {
		nameVal := name.(string)
		user.Name = &nameVal
	}

	if domain, ok := d.GetOk("domain"); ok {
		user.Domain = domain.(string)
	}

	if roleType, ok := d.GetOk("type"); ok {
		roleTypeVal := roleType.(string)
		user.Type = &roleTypeVal
	}

	if roleName, ok := d.GetOk("role_name"); ok {
		roleNameVal := roleName.(string)

		roleResult, err := client.Users.GetRoles(nil)
		if err != nil {
			log.Println("error = ", err)
			return diag.FromErr(err)
		}

		roleFound := false
		for _, role := range roleResult.Payload.Elements {
			if *role.Name == roleNameVal {
				user.Role = &models.RoleReference{ID: role.ID}
				roleFound = true
				break
			}
		}

		if !roleFound {
			log.Println("Did not find role ", roleNameVal)
			return diag.Errorf(fmt.Sprintf("Did not find role %s", roleNameVal))
		}
	}
	params.Users = []*models.User{&user}

	_, created, err := client.Users.AddUsers(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	createdUser := created.Payload.Elements[0]
	d.SetId(createdUser.ID)
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*SddcManagerClient).ApiClient

	id := d.Id()

	ok, err := client.Users.GetUsers(nil)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	// Check if the resource with the known id exists
	for _, user := range ok.Payload.Elements {
		if user.ID == id {
			_ = d.Set("api_key", user.APIKey)
			_ = d.Set("creation_timestamp", user.CreationTimestamp)
			return nil
		}
	}

	return nil
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*SddcManagerClient).ApiClient

	params := users.NewDeleteUserParams()
	params.ID = d.Id()

	log.Println(params)
	_, err := client.Users.DeleteUser(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Printf("%s: Delete complete", d.Id())
	d.SetId("")
	return nil
}
