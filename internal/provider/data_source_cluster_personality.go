// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/vcf-sdk-go/vcf"
)

type DataSourceClusterPersonality struct {
	client *vcf.ClientWithResponses
}

type DataSourceClusterPersonalityModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *DataSourceClusterPersonality) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api_client.SddcManagerClient).ApiClient
}

func (d *DataSourceClusterPersonality) Metadata(_ context.Context, _ datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = "vcf_cluster_personality"
}

func (d *DataSourceClusterPersonality) Schema(_ context.Context, _ datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the personality",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (d *DataSourceClusterPersonality) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data DataSourceClusterPersonalityModel
	res.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	tflog.Debug(ctx, fmt.Sprintf("Looking for personality '%s'", data.Name.ValueString()))
	response, err := d.client.GetPersonalitiesWithResponse(ctx, &vcf.GetPersonalitiesParams{
		PersonalityName: data.Name.ValueStringPointer(),
	})
	if err != nil {
		res.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}

	responsePage, vcfErr := api_client.GetResponseAs[vcf.PageOfPersonality](response)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		res.Diagnostics.Append(diag.NewErrorDiagnostic(*vcfErr.ErrorType, *vcfErr.Message))
		return
	}

	if responsePage == nil || responsePage.Elements == nil || len(*responsePage.Elements) == 0 {
		res.Diagnostics.Append(diag.NewErrorDiagnostic(
			fmt.Sprintf("personality with name '%s' not found", data.Name.ValueString()), ""))
		return
	}

	personalities := *responsePage.Elements
	personality := personalities[0]

	tflog.Debug(ctx, fmt.Sprintf("Personality '%s' found, reading data", data.Name.ValueString()))

	data.ID = types.StringValue(*personality.PersonalityId)
	data.Name = types.StringValue(*personality.PersonalityName)

	res.State.Set(ctx, data)
}
