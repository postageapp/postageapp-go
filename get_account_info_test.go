package postage_app

import (
	"testing"
)

func TestGetAccountInfoSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	response, _ := cl.GetAccountInfo()
	if response.Response.Status != "ok" {
		t.Log(response.Response.Status)
		t.Fail()
	}
	if response.Data == nil {
		t.Log("Data is nil")
		t.Fail()
	}

	if response.Data.Url == "" {
		t.Log("Data.Url is nil")
		t.Fail()
	}
}
