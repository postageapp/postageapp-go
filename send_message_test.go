package postage_app

import (
	"fmt"
	"github.com/golibs/uuid"
	"testing"
)

const (
	ApiKey            = "YOU API KEY"
	RecipientOverride = ""
)

func TestSendMessageWithAttachment(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()

	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"
	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "This is my text content ☃☃"
	attachment := new(Attachment)
	attachment.FileName = "readme.txt"
	attachment.ContentType = "text/plain"
	attachment.ContentBytes = []byte("file contents!\n\n")
	message.Attachments = append(message.Attachments, attachment)
	response, _ := cl.SendMessage(message)

	if response.Response.Status != "ok" {
		fmt.Println(response.Response.Status)
		t.Fail()
	}

	if &response.Response.Uid == nil {
		fmt.Println("Uid is nil")
		t.Fail()
	}

	if &response.Data.Id == nil {
		t.Log("Data.Id is nil")
		t.Fail()
	}
}

func TestSendMessageInvalidDomain(t *testing.T) {
	cl := new(Client)
	cl.BaseUrl = "http://0.0.0.0/"
	cl.ApiKey = ApiKey

	message := new(Message)
	message.Uid = "6e36017c-b662-441b-92cb-3acba3d556f4"

	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"
	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "This is my text content"
}

func TestSendMessageBadRequest(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey

	message := new(Message)
	_, err := cl.SendMessage(message)

	if err == nil {
		t.Log("Error is nil", err)
		t.Fail()
	}

	if err.Error() != "bad_request" {
		t.Log("Error is not 'bad_request', is", err)
		t.Fail()
	}

}
func TestSendMessageSuccess(t *testing.T) {
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
	if response.Response.Status != "ok" {
		fmt.Println(response.Response.Status)
		t.Fail()
	}

	if &response.Response.Uid == nil {
		fmt.Println("Uid is nil")
		t.Fail()
	}

	if &response.Data.Id == nil {
		t.Log("Data.Id is nil")
		t.Fail()
	}
}

func TestSendUnicodeMessageSuccess(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()

	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"
	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "There are my snowmans ☃☃☃☃☃☃"
	response, _ := cl.SendMessage(message)
	if response.Response.Status != "ok" {
		fmt.Println("Expected 'ok' but was :", response.Response.Status)
		t.Fail()
	}

	if &response.Response.Uid == nil {
		fmt.Println("Uid is nil")
		t.Fail()
	}

	if &response.Data.Id == nil {
		t.Log("Data.Id is nil")
		t.Fail()
	}
}

func TestSendMessagePreconditionFailed(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()
	message.Template = "some-unknown-template-xxxxxxxxxxxxxx"

	_, err := cl.SendMessage(message)
	if err.Error() != "precondition_failed" {
		fmt.Println("Expected 'precondition_failed' but was :", err)
		t.Fail()
	}
}

func TestSendMessageUidRoundTrip(t *testing.T) {
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
	if response.Response.Status != "ok" {
		fmt.Println(response.Response.Status)
		t.Fail()
	}
	if response.Response.Uid != message.Uid {
		fmt.Println(response.Response.Uid, "!=", message.Uid)
		t.Fail()
	}
}

func TestSendMessageWithHeaders(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()

	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"
	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "This is my text content"
	message.Headers = make(map[string]string)
	message.Headers["From"] = "test.robot@null.postageapp.com"
	message.Headers["Subject"] = "This is a custom subject line"

	response, _ := cl.SendMessage(message)
	if response.Response.Status != "ok" {
		fmt.Println(response.Response.Status)
		t.Fail()
	}
}

func TestSendMessageWithHtmlContent(t *testing.T) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = uuid.Rand().Hex()
	message.Subject = "Html body"
	recipient := new(Recipient)
	recipient.Email = "test@null.postageapp.com"

	message.Recipients = append(message.Recipients, recipient)
	message.RecipientOverride = RecipientOverride
	message.Text = "This email should have some html content."
	message.Html = "<h1>Title</h1><p>This is an <em>html email</em></p></h1>"

	response, _ := cl.SendMessage(message)
	if response.Response.Status != "ok" {
		fmt.Println(response.Response.Status)
		t.Fail()
	}
}
