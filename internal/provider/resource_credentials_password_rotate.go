// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/credentials"
)

func ResourceCredentialsRotate() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceCredentialsPasswordRotationRead,
		CreateContext: resourceCredentialsPasswordRotationCreate,
		DeleteContext: resourceCredentialsPasswordRotationDelete,
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource which credentials will be rotated",
				ForceNew:    true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of the resource which credentials will be rotated",
				ValidateFunc: validation.StringInSlice(credentials.AllResourceTypes(), false),
				ForceNew:     true,
			},
			"credentials": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The credentials that should be rotated",
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
							Computed:    true,
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
			"last_rotate_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of the last password rotation.",
			},
		},
	}
}

func resourceCredentialsPasswordRotationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceCredentialsPasswordRotationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("last_rotate_time").(string) != "" && d.Get("once_only").(bool) {
		log.Print("[DEBUG] Skipping password rotation")
		return nil
	}

	err := credentials.RotatePasswords(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("last_rotate_time", time.Now().Format(time.RFC3339))

	return resourceCredentialsPasswordRotationRead(ctx, d, meta)
}

func resourceCredentialsPasswordRotationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
