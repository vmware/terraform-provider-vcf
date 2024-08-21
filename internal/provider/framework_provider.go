// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

type FrameworkProviderModel struct {
	SddcManagerUsername types.String `tfsdk:"sddc_manager_username"`
	SddcManagerPassword types.String `tfsdk:"sddc_manager_password"`
	SddcManagerHost     types.String `tfsdk:"sddc_manager_host"`

	CloudBuilderUsername types.String `tfsdk:"cloud_builder_username"`
	CloudBuilderPassword types.String `tfsdk:"cloud_builder_password"`
	CloudBuilderHost     types.String `tfsdk:"cloud_builder_host"`

	AllowUnverifiedTls types.Bool `tfsdk:"allow_unverified_tls"`
}

type FrameworkProvider struct {
	// The clients are exposed for the purpose of running the existing tests
	// Individual resources should obtain access to these in their Configure methods
	// via the ConfigureRequest
	SddcManagerClient  *api_client.SddcManagerClient
	CloudBuilderClient *api_client.CloudBuilderClient
}

func New() provider.Provider {
	return &FrameworkProvider{}
}

func (frameworkProvider *FrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sddc_manager_username": schema.StringAttribute{
				Optional:    true,
				Description: "Username to authenticate to SDDC Manager",
				Validators: []validator.String{
					getSddcManagerConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("sddc_manager_password"),
							path.MatchRoot("sddc_manager_host"),
						}...),
				},
			},
			"sddc_manager_password": schema.StringAttribute{
				Optional:    true,
				Description: "Password to authenticate to SDDC Manager",
				Validators: []validator.String{
					getSddcManagerConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("sddc_manager_username"),
							path.MatchRoot("sddc_manager_host"),
						}...),
				},
			},
			"sddc_manager_host": schema.StringAttribute{
				Optional:    true,
				Description: "Fully qualified domain name or IP address of the SDDC Manager",
				Validators: []validator.String{
					getSddcManagerConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("sddc_manager_username"),
							path.MatchRoot("sddc_manager_password"),
						}...),
				},
			},
			"cloud_builder_username": schema.StringAttribute{
				Optional:    true,
				Description: "Username to authenticate to CloudBuilder",
				Validators: []validator.String{
					getCloudBuilderConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("cloud_builder_password"),
							path.MatchRoot("cloud_builder_host"),
						}...),
				},
			},
			"cloud_builder_password": schema.StringAttribute{
				Optional:    true,
				Description: "Password to authenticate to CloudBuilder",
				Validators: []validator.String{
					getCloudBuilderConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("cloud_builder_username"),
							path.MatchRoot("cloud_builder_host"),
						}...),
				},
			},
			"cloud_builder_host": schema.StringAttribute{
				Optional:    true,
				Description: "Fully qualified domain name or IP address of the CloudBuilder",
				Validators: []validator.String{
					getCloudBuilderConflictsValidator(),
					stringvalidator.AlsoRequires(
						path.Expressions{
							path.MatchRoot("cloud_builder_username"),
							path.MatchRoot("cloud_builder_password"),
						}...),
				},
			},
			"allow_unverified_tls": schema.BoolAttribute{
				Optional:    true,
				Description: "If set, VMware VCF client will permit unverifiable TLS certificates.",
			},
		},
	}
}

func (frameworkProvider *FrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "vcf"
}

func (frameworkProvider *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &ResourceNetworkPool{}
		},
	}
}

func (frameworkProvider *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (frameworkProvider *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	var data FrameworkProviderModel

	res.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sddcManagerUsername := getAttributeValue(data.SddcManagerUsername.ValueString(), constants.VcfTestUsername).(string)

	if sddcManagerUsername != "" {
		// Connect to SDDC Manager
		client := api_client.NewSddcManagerClient(
			sddcManagerUsername,
			getAttributeValue(data.SddcManagerPassword.ValueString(), constants.VcfTestPassword).(string),
			getAttributeValue(data.SddcManagerHost.ValueString(), constants.VcfTestUrl).(string),
			getAttributeValue(data.AllowUnverifiedTls.ValueBool(), constants.VcfTestAllowUnverifiedTls).(bool),
		)

		if err := client.Connect(); err != nil {
			res.Diagnostics.Append(diag.NewErrorDiagnostic("Failed to connect to the SDDC Manager", err.Error()))
		}

		frameworkProvider.SddcManagerClient = client
		res.ResourceData = client
	} else {
		// Connect to Cloud Builder
		client := api_client.NewCloudBuilderClient(
			getAttributeValue(data.CloudBuilderUsername.ValueString(), constants.CloudBuilderTestUsername).(string),
			getAttributeValue(data.CloudBuilderPassword.ValueString(), constants.CloudBuilderTestPassword).(string),
			getAttributeValue(data.CloudBuilderHost.ValueString(), constants.CloudBuilderTestUrl).(string),
			getAttributeValue(data.AllowUnverifiedTls.ValueBool(), constants.VcfTestAllowUnverifiedTls).(bool),
		)

		frameworkProvider.CloudBuilderClient = client
		res.ResourceData = client
	}
}

func getAttributeValue[T string | bool](data T, envVar string) interface{} {
	if envVal := os.Getenv(envVar); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			return val
		}

		return envVal
	}

	return data
}

func getSddcManagerConflictsValidator() validator.String {
	return stringvalidator.ConflictsWith(
		path.Expressions{
			path.MatchRoot("cloud_builder_username"),
			path.MatchRoot("cloud_builder_password"),
			path.MatchRoot("cloud_builder_host"),
		}...)
}

func getCloudBuilderConflictsValidator() validator.String {
	return stringvalidator.ConflictsWith(
		path.Expressions{
			path.MatchRoot("sddc_manager_username"),
			path.MatchRoot("sddc_manager_password"),
			path.MatchRoot("sddc_manager_host"),
		}...)
}
