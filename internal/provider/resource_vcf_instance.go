/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/terraform-provider-vcf/internal/sddc"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	sddc_api "github.com/vmware/vcf-sdk-go/client/sddc"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

var dvSwitchVersions = []string{"7.0.0", "7.0.2", "7.0.3"}

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

// TODO add support for "subscriptionLicensing" property in future releases.
func resourceVcfInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Description:  "Client string that identifies an SDDC by name or instance name. Used for management domain name. Can contain only letters, numbers and the following symbols: '-'. Example: \"sfo01-m01\", Length 3-20 characters",
			Required:     true,
			ValidateFunc: validation_utils.ValidateSddcId,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_timestamp": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sddc_manager_fqdn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sddc_manager_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sddc_manager_version": {
			Type:     schema.TypeString,
			Computed: true,
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
			Description:  "The version of the distributed virtual switches to be used. One among: 7.0.0, 7.0.2, 7.0.3",
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
				ValidateFunc: validation.IsIPAddress,
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

func resourceVcfInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	sddcSpec := buildSddcSpec(d)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	bringUpID, diags := invokeBringupWorkflow(ctx, client, sddcSpec, bringUpInfo)
	if diags != nil {
		return diags
	}

	return waitForBringupProcess(ctx, bringUpID, client)
}

func buildSddcSpec(data *schema.ResourceData) *models.SDDCSpec {
	sddcSpec := &models.SDDCSpec{}
	if rawCeipEnabled, ok := data.GetOk("ceip_enabled"); ok {
		ceipEnabled := rawCeipEnabled.(bool)
		sddcSpec.CEIPEnabled = ceipEnabled
	}
	if clusterSpec, ok := data.GetOk("cluster"); ok {
		sddcSpec.ClusterSpec = sddc.GetSddcClusterSpecFromSchema(clusterSpec.([]interface{}))
	}
	if dnsSpec, ok := data.GetOk("dns"); ok {
		sddcSpec.DNSSpec = sddc.GetDnsSpecFromSchema(dnsSpec.([]interface{}))
	}
	if dvsSpecs, ok := data.GetOk("dvs"); ok {
		sddcSpec.DvsSpecs = sddc.GetDvsSpecsFromSchema(dvsSpecs.([]interface{}))
	}
	if dvSwitchVersion, ok := data.GetOk("dv_switch_version"); ok {
		sddcSpec.DvSwitchVersion = dvSwitchVersion.(string)
	}
	if esxLicense, ok := data.GetOk("esx_license"); ok {
		sddcSpec.EsxLicense = esxLicense.(string)
	}
	if rawFipsEnabled, ok := data.GetOk("fips_enabled"); ok {
		fipsEnabled := rawFipsEnabled.(bool)
		sddcSpec.FIPSEnabled = fipsEnabled
	}
	if hostSpecs, ok := data.GetOk("host"); ok {
		sddcSpec.HostSpecs = sddc.GetSddcHostSpecsFromSchema(hostSpecs.([]interface{}))
	}
	if managementPoolName, ok := data.GetOk("management_pool_name"); ok {
		sddcSpec.ManagementPoolName = managementPoolName.(string)
	}
	if networkSpecs, ok := data.GetOk("network"); ok {
		sddcSpec.NetworkSpecs = sddc.GetNetworkSpecsBindingFromSchema(networkSpecs.([]interface{}))
	}
	if nsxSpec, ok := data.GetOk("nsx"); ok {
		sddcSpec.NSXTSpec = sddc.GetNsxSpecFromSchema(nsxSpec.([]interface{}))
	}
	if ntpServers, ok := data.GetOk("ntp_servers"); ok {
		sddcSpec.NtpServers = utils.ToStringSlice(ntpServers.([]interface{}))
	}
	if pscSpecs, ok := data.GetOk("psc"); ok {
		sddcSpec.PscSpecs = sddc.GetPscSpecsFromSchema(pscSpecs.([]interface{}))
	}
	if sddcID, ok := data.GetOk("instance_id"); ok {
		sddcSpec.SDDCID = utils.ToStringPointer(sddcID)
	}
	if sddcManagerSpec, ok := data.GetOk("sddc_manager"); ok {
		sddcSpec.SDDCManagerSpec = sddc.GetSddcManagerSpecFromSchema(sddcManagerSpec.([]interface{}))
	}
	if securitySpec, ok := data.GetOk("security"); ok {
		sddcSpec.SecuritySpec = sddc.GetSecuritySpecSchema(securitySpec.([]interface{}))
	}
	if skipEsxThumbPrintValidation, ok := data.GetOk("skip_esx_thumbprint_validation"); ok {
		sddcSpec.SkipEsxThumbprintValidation = skipEsxThumbPrintValidation.(bool)
	}
	if taskName, ok := data.GetOk("task_name"); ok {
		sddcSpec.TaskName = utils.ToStringPointer(taskName)
	}
	if vcenterSpec, ok := data.GetOk("vcenter"); ok {
		sddcSpec.VcenterSpec = sddc.GetVcenterSpecFromSchema(vcenterSpec.([]interface{}))
	}
	if vsanSpec, ok := data.GetOk("vsan"); ok {
		sddcSpec.VSANSpec = sddc.GetVsanSpecFromSchema(vsanSpec.([]interface{}))
	}
	if vxManagerSpec, ok := data.GetOk("vx_manager"); ok {
		sddcSpec.VxManagerSpec = sddc.GetVxManagerSpecFromSchema(vxManagerSpec.([]interface{}))
	}
	return sddcSpec
}

func resourceVcfInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	bringupId := bringUpInfo.ID

	d.SetId(bringupId)
	_ = d.Set("status", bringUpInfo.Status)
	_ = d.Set("creation_timestamp", bringUpInfo.CreationTimestamp)

	sddcManagerInfo, err := getSddcManagerInfo(ctx, bringupId, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	_ = d.Set("sddc_manager_fqdn", sddcManagerInfo.Fqdn)
	_ = d.Set("sddc_manager_id", sddcManagerInfo.ID)
	_ = d.Set("sddc_manager_version", sddcManagerInfo.Version)

	return nil
}
func resourceVcfInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// no op
	return resourceVcfInstanceRead(ctx, d, meta)
}
func resourceVcfInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// no op
	return nil
}

func invokeBringupWorkflow(ctx context.Context, client *api_client.CloudBuilderClient, sddcSpec *models.SDDCSpec, lastBringup *models.SDDCTask) (string, diag.Diagnostics) {
	var bringUpID string
	if lastBringup != nil && lastBringup.Status != "COMPLETED_WITH_SUCCESS" {
		bringUpID = lastBringup.ID
		diags := validateBringupSpec(ctx, client, sddcSpec, true)
		if diags != nil {
			return bringUpID, diags
		}

		retryBringupParams := sddc_api.NewRetrySDDCParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout).WithID(bringUpID).WithSDDCSpec(sddcSpec)
		okResponse, acceptedResponse, err := client.ApiClient.SDDC.RetrySDDC(retryBringupParams)
		if okResponse != nil {
			bringUpID = okResponse.Payload.ID
		}
		if acceptedResponse != nil {
			bringUpID = acceptedResponse.Payload.ID
		}
		if err != nil {
			return "", diag.FromErr(err)
		}
	} else {
		diags := validateBringupSpec(ctx, client, sddcSpec, false)
		if diags != nil {
			return bringUpID, diags
		}

		bringupParams := sddc_api.NewCreateSDDCParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout).WithSDDCSpec(sddcSpec)

		okResponse, acceptedResponse, err := client.ApiClient.SDDC.CreateSDDC(bringupParams)
		if okResponse != nil {
			bringUpID = okResponse.Payload.ID
		}
		if acceptedResponse != nil {
			bringUpID = acceptedResponse.Payload.ID
		}
		if err != nil {
			return "", diag.FromErr(err)
		}

		tflog.Info(ctx, fmt.Sprintf("Bring-Up workflow with ID %s has started", bringUpID))
	}
	return bringUpID, nil
}

func waitForBringupProcess(ctx context.Context, bringUpID string, client *api_client.CloudBuilderClient) diag.Diagnostics {
	for {
		task, err := getBringUp(ctx, bringUpID, client)
		if err != nil {
			return diag.FromErr(err)
		}

		if task.Status == "IN_PROGRESS" {
			time.Sleep(20 * time.Second)
			continue
		}

		if task.Status == "COMPLETED_WITH_FAILURE" {
			errorMsg := fmt.Sprintf("Task with ID = %s , Name: %q is in state %s", bringUpID, task.Name, task.Status)

			tflog.Error(ctx, errorMsg)
			return diag.FromErr(fmt.Errorf(errorMsg))
		}

		return nil
	}
}

func getLastBringUp(ctx context.Context, client *api_client.CloudBuilderClient) (*models.SDDCTask, error) {
	retrieveAllSddcsResp, err := client.ApiClient.SDDC.RetrieveAllSddcs(
		sddc_api.NewRetrieveAllSddcsParamsWithTimeout(constants.DefaultVcfApiCallTimeout).WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if len(retrieveAllSddcsResp.Payload.Elements) > 0 {
		return retrieveAllSddcsResp.Payload.Elements[0], nil
	}
	return nil, fmt.Errorf("no bringups executed, cannot determine last successful bringup")
}

func validateBringupSpec(ctx context.Context, client *api_client.CloudBuilderClient, sddcSpec *models.SDDCSpec, isRetry bool) diag.Diagnostics {
	validateSddcSpec := sddc_api.NewValidateSDDCSpecParams().WithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).WithSDDCSpec(sddcSpec).WithRedo(utils.ToBoolPointer(isRetry))

	var validationResponse *models.Validation
	okResponse, acceptedResponse, err := client.ApiClient.SDDC.ValidateSDDCSpec(validateSddcSpec)
	if okResponse != nil {
		validationResponse = okResponse.Payload
	}
	if acceptedResponse != nil {
		validationResponse = acceptedResponse.Payload
	}
	if err != nil {
		return validation_utils.ConvertVcfErrorToDiag(err)
	}
	if validation_utils.HasValidationFailed(validationResponse) {
		return validation_utils.ConvertValidationResultToDiag(validationResponse)
	}
	return nil
}

func getBringUp(ctx context.Context, bringupId string, client *api_client.CloudBuilderClient) (*models.SDDCTask, error) {
	retrieveSddcResponse, err := client.ApiClient.SDDC.RetrieveSDDC(
		sddc_api.NewRetrieveSDDCParamsWithContext(ctx).WithID(bringupId).WithTimeout(constants.DefaultVcfApiCallTimeout))
	if err != nil {
		return nil, err
	}
	return retrieveSddcResponse.Payload, nil
}

func getSddcManagerInfo(ctx context.Context, bringupId string, client *api_client.CloudBuilderClient) (*models.SDDCManagerInfo, error) {
	getSddcManagerInfoResponse, err := client.ApiClient.SDDC.GetSDDCManagerInfo(
		sddc_api.NewGetSDDCManagerInfoParamsWithContext(ctx).WithID(bringupId).WithTimeout(constants.DefaultVcfApiCallTimeout))
	if err != nil {
		return nil, err
	}
	return getSddcManagerInfoResponse.Payload, nil
}
