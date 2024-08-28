// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CertificateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain of the resource certificate",
				Computed:    true,
			},
			"expiration_status": {
				Type:        schema.TypeString,
				Description: "Certificate expiry status. One among: ACTIVE, ABOUT_TO_EXPIRE, EXPIRED",
				Computed:    true,
			},
			"certificate_error": {
				Type:        schema.TypeString,
				Description: "Error if certificate cannot be fetched. Example: Status : NOT_TRUSTED, Message : Certificate Expired",
				Computed:    true,
			},
			"issued_by": {
				Type:        schema.TypeString,
				Description: "The certificate authority that issued the certificate",
				Computed:    true,
			},
			"issued_to": {
				Type:        schema.TypeString,
				Description: "To whom the certificate is issued",
				Computed:    true,
			},
			"key_size": {
				Type:        schema.TypeString,
				Description: "The key size of the certificate",
				Computed:    true,
			},
			"not_after": {
				Type:        schema.TypeString,
				Description: "The timestamp after which certificate is not valid",
				Computed:    true,
			},
			"not_before": {
				Type:        schema.TypeString,
				Description: "The timestamp before which certificate is not valid",
				Computed:    true,
			},
			"number_of_days_to_expire": {
				Type:        schema.TypeInt,
				Description: "Number of days left for the certificate to expire",
				Computed:    true,
			},
			"pem_encoded": {
				Type:        schema.TypeString,
				Description: "The PEM encoded certificate content",
				Sensitive:   true,
				Computed:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key of the certificate",
				Computed:    true,
			},
			"public_key_algorithm": {
				Type:        schema.TypeString,
				Description: "The public key algorithm of the certificate",
				Computed:    true,
			},
			"serial_number": {
				Type:        schema.TypeString,
				Description: "The serial number of the certificate",
				Computed:    true,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Description: "Algorithm used to sign the certificate",
				Computed:    true,
			},
			"subject": {
				Type:        schema.TypeString,
				Description: "Complete distinguished name to which the certificate is issued",
				Computed:    true,
			},
			"subject_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_cn(common name) of the certificate.",
			},
			"subject_ou": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_ou(org unit) of the certificate.",
			},
			"subject_org": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_org of the certificate.",
			},
			"subject_locality": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_locality of the certificate.",
			},
			"subject_st": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_st(state) of the certificate.",
			},
			"subject_country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subject_country of the certificate.",
			},
			"subject_alternative_name": {
				Type:        schema.TypeList,
				Description: "The alternative names to which the certificate is issued",
				Computed:    true,
				Elem:        schema.TypeString,
			},
			"thumbprint": {
				Type:        schema.TypeString,
				Description: "Thumbprint generated using certificate content",
				Computed:    true,
			},
			"thumbprint_algorithm": {
				Type:        schema.TypeString,
				Description: "Algorithm used to generate thumbprint",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The X.509 version of the certificate",
				Computed:    true,
			},
		},
	}
}
