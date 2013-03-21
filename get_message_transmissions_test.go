package postage_app

import (
	"testing"
)

func TestGetMessageTransmissionsSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	response, _ := cl.GetMessages()

	for k, _ := range response.Data {
		response2, _ := cl.GetMessageTransmissions(k)
		if response2.Response.Status != "ok" {
			t.Log(response2.Response.Status)
			t.Fail()
		}

		if response2.Data == nil {
			t.Log("Data is nil")
			t.Fail()
		}
		break
	}
}

func TestGetMessageTransmissionsNotFound(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	_, err := cl.GetMessageTransmissions("strange UID")
	if err == nil {
		t.Log("Error is nil", err)
		t.Fail()
	}

	if err.Error() != "not_found" {
		t.Log("Error is not 'not_found'", err)
		t.Fail()
	}
}
