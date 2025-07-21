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
	"github.com/vmware/vcf-sdk-go/installer"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/terraform-provider-vcf/internal/sddc"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func ResourceVcfInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcfInstanceCreate,
		ReadContext:   resourceVcfInstanceRead,
		UpdateContext: resourceVcfInstanceUpdate,
		DeleteContext: resourceVcfInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Hour), // it takes a while
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
		"host":    sddc.GetSddcHostSchema(),
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
		"sddc_manager": sddc.GetSddcManagerSchema(),
		"security":     sddc.GetSecuritySchema(),
		"skip_esx_thumbprint_validation": {
			Type:        schema.TypeBool,
			Description: "Skip ESXi thumbprint validation",
			Required:    true,
		},
		"vcenter":                     sddc.GetVcenterSchema(),
		"vsan":                        sddc.GetVsanSchema(),
		"automation":                  sddc.GetVcfAutomationSchema(),
		"operations":                  sddc.GetVcfOperationsSchema(),
		"operations_collector":        sddc.GetVcfOperationsCollectorSchema(),
		"operations_fleet_management": sddc.GetVcfOperationsFleetManagementSchema(),
		"version": {
			Type:        schema.TypeString,
			Description: "VCF version",
			Optional:    true,
		},
	}
}

func buildSddcSpec(data *schema.ResourceData) *installer.SddcSpec {
	sddcSpec := &installer.SddcSpec{}
	if rawCeipEnabled, ok := data.GetOk("ceip_enabled"); ok {
		sddcSpec.CeipEnabled = utils.ToBoolPointer(rawCeipEnabled)
	}
	if clusterSpec, ok := data.GetOk("cluster"); ok {
		sddcSpec.ClusterSpec = sddc.GetSddcClusterSpecFromSchema(clusterSpec.([]interface{}))
	}
	if vsanSpec, ok := data.GetOk("vsan"); ok {
		// TODO support NFS & VMFS if necessary
		sddcSpec.DatastoreSpec = &installer.SddcDatastoreSpec{}
		sddcSpec.DatastoreSpec.VsanSpec = sddc.GetVsanSpecFromSchema(vsanSpec.([]interface{}))
	}
	if dnsSpec, ok := data.GetOk("dns"); ok {
		spec := sddc.GetDnsSpecFromSchema(dnsSpec.([]interface{}))
		sddcSpec.DnsSpec = *spec
	}
	if dvsSpecs, ok := data.GetOk("dvs"); ok {
		sddcSpec.DvsSpecs = sddc.GetDvsSpecsFromSchema(dvsSpecs.([]interface{}))
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
		ntpServersValue := utils.ToStringSlice(ntpServers.([]interface{}))
		sddcSpec.NtpServers = &ntpServersValue
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
	if vcenterSpec, ok := data.GetOk("vcenter"); ok {
		if spec := sddc.GetVcenterSpecFromSchema(vcenterSpec.([]interface{})); spec != nil {
			sddcSpec.VcenterSpec = *spec
		}
	}
	if automationSpec, ok := data.GetOk("automation"); ok {
		if spec := sddc.GetVcfAutomationSpecFromSchema(automationSpec.([]interface{})); spec != nil {
			sddcSpec.VcfAutomationSpec = spec
		}
	}
	if operationsSpec, ok := data.GetOk("operations"); ok {
		if spec := sddc.GetVcfOperationsSpecFromSchema(operationsSpec.([]interface{})); spec != nil {
			sddcSpec.VcfOperationsSpec = spec
		}
	}
	if operationsCollectorSpec, ok := data.GetOk("operations_collector"); ok {
		if spec := sddc.GetVcfOperationsCollectorSpecFromSchema(operationsCollectorSpec.([]interface{})); spec != nil {
			sddcSpec.VcfOperationsCollectorSpec = spec
		}
	}
	if operationsFleetManagementSpec, ok := data.GetOk("operations_fleet_management"); ok {
		if spec := sddc.GetVcfOperationsFleetManagementSpecFromSchema(operationsFleetManagementSpec.([]interface{})); spec != nil {
			sddcSpec.VcfOperationsFleetManagementSpec = spec
		}
	}
	if version, ok := data.GetOk("version"); ok {
		sddcSpec.Version = utils.ToStringPointer(version)
	}
	return sddcSpec
}

func resourceVcfInstanceCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.InstallerClient)

	sddcSpec := buildSddcSpec(data)

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
	client := meta.(*api_client.InstallerClient)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	bringupId := bringUpInfo.Id

	data.SetId(*bringupId)
	_ = data.Set("status", bringUpInfo.Status)
	_ = data.Set("creation_timestamp", bringUpInfo.CreationTimestamp)

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

func invokeBringupWorkflow(ctx context.Context, client *api_client.InstallerClient, sddcSpec *installer.SddcSpec, lastBringup *installer.SddcTask) (string, diag.Diagnostics) {
	var bringUpId string
	if lastBringup != nil && *lastBringup.Status != "COMPLETED_WITH_SUCCESS" {
		bringUpId = *lastBringup.Id
		diags := validateBringupSpec(ctx, client, sddcSpec)
		if diags != nil {
			return bringUpId, diags
		}

		res, err := client.ApiClient.RetrySddcWithResponse(ctx, bringUpId, nil, *sddcSpec)
		if err != nil {
			return "", diag.FromErr(err)
		}
		sddcTask, vcfErr := api_client.GetResponseAs[installer.SddcTask](res)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
		} else {
			bringUpId = *sddcTask.Id
		}
	} else {
		diags := validateBringupSpec(ctx, client, sddcSpec)
		if diags != nil {
			return bringUpId, diags
		}

		res, err := client.ApiClient.DeploySddcWithResponse(ctx, nil, *sddcSpec)
		if err != nil {
			return "", diag.FromErr(err)
		}
		sddcTask, vcfErr := api_client.GetResponseAs[installer.SddcTask](res)
		if err != nil {
			return "", diag.FromErr(err)
		}
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
		} else {
			bringUpId = *sddcTask.Id
		}

		tflog.Info(ctx, fmt.Sprintf("Bring-Up workflow with ID %s has started", bringUpId))
	}
	return bringUpId, nil
}

func waitForBringupProcess(ctx context.Context, bringUpID string, client *api_client.InstallerClient) diag.Diagnostics {
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
			logFailedTask(ctx, task)
			return diag.FromErr(err)
		}

		return nil
	}
}

func getLastBringUp(ctx context.Context, client *api_client.InstallerClient) (*installer.SddcTask, error) {
	retrieveAllSddcsResp, err := client.ApiClient.GetSddcTasksWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	page, vcfErr := api_client.GetResponseAs[installer.PageOfSddcTask](retrieveAllSddcsResp)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}
	if page != nil && len(*page.Elements) > 0 {
		elements := *page.Elements
		return &(elements)[0], nil
	}
	return nil, nil
}

func validateBringupSpec(ctx context.Context, client *api_client.InstallerClient, sddcSpec *installer.SddcSpec) diag.Diagnostics {
	validateSpecRes, err := client.ApiClient.ValidateSddcSpecWithResponse(ctx, *sddcSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	validationResult, vcfErr := api_client.GetResponseAs[installer.Validation](validateSpecRes)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	vcfValidationResult := api_client.ConvertToVcfValidation(*validationResult)

	if validationutils.HasValidationFailed(&vcfValidationResult) {
		return validationutils.ConvertValidationResultToDiag(&vcfValidationResult)
	}
	for {
		getValidationResponse, err := client.ApiClient.GetSddcSpecValidationWithResponse(ctx, *validationResult.Id)
		if err != nil {
			return validationutils.ConvertVcfErrorToDiag(err)
		}
		validationResult, vcfErr = api_client.GetResponseAs[installer.Validation](getValidationResponse)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
			return diag.FromErr(errors.New(*vcfErr.Message))
		}
		vcfValidationResult = api_client.ConvertToVcfValidation(*validationResult)
		if validationutils.HaveValidationChecksFinished(*vcfValidationResult.ValidationChecks) {
			break
		}
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HasValidationFailed(&vcfValidationResult) {
		return validationutils.ConvertValidationResultToDiag(&vcfValidationResult)
	}

	return nil
}

func getBringUp(ctx context.Context, bringupId string, client *api_client.InstallerClient) (*installer.SddcTask, error) {
	retrieveSddcResponse, err := client.ApiClient.GetSddcTaskByIDWithResponse(ctx, bringupId)
	if err != nil {
		return nil, err
	}
	sddcTask, vcfErr := api_client.GetResponseAs[installer.SddcTask](retrieveSddcResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}
	return sddcTask, nil
}

func logFailedTask(ctx context.Context, task *installer.SddcTask) {
	tflog.Error(ctx, fmt.Sprintf("task with ID = %s, Name: %q is in state %s", *task.Id, *task.Name, *task.Status))

	if task.SddcSubTasks != nil {
		for _, subtask := range *task.SddcSubTasks {
			tflog.Error(ctx, fmt.Sprintf("subtask %q is in state %s", *subtask.Name, *subtask.Status))
		}
	}
}
