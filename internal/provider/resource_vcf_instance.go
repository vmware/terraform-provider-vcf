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

var dvSwitchVersions = []string{"6.0.0", "6.5.0", "7.0.0"}
var excludedComponents = []string{"Foundation", "VsphereHostProfiles", "LogInsight", "NSX", "VrealizeNetwork", "VSAN", "VSANCleanup",
	"VROPS", "VRA", "DRDeployment", "DRConfiguration", "ConfigurationBackup", "VRB", "VRSLCM", "Inventory", "UMDS", "EsxThumbprintValidation",
	"AVN", "CEIP", "Backup", "EBGP"}

func ResourceVcfInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcfInstanceCreate,
		ReadContext:   resourceVcfInstanceRead,
		UpdateContext: resourceVcfInstanceUpdate,
		DeleteContext: resourceVcfInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"ceip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Customer Experience Improvement Program is to be enabled",
			},
			"certificates_passphrase": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster": sddc.GetSddcClusterSchema(),
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns": sddc.GetDnsSchema(),
			"dvs": sddc.GetDvsSchema(),
			"dv_switch_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(dvSwitchVersions, false),
			},
			"esx_license": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"excluded_components": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(excludedComponents, false),
				},
				Optional: true,
			},
			"fips_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"host": sddc.GetSddcHostSchema(),
			"management_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": sddc.GetNetworkSpecsSchema(),
			"nsx":     sddc.GetNsxSpecSchema(),
			"ntp_servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},
			"psc": sddc.GetPscSchema(),
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Client string that identifies an SDDC by name or instance name. Used for management domain name. Can contain only letters, numbers and the following symbols: '-'. Example: \"sfo01-m01\", Length 3-20 characters",
				Required:    true,
			},
			"sddc_manager": sddc.GetSddcManagerSchema(),
			"security":     sddc.GetSecuritySchema(),
			"should_cleanup_vsan": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skip_esx_thumbprint_validation": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "workflowconfig/workflowspec-ems.json",
			},
			"vcenter":    sddc.GetVcenterSchema(),
			"vsan":       sddc.GetVsanSchema(),
			"vx_manager": sddc.GetVxManagerSchema(),
		},
	}
}

func resourceVcfInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	sddcSpec := &models.SDDCSpec{}
	if rawCeipEnabled, ok := d.GetOk("ceip_enabled"); ok {
		ceipEnabled := rawCeipEnabled.(bool)
		sddcSpec.CEIPEnabled = ceipEnabled
	}
	if certificatesPassphrase, ok := d.GetOk("certificates_passphrase"); ok {
		sddcSpec.CertificatesPassphrase = certificatesPassphrase.(string)
	}
	if clusterSpec, ok := d.GetOk("cluster"); ok {
		sddcSpec.ClusterSpec = sddc.GetSddcClusterSpecFromSchema(clusterSpec.([]interface{}))
	}
	if dnsSpec, ok := d.GetOk("dns"); ok {
		sddcSpec.DNSSpec = sddc.GetDnsSpecFromSchema(dnsSpec.([]interface{}))
	}
	if dvsSpecs, ok := d.GetOk("dvs"); ok {
		sddcSpec.DvsSpecs = sddc.GetDvsSpecsFromSchema(dvsSpecs.([]interface{}))
	}
	if dvSwitchVersion, ok := d.GetOk("dv_switch_version"); ok {
		sddcSpec.DvSwitchVersion = dvSwitchVersion.(string)
	}
	if esxLicense, ok := d.GetOk("esx_license"); ok {
		sddcSpec.EsxLicense = esxLicense.(string)
	}
	if rawFipsEnabled, ok := d.GetOk("fips_enabled"); ok {
		fipsEnabled := rawFipsEnabled.(bool)
		sddcSpec.FIPSEnabled = fipsEnabled
	}
	if hostSpecs, ok := d.GetOk("host"); ok {
		sddcSpec.HostSpecs = sddc.GetSddcHostSpecsFromSchema(hostSpecs.([]interface{}))
	}
	if managementPoolName, ok := d.GetOk("management_pool_name"); ok {
		sddcSpec.ManagementPoolName = managementPoolName.(string)
	}
	if networkSpecs, ok := d.GetOk("network"); ok {
		sddcSpec.NetworkSpecs = sddc.GetNetworkSpecsBindingFromSchema(networkSpecs.([]interface{}))
	}
	if nsxSpec, ok := d.GetOk("nsx"); ok {
		sddcSpec.NSXTSpec = sddc.GetNsxSpecFromSchema(nsxSpec.([]interface{}))
	}
	if ntpServers, ok := d.GetOk("ntp_servers"); ok {
		sddcSpec.NtpServers = utils.ToStringSlice(ntpServers.([]interface{}))
	}
	if pscSpecs, ok := d.GetOk("psc"); ok {
		sddcSpec.PscSpecs = sddc.GetPscSpecsFromSchema(pscSpecs.([]interface{}))
	}
	if sddcID, ok := d.GetOk("instance_id"); ok {
		sddcSpec.SDDCID = utils.ToStringPointer(sddcID)
	}
	if sddcManagerSpec, ok := d.GetOk("sddc_manager"); ok {
		sddcSpec.SDDCManagerSpec = sddc.GetSddcManagerSpecFromSchema(sddcManagerSpec.([]interface{}))
	}
	if securitySpec, ok := d.GetOk("security"); ok {
		sddcSpec.SecuritySpec = sddc.GetSecuritySpecSchema(securitySpec.([]interface{}))
	}
	if rawShouldCleanupVsan, ok := d.GetOk("should_cleanup_vsan"); ok {
		shouldCleanupVsan := rawShouldCleanupVsan.(bool)
		sddcSpec.ShouldCleanupVSAN = shouldCleanupVsan
	}
	if skipEsxThumbPrintValidation, ok := d.GetOk("skip_esx_thumbprint_validation"); ok {
		sddcSpec.SkipEsxThumbprintValidation = skipEsxThumbPrintValidation.(bool)
	}
	if taskName, ok := d.GetOk("task_name"); ok {
		sddcSpec.TaskName = utils.ToStringPointer(taskName)
	}
	if vcenterSpec, ok := d.GetOk("vcenter"); ok {
		sddcSpec.VcenterSpec = sddc.GetVcenterSpecFromSchema(vcenterSpec.([]interface{}))
	}
	if vsanSpec, ok := d.GetOk("vsan"); ok {
		sddcSpec.VSANSpec = sddc.GetVsanSpecFromSchema(vsanSpec.([]interface{}))
	}
	if vxManagerSpec, ok := d.GetOk("vx_manager"); ok {
		sddcSpec.VxManagerSpec = sddc.GetVxManagerSpecFromSchema(vxManagerSpec.([]interface{}))
	}
	if excludedComponentsAttribute, ok := d.GetOk("excluded_components"); ok {
		sddcSpec.ExcludedComponents = utils.ToStringSlice(excludedComponentsAttribute.([]interface{}))
	}

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
func resourceVcfInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.CloudBuilderClient)

	bringUpInfo, err := getLastBringUp(ctx, client)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}

	d.SetId(bringUpInfo.ID)
	_ = d.Set("status", bringUpInfo.Status)
	_ = d.Set("creation_timestamp", bringUpInfo.CreationTimestamp)
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
