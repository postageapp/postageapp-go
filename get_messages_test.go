package postage_app

import (
	"testing"
)

func TestGetMessagesSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	response, _ := cl.GetMessages()
	if response.Response.Status != "ok" {
		t.Log(response.Response.Status)
		t.Fail()
	}
	if response.Data == nil {
		t.Log("Data is nil")
		t.Fail()
	}
}

func TestGetMessagesUnauthorized(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = "abc123ThisIsNotValid"
	_, err := cl.GetMessages()
	if err == nil {
		t.Log("Error is nil")
		t.Fail()
	}

	if err.Error() != "unauthorized" {
		t.Fail()
	}
}
