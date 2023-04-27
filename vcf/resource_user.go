package vcf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain of the user",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the user",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role to assign to the user",
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API Key for a service user",
			},
		},
	}
}

func resourceUserCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		ok, err := client.Users.GETRoles(nil)
		if err != nil {
			log.Println("error = ", err)
			return diag.FromErr(err)
		}

		found := false
		for _, role := range ok.Payload.Elements {
			if *role.Name == roleNameVal {
				user.Role = &models.RoleReference{ID: role.ID}
				found = true
				break
			}
		}

		if !found {
			log.Println("Did not find role ", roleNameVal)
			return diag.Errorf(fmt.Sprintf("Did not find role %s", roleNameVal))
		}
	}
	// AddUsers(params *AddUsersParams, opts ...ClientOption) (*AddUsersOK, *AddUsersCreated, error)
	params.Users = []*models.User{&user}

	_, created, err := client.Users.AddUsers(params)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	log.Println("created = ", created)
	createdUser := created.Payload.Elements[0]
	d.SetId(createdUser.ID)
	_ = d.Set("api_key", createdUser.APIKey)

	return nil
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*SddcManagerClient).ApiClient

	id := d.Id()

	ok, err := client.Users.GETUsers(nil)
	if err != nil {
		log.Println("error = ", err)
		return diag.FromErr(err)
	}

	jsonp, _ := json.MarshalIndent(ok.Payload, " ", " ")
	log.Println(string(jsonp))

	// Check if the resource with the known id exists
	for _, user := range ok.Payload.Elements {
		if user.ID == id {
			_ = d.Set("api_key", user.APIKey)
			return nil
		}
	}

	// Did not find the resource, set ID to ""
	log.Println("Did not find user with id ", id)
	d.SetId("")

	return nil
}

/**
 * Updating users is not supported in VCF API.
 */
func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceUserRead(ctx, d, meta)
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
