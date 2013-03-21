package postage_app

import (
	"testing"
)

func TestGetMetricsSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	response, _ := cl.GetMetrics()
	if response.Response.Status != "ok" {
		t.Log(response.Response.Status)
		t.Fail()
	}
	if response.Data == nil {
		t.Log("Data is nil")
		t.Fail()
	}

	if response.Data.Date == nil {
		t.Log("Data.Date is nil")
		t.Fail()
	}

	if response.Data.Hour == nil {
		t.Log("Data.Hour is nil")
		t.Fail()
	}

	if response.Data.Month == nil {
		t.Log("Data.Month is nil")
		t.Fail()
	}

	if response.Data.Week == nil {
		t.Log("Data.Week is nil")
		t.Fail()
	}

}
