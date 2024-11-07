package api_client

import (
	"encoding/json"
	"testing"

	"github.com/vmware/vcf-sdk-go/vcf"
)

func TestGetResponseAs_pos(t *testing.T) {
	taskId := "id"
	model := vcf.Task{Id: &taskId}
	body, _ := json.Marshal(model)

	res, err := GetResponseAs[vcf.Task](body, 200)

	if err != nil {
		t.Fatal("received an unexpected error", err)
	}

	if res == nil {
		t.Fatal("response is nil")
	}

	if *res.Id != taskId {
		t.Fatal("response does not contain correct payload")
	}
}

func TestGetResponseAs_neg(t *testing.T) {
	message := "message"
	model := vcf.Error{Message: &message}
	body, _ := json.Marshal(model)

	res, err := GetResponseAs[vcf.Task](body, 400)

	if res != nil {
		t.Fatal("received an unexpected response", res)
	}

	if err == nil {
		t.Fatal("error is nil")
	}

	if *err.Message != message {
		t.Fatal("response does not contain correct payload")
	}
}
