// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/terraform-provider-vcf/internal/sddc"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

var dvSwitchVersions = []string{"7.0.0", "7.0.2", "7.0.3", "8.0.0", "8.0.3"}

func ResourceVcfInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcfInstanceCreate,
		ReadContext:   resourceVcfInstanceRead,
		UpdateContext: resourceVcfInstanceUpdate,
		DeleteContext: resourceVcfInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Hour),
		},
		Schema: resourceVcfInstanceSchema(),
	}
}

func resourceVcfInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Description:  "Client string that identifies an SDDC by name or instance name. Used for management domain name. Can contain only letters, numbers and the following symbols: '-'. Example: \"sfo01-m01\", Length 3-20 characters",
			Required:     true,
			ValidateFunc: validationutils.ValidateSddcId,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "SDDC creation Task status",
			Computed:    true,
		},
		"creation_timestamp": {
			Type:        schema.TypeString,
			Description: "SDDC Task creation timestamp",
			Computed:    true,
		},
		"sddc_manager_fqdn": {
			Type:        schema.TypeString,
			Description: "FQDN of the resulting SDDC Manager",
			Computed:    true,
		},
		"sddc_manager_id": {
			Type:        schema.TypeString,
			Description: "ID of the resulting SDDC Manager",
			Computed:    true,
		},
		"sddc_manager_version": {
			Type:        schema.TypeString,
			Description: "Version of the resulting SDDC Manager",
			Computed:    true,
		},
		"ceip_enabled": {
			Type:        schema.TypeBool,
			Description: "Enable VCF Customer Experience Improvement Program",
			Optional:    true,
		},
		"fips_enabled": {
			Type:        schema.TypeBool,
			Description: "Enable Federal Information Processing Standards",
			Optional:    true,
		},
		"cluster": sddc.GetSddcClusterSchema(),
		"dns":     sddc.GetDnsSchema(),
		"dvs":     sddc.GetDvsSchema(),
		"dv_switch_version": {
			Type:         schema.TypeString,
			Description:  "The version of the distributed virtual switches to be used. One among: 7.0.0, 7.0.2, 7.0.3, 8.0.0, 8.0.3",
			Required:     true,
			ValidateFunc: validation.StringInSlice(dvSwitchVersions, false),
		},
		"esx_license": {
			Type:      schema.TypeString,
			Sensitive: true,
			Optional:  true,
		},
		"host": sddc.GetSddcHostSchema(),
		"management_pool_name": {
			Type:        schema.TypeString,
			Description: "A string identifying the network pool associated with the management domain",
			Required:    true,
		},
		"network": sddc.GetNetworkSpecsSchema(),
		"nsx":     sddc.GetNsxSpecSchema(),
		"ntp_servers": {
			Type:        schema.TypeList,
			Description: "List of NTP servers",
			Required:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"psc":          sddc.GetPscSchema(),
		"sddc_manager": sddc.GetSddcManagerSchema(),
		"security":     sddc.GetSecuritySchema(),
		"skip_esx_thumbprint_validation": {
			Type:        schema.TypeBool,
			Description: "Skip ESXi thumbprint validation",
			Required:    true,
		},
		"task_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "workflowconfig/workflowspec-ems.json",
		},
		"vcenter":    sddc.GetVcenterSchema(),
		"vsan":       sddc.GetVsanSchema(),
		"vx_manager": sddc.GetVxManagerSchema(),
	}
}

func buildSddcSpec(data *schema.ResourceData, apiClient *vcf.ClientWithResponses) *vcf.SddcSpec {
	sddcSpec := &vcf.SddcSpec{}
	if rawCeipEnabled, ok := data.GetOk("ceip_enabled"); ok {
		sddcSpec.CeipEnabled = utils.ToBoolPointer(rawCeipEnabled)
	}
	if clusterSpec, ok := data.GetOk("cluster"); ok {
		sddcSpec.ClusterSpec = sddc.GetSddcClusterSpecFromSchema(clusterSpec.([]interface{}), apiClient)
	}
	if dnsSpec, ok := data.GetOk("dns"); ok {
		spec := sddc.GetDnsSpecFromSchema(dnsSpec.([]interface{}))
		// TODO throw error and make dns mandatory
		sddcSpec.DnsSpec = *spec
	}
	if dvsSpecs, ok := data.GetOk("dvs"); ok {
		sddcSpec.DvsSpecs = sddc.GetDvsSpecsFromSchema(dvsSpecs.([]interface{}))
	}
	if dvSwitchVersion, ok := data.GetOk("dv_switch_version"); ok {
		sddcSpec.DvSwitchVersion = utils.ToStringPointer(dvSwitchVersion)
	}
	if esxLicense, ok := data.GetOk("esx_license"); ok {
		sddcSpec.EsxLicense = utils.ToStringPointer(esxLicense)
	}
	if fipsEnabled, ok := data.GetOk("fips_enabled"); ok {
		sddcSpec.FipsEnabled = utils.ToBoolPointer(fipsEnabled)
	}
	if hostSpecs, ok := data.GetOk("host"); ok {
		sddcSpec.HostSpecs = sddc.GetSddcHostSpecsFromSchema(hostSpecs.([]interface{}))
	}
	if managementPoolName, ok := data.GetOk("management_pool_name"); ok {
		sddcSpec.ManagementPoolName = utils.ToStringPointer(managementPoolName)
	}
	if networkSpecs, ok := data.GetOk("network"); ok {
		sddcSpec.NetworkSpecs = sddc.GetNetworkSpecsBindingFromSchema(networkSpecs.([]interface{}))
	}
	if nsxSpec, ok := data.GetOk("nsx"); ok {
		sddcSpec.NsxtSpec = sddc.GetNsxSpecFromSchema(nsxSpec.([]interface{}))
	}
	if ntpServers, ok := data.GetOk("ntp_servers"); ok {
		sddcSpec.NtpServers = utils.ToStringSlice(ntpServers.([]interface{}))
	}
	if pscSpecs, ok := data.GetOk("psc"); ok {
		sddcSpec.PscSpecs = sddc.GetPscSpecsFromSchema(pscSpecs.([]interface{}))
	}
	if sddcID, ok := data.GetOk("instance_id"); ok {
		sddcSpec.SddcId = sddcID.(string)
	}
	if sddcManagerSpec, ok := data.GetOk("sddc_manager"); ok {
		sddcSpec.SddcManagerSpec = sddc.GetSddcManagerSpecFromSchema(sddcManagerSpec.([]interface{}))
	}
	if securitySpec, ok := data.GetOk("security"); ok {
		sddcSpec.SecuritySpec = sddc.GetSecuritySpecSchema(securitySpec.([]interface{}))
	}
	if skipEsxThumbPrintValidation, ok := data.GetOk("skip_esx_thumbprint_validation"); ok {
		sddcSpec.SkipEsxThumbprintValidation = utils.ToBoolPointer(skipEsxThumbPrintValidation)
	}
	if taskName, ok := data.GetOk("task_name"); ok {
		sddcSpec.TaskName = utils.ToStringPointer(taskName)
	}
	if vcenterSpec, ok := data.GetOk("vcenter"); ok {
		if spec := sddc.GetVcenterSpecFromSchema(vcenterSpec.([]interface{})); spec != nil {
			sddcSpec.VcenterSpec = *spec
		}
	}
	if vsanSpec, ok := data.GetOk("vsan"); ok {
		sddcSpec.VsanSpec = sddc.GetVsanSpecFromSchema(vsanSpec.([]interface{}))
	}
	if vxManagerSpec, ok := data.GetOk("vx_manager"); ok {
		sddcSpec.VxManagerSpec = sddc.GetVxManagerSpecFromSchema(vxManagerSpec.([]interface{}))
	}
	return sddcSpec
}

func resourceVcfInstanceCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	sddcSpec := buildSddcSpec(data, client.ApiClient)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	bringUpID, diags := invokeBringupWorkflow(ctx, client, sddcSpec, bringUpInfo)
	if diags != nil {
		return diags
	}

	diags = waitForBringupProcess(ctx, bringUpID, client)
	if diags != nil {
		return diags
	}

	return resourceVcfInstanceRead(ctx, data, meta)
}

func resourceVcfInstanceRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	bringupId := bringUpInfo.Id

	data.SetId(*bringupId)
	_ = data.Set("status", bringUpInfo.Status)
	_ = data.Set("creation_timestamp", bringUpInfo.CreationTimestamp)

	sddcManagerInfo, err := getSddcManagerInfo(ctx, *bringupId, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	_ = data.Set("sddc_manager_fqdn", sddcManagerInfo.Fqdn)
	_ = data.Set("sddc_manager_id", sddcManagerInfo.Id)
	_ = data.Set("sddc_manager_version", sddcManagerInfo.Version)

	return nil
}
func resourceVcfInstanceUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// no op
	return resourceVcfInstanceRead(ctx, data, meta)
}
func resourceVcfInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// no op
	return nil
}

func invokeBringupWorkflow(ctx context.Context, client *api_client.CloudBuilderClient, sddcSpec *vcf.SddcSpec, lastBringup *vcf.SddcTask) (string, diag.Diagnostics) {
	var bringUpId string
	if lastBringup != nil && *lastBringup.Status != "COMPLETED_WITH_SUCCESS" {
		bringUpId = *lastBringup.Id
		diags := validateBringupSpec(ctx, client, sddcSpec)
		if diags != nil {
			return bringUpId, diags
		}

		res, err := client.ApiClient.RetrySddcWithResponse(ctx, bringUpId, *sddcSpec)

		sddcTask, vcfErr := api_client.GetResponseAs[vcf.SddcTask](res)
		if err != nil {
			return "", diag.FromErr(err)
		}
		if vcfErr != nil {
			api_client.LogError(vcfErr)
		} else {
			bringUpId = *sddcTask.Id
		}
	} else {
		diags := validateBringupSpec(ctx, client, sddcSpec)
		if diags != nil {
			return bringUpId, diags
		}

		res, err := client.ApiClient.StartBringupWithResponse(ctx, *sddcSpec)
		if err != nil {
			return "", diag.FromErr(err)
		}
		sddcTask, vcfErr := api_client.GetResponseAs[vcf.SddcTask](res)
		if err != nil {
			return "", diag.FromErr(err)
		}
		if vcfErr != nil {
			api_client.LogError(vcfErr)
		} else {
			bringUpId = *sddcTask.Id
		}

		tflog.Info(ctx, fmt.Sprintf("Bring-Up workflow with ID %s has started", bringUpId))
	}
	return bringUpId, nil
}

func waitForBringupProcess(ctx context.Context, bringUpID string, client *api_client.CloudBuilderClient) diag.Diagnostics {
	for {
		task, err := getBringUp(ctx, bringUpID, client)
		if err != nil {
			return diag.FromErr(err)
		}

		if *task.Status == "IN_PROGRESS" {
			time.Sleep(20 * time.Second)
			continue
		}

		if *task.Status == "COMPLETED_WITH_FAILURE" {
			err := fmt.Errorf("task with ID = %s , Name: %q is in state %s", bringUpID, *task.Name, *task.Status)
			return diag.FromErr(err)
		}

		return nil
	}
}

func getLastBringUp(ctx context.Context, client *api_client.CloudBuilderClient) (*vcf.SddcTask, error) {
	retrieveAllSddcsResp, err := client.ApiClient.GetBringupTasksWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfSddcTask](retrieveAllSddcsResp)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}
	if page != nil && len(*page.Elements) > 0 {
		elements := *page.Elements
		return &(elements)[0], nil
	}
	return nil, nil
}

func validateBringupSpec(ctx context.Context, client *api_client.CloudBuilderClient, sddcSpec *vcf.SddcSpec) diag.Diagnostics {
	bringupParams := &vcf.ValidateBringupSpecParams{
		Redo: utils.ToBoolPointer(true),
	}
	validateSpecRes, err := client.ApiClient.ValidateBringupSpecWithResponse(ctx, bringupParams, *sddcSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	validationResult, vcfErr := api_client.GetResponseAs[vcf.Validation](validateSpecRes)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HasValidationFailed(validationResult) {
		return validationutils.ConvertValidationResultToDiag(validationResult)
	}
	for {
		getValidationResponse, err := client.ApiClient.GetBringupValidationWithResponse(ctx, *validationResult.Id)
		if err != nil {
			return validationutils.ConvertVcfErrorToDiag(err)
		}
		validationResult, vcfErr = api_client.GetResponseAs[vcf.Validation](getValidationResponse)
		if vcfErr != nil {
			api_client.LogError(vcfErr)
			return diag.FromErr(errors.New(*vcfErr.Message))
		}
		if validationResult != nil && validationutils.HaveValidationChecksFinished(*validationResult.ValidationChecks) {
			break
		}
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HasValidationFailed(validationResult) {
		return validationutils.ConvertValidationResultToDiag(validationResult)
	}

	return nil
}

func getBringUp(ctx context.Context, bringupId string, client *api_client.CloudBuilderClient) (*vcf.SddcTask, error) {
	retrieveSddcResponse, err := client.ApiClient.GetBringupTaskByIDWithResponse(ctx, bringupId)
	if err != nil {
		return nil, err
	}
	sddcTask, vcfErr := api_client.GetResponseAs[vcf.SddcTask](retrieveSddcResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}
	return sddcTask, nil
}

func getSddcManagerInfo(ctx context.Context, bringupId string, client *api_client.CloudBuilderClient) (*vcf.SddcManagerInfo, error) {
	getSddcManagerInfoResponse, err := client.ApiClient.GetSddcManagerInfoWithResponse(ctx, bringupId)
	if err != nil {
		return nil, err
	}
	info, vcfErr := api_client.GetResponseAs[vcf.SddcManagerInfo](getSddcManagerInfoResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}

	return info, nil
}
