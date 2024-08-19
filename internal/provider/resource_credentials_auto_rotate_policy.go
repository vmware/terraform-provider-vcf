package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/credentials"
)

func ResourceCredentialsAutoRotatePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialsAutoRotatePolicyCreate,
		ReadContext:   resourceCredentialsAutoRotatePolicyRead,
		DeleteContext: resourceCredentialsAutoRotatePolicyDelete,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the resource which credentials autorotate policy will be managed",
				ForceNew:    true,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the resource which credentials autorotate policy will be managed",
				ForceNew:    true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of the resource which credentials autorotate policy will be managed",
				ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile("^ESXI$"), "Schedule auto rotate not supported for the ESXI entity type."),
				ForceNew:     true,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account name which autorotate policy will be managed",
				ForceNew:    true,
			},
			"enable_auto_rotation": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable the automatic credential rotation",
				ForceNew:    true,
			},
			"auto_rotate_days": {
				Type:         schema.TypeInt,
				Default:      credentials.AutorotateDays30,
				Optional:     true,
				Description:  fmt.Sprintf("The number of days after the credentials will be automatically rotated. Must be between %v and %v", credentials.AutoRotateDaysMin, credentials.AutorotateDaysMax),
				ValidateFunc: validation.All(validation.IntAtLeast(credentials.AutoRotateDaysMin), validation.IntAtMost(credentials.AutorotateDays90)),
				ForceNew:     true,
			},
			"auto_rotate_next_schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The next time automatic rotation will be started",
			},
		},
	}
}

func resourceCredentialsAutoRotatePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := credentials.CreateAutoRotatePolicy(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialsAutoRotatePolicyRead(ctx, d, meta)
}

func resourceCredentialsAutoRotatePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient
	matchedCredentials, err := credentials.ReadCredentials(ctx, d, apiClient)
	matchedCredentials = filterCredentials(d.Get("user_name").(string), d.Get("resource_id").(string), matchedCredentials)

	if err != nil {
		return diag.FromErr(err)
	}

	lenCredentials := len(matchedCredentials)
	if lenCredentials != 1 {
		err := fmt.Errorf("only one credential expected, received %v", lenCredentials)
		return diag.FromErr(err)
	}

	id, err := createAutorotateID(d)
	if err != nil {
		return diag.Errorf("error during id generation %s", err)
	}

	d.SetId(id)

	if matchedCredentials[0].AutoRotatePolicy != nil {
		_ = d.Set("enable_auto_rotation", true)
		_ = d.Set("auto_rotate_days", matchedCredentials[0].AutoRotatePolicy.FrequencyInDays)
		_ = d.Set("auto_rotate_next_schedule", matchedCredentials[0].AutoRotatePolicy.NextSchedule)
	} else {
		_ = d.Set("enable_auto_rotation", false)
		_ = d.Set("auto_rotate_days", 0)
		_ = d.Set("auto_rotate_next_schedule", "")
	}

	return nil
}

func resourceCredentialsAutoRotatePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := credentials.RemoveAutoRotatePolicy(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createAutorotateID(data *schema.ResourceData) (string, error) {
	params := []string{
		data.Get("resource_id").(string),
		data.Get("resource_name").(string),
		data.Get("resource_type").(string),
		data.Get("user_name").(string),
	}

	return credentials.HashFields(params)
}

func filterCredentials(userName, resourceId string, creds []*models.Credential) []*models.Credential {
	result := make([]*models.Credential, 0)
	for _, cred := range creds {
		if *cred.Username == userName && cred.Resource != nil && *cred.Resource.ResourceID == resourceId {
			result = append(result, cred)
		}
	}

	return result
}
