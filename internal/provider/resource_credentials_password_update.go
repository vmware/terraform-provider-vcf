package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/credentials"
	"log"
	"time"
)

func ResourceCredentialsUpdate() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceCredentialsPasswordUpdateRead,
		CreateContext: resourceCredentialsPasswordUpdateCreate,
		DeleteContext: resourceCredentialsPasswordUpdateDelete,
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource which credentials will be updated",
				ForceNew:    true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the resource which credentials will be updated",
				ForceNew:    true,
			},
			"credentials": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The credentials that should be updated",
				MinItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type(s) of the account. One among: SSO, SSH, API, FTP, AUDIT",
							ValidateFunc: validation.StringInSlice(
								credentials.AllCredentialTypes(), true),
						},
						"user_name": {
							Type:        schema.TypeString,
							Description: "The user name of the account.",
							Required:    true,
						},
						"password": {
							Type:        schema.TypeString,
							Description: "The password for the account.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"once_only": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				ForceNew:    true,
				Description: "If set to true operation is executed only once otherwise rotation is done each time.",
			},
			"last_update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of the last password rotation.",
			},
		},
	}
}

func resourceCredentialsPasswordUpdateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient
	creds, err := credentials.ReadCredentials(ctx, data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	dataCreds := data.Get("credentials").([]interface{})
	for _, dataCredential := range dataCreds {
		entry := dataCredential.(map[string]interface{})
		for _, cred := range creds {
			username := entry["user_name"].(string)
			credentialType := entry["credential_type"].(string)
			if *cred.Username == username && *cred.CredentialType == credentialType {
				entry["password"] = cred.Password
			}
		}
	}

	_ = data.Set("credentials", dataCreds)
	id, err := credentials.CreatePasswordChangeID(data, credentials.Rotate)
	if err != nil {
		return diag.Errorf("error during id generation %s", err)
	}

	data.SetId(id)

	return nil
}

func resourceCredentialsPasswordUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("last_update_time").(string) != "" && d.Get("once_only").(bool) {
		log.Print("[DEBUG] Skipping password rotation")
		return nil
	}

	err := credentials.UpdatePasswords(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("last_update_time", time.Now().Format(time.RFC3339))
	id, err := credentials.CreatePasswordChangeID(d, credentials.Update)
	if err != nil {
		return diag.Errorf("error during id generation %s", err)
	}
	d.SetId(id)

	return nil
}

func resourceCredentialsPasswordUpdateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
