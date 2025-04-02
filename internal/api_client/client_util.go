package api_client

import (
	"github.com/vmware/vcf-sdk-go/installer"
	"github.com/vmware/vcf-sdk-go/vcf"
)

func ConvertToVcfError(err installer.Error) vcf.Error {
	var causes *[]vcf.ErrorCause
	if err.Causes != nil {
		mappedCauses := mapValues(*err.Causes, convertToVcfErrorCause)
		causes = &mappedCauses
	}

	var nestedErrors *[]vcf.Error
	if err.NestedErrors != nil {
		mappedNestedErrors := mapValues(*err.NestedErrors, ConvertToVcfError)
		nestedErrors = &mappedNestedErrors
	}

	return vcf.Error{
		Arguments:          err.Arguments,
		Causes:             causes,
		Context:            err.Context,
		ErrorCode:          err.ErrorCode,
		ErrorType:          err.ErrorType,
		Message:            err.Message,
		NestedErrors:       nestedErrors,
		ReferenceToken:     err.ReferenceToken,
		RemediationMessage: err.RemediationMessage,
	}
}

func ConvertToVcfValidation(val installer.Validation) vcf.Validation {
	var validationChecks *[]vcf.ValidationCheck
	if val.ValidationChecks != nil {
		mappedValidationChecks := mapValues(*val.ValidationChecks, convertToVcfValidationCheck)
		validationChecks = &mappedValidationChecks
	}

	return vcf.Validation{
		AdditionalProperties: val.AdditionalProperties,
		Description:          val.Description,
		ExecutionStatus:      val.ExecutionStatus,
		Id:                   val.Id,
		ResultStatus:         val.ResultStatus,
		ValidationChecks:     validationChecks,
	}
}

func convertToVcfErrorCause(cause installer.ErrorCause) vcf.ErrorCause {
	return vcf.ErrorCause{
		Message: cause.Message,
		Type:    cause.Type,
	}
}

func convertToVcfValidationCheck(check installer.ValidationCheck) vcf.ValidationCheck {
	var errorResponse *vcf.Error
	if check.ErrorResponse != nil {
		value := ConvertToVcfError(*check.ErrorResponse)
		errorResponse = &value
	}
	return vcf.ValidationCheck{
		Description:   check.Description,
		ErrorResponse: errorResponse,
		ResultStatus:  check.ResultStatus,
		Severity:      check.Severity,
	}
}

func mapValues[T, K interface{}](src []T, f func(T) K) []K {
	r := make([]K, len(src))
	for i, v := range src {
		r[i] = f(v)
	}
	return r
}
