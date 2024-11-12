package api_client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/vmware/vcf-sdk-go/vcf"
)

func TestGetResponseAs_pos(t *testing.T) {
	taskId := "id"
	model := vcf.Task{Id: &taskId}
	body, _ := json.Marshal(model)
	httpResponse := http.Response{StatusCode: 200}
	response := vcf.GetTaskResponse{
		Body:         body,
		HTTPResponse: &httpResponse,
	}

	result, err := GetResponseAs[vcf.Task](response)

	if err != nil {
		t.Fatal("received an unexpected error", err)
	}

	if result == nil {
		t.Fatal("response is nil")
	}

	if *result.Id != taskId {
		t.Fatal("response does not contain correct payload")
	}
}

func TestGetResponseAs_neg(t *testing.T) {
	message := "message"
	model := vcf.Error{Message: &message}
	body, _ := json.Marshal(model)
	httpResponse := http.Response{StatusCode: 400}
	response := vcf.GetTaskResponse{
		Body:         body,
		HTTPResponse: &httpResponse,
	}

	result, err := GetResponseAs[vcf.Task](response)

	if result != nil {
		t.Fatal("received an unexpected response", result)
	}

	if err == nil {
		t.Fatal("error is nil")
	}

	if *err.Message != message {
		t.Fatal("response does not contain correct payload")
	}
}
