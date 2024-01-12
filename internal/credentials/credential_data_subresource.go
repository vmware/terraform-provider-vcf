package credentials

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CredentialDataSubresource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID hash of the credential",
				Required:    true,
			},
			"user_name": {
				Type:        schema.TypeString,
				Description: "The username of the account to which the credential belong",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account to which the credential belong",
				Computed:    true,
			},
			"account_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "One among USER, SYSTEM, SERVICE",
			},
			"credential_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the credential. For example FTP, SSH, etc.",
			},
			"auto_rotate_frequency_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "After how many days the credentials will be auto rotated. One among 30, 60, 90",
			},
			"auto_rotate_next_schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of the next rotation",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the credential is created",
			},
			"modification_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the credentials are changed",
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID hash of the resource to which the credential belongs",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCF domain to which the resource belongs",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ip address of the resource related to the credential",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the resource as registered in SDDC Manager inventory",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the resource.One among ESXI, VCENTER, PSC, NSX_MANAGER, NSX_CONTROLLER, NSX_EDGE, NSXT_MANAGER, VRLI, VROPS, VRA, WSA, VRSLCM, VXRAIL_MANAGER, NSX_ALB, BACKUP",
						},
					},
				},
			},
		},
	}
}
