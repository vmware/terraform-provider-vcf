package vcf

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the resource configuration of the VCF provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sddc_manager_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC manager username",
			},
			"sddc_manager_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC manager password.",
			},
			"sddc_manager_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC manager host.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username, isSetUsername := data.GetOk("sddc_manager_username")
	password, isSetPassword := data.GetOk("sddc_manager_password")
	host, isSetHost := data.GetOk("sddc_manager_host")

	if !isSetUsername || !isSetPassword || !isSetHost {
		return nil, diag.Errorf("SDDC Manager username, password and hostname must be provided")
	}

	var newClient = NewSddcManagerClient(username.(string), password.(string), host.(string))
	newClient.Connect()
	return newClient, nil
}
