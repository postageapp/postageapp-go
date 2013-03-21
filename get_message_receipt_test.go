package postage_app

import (
	"github.com/golibs/uuid"
	"testing"
)

func TestGetMessageReceiptSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()

	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"
	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "This is my text content"
	response, _ := cl.SendMessage(message)

	response2, _ := cl.GetMessageReceipt(response.Response.Uid)
	if response2.Response.Status != "ok" {
		t.Log(response2.Response.Status)
		t.Fail()
	}

	if &response2.Data.Id == nil {
		t.Log("Data.Id is nil")
		t.Fail()
	}
}

func TestGetMessageReceiptNotFound(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	_, err := cl.GetMessageReceipt("strange UID")
	if err == nil {
		t.Log("Error is nil", err)
		t.Fail()
	}

	if err.Error() != "not_found" {
		t.Log("Error is not 'not_found'", err)
		t.Fail()
	}
}
