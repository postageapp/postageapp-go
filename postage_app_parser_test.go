package postage_app

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func InitMessage() (*Client, *Message) {
	cl := new(Client)
	cl.ApiKey = ApiKey
	message := new(Message)
	message.Uid = "6e36017c-b662-441b-92cb-3acba3d556f4"
	return cl, message
}

func unmarshal(str string) map[string]interface{} {
	var f interface{}
	bs := []byte(str)
	json.Unmarshal(bs, &f)
	// fmt.Println(str)
	// fmt.Println(f)
	return f.(map[string]interface{})
}

func TestMessageParseIncludesApiKey(t *testing.T) {
	cl, message := InitMessage()
	b, _ := cl.MarshalMessage(message)

	if !strings.Contains(string(b), fmt.Sprintf(`"api_key":"%s"`, ApiKey)) {
		t.Log(string(b))
		t.Fail()
	}
}

func TestMessageParseIncludesText(t *testing.T) {
	cl, message := InitMessage()
	message.Text = "This is my text content"
	b, _ := cl.MarshalMessage(message)

	if !strings.Contains(string(b), `"arguments":{"content":{"text/plain":"This is my text content"}}`) {
		t.Log(string(b))
		t.Fail()
	}
}

func TestMessageParseIncludesHtml(t *testing.T) {
	cl, message := InitMessage()
	message.Html = "<h1>my content</h1>"
	b, _ := cl.MarshalMessage(message)

	if !strings.Contains(string(b), `"arguments":{"content":{"text/html":"\u003ch1\u003emy content\u003c/h1\u003e"}`) {
		t.Log(string(b))
		t.Fail()
	}
}

func TestMessageParseIncludesRecipients(t *testing.T) {
	cl, message := InitMessage()
	recipient := new(Recipient)
	recipient.Email = "Alan Smithee <alan.smithee@gmail.com>"
	recipient.Variables = make(map[string]string)
	recipient.Variables["first_name"] = "Alan"
	recipient.Variables["last_name"] = "Smithee"
	recipient.Variables["order_id"] = "555"

	recipient2 := new(Recipient)
	recipient2.Email = "Rick James <rick.james@gmail.com>"
	recipient2.Variables = make(map[string]string)
	recipient2.Variables["first_name"] = "Rick"
	recipient2.Variables["last_name"] = "James"
	recipient2.Variables["order_id"] = "556"
	message.Recipients = append(message.Recipients, recipient, recipient2)
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"recipients":{"Alan Smithee \u003calan.smithee@gmail.com\u003e":{"first_name":"Alan","last_name":"Smithee","order_id":"555"},"Rick James \u003crick.james@gmail.com\u003e":{"first_name":"Rick","last_name":"James","order_id":"556"}}}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesRecipientOverride(t *testing.T) {
	cl, message := InitMessage()
	message.RecipientOverride = RecipientOverride
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), fmt.Sprintf(`"arguments":{"recipient_override":"%s"}`, RecipientOverride)) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesTemplate(t *testing.T) {
	cl, message := InitMessage()
	message.Template = "some-template-slug"
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"template":"some-template-slug"}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesVariables(t *testing.T) {
	cl, message := InitMessage()
	message.Variables = make(map[string]string)
	message.Variables["movie"] = "Pee Wee's Big Adventure"
	message.Variables["actor"] = "Meryl Streep"
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"variables":{"actor":"Meryl Streep","movie":"Pee Wee's Big Adventure"}}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesSubjectHeader(t *testing.T) {
	cl, message := InitMessage()
	message.Subject = "my content"
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"headers":{"subject":"my content"}}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesReplyToHeader(t *testing.T) {
	cl, message := InitMessage()
	message.ReplyTo = "test@null.postageapp.com"
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"headers":{"reply-to":"test@null.postageapp.com"}}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesHeaders(t *testing.T) {
	cl, message := InitMessage()
	message.Headers = make(map[string]string)
	message.Headers["Subject"] = "Hello friend!"
	message.Headers["X-Accept-Language"] = "en-us, en"
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"headers":{"Subject":"Hello friend!","X-Accept-Language":"en-us, en"}}`) {
		t.Log(string(b))
		t.Fail()
	}
}
func TestMessageParseIncludesAttachments(t *testing.T) {
	cl, message := InitMessage()
	attachment := new(Attachment)
	attachment.FileName = "readme.txt"
	attachment.ContentType = "text/plain"
	attachment.ContentBytes = []byte("file contents!\n\n")
	message.Attachments = append(message.Attachments, attachment)
	b, _ := cl.MarshalMessage(message)
	if !strings.Contains(string(b), `"arguments":{"attachments":{"readme.txt":{"content":"ZmlsZSBjb250ZW50cyEKCg==","content_type":"text/plain"}}}`) {
		t.Log(string(b))
		t.Fail()
	}
}

func TestMessageTransmissionsParse(t *testing.T) {
	cl, _ := InitMessage()
	js := unmarshal(`{"message":{"id":34968902},"transmissions":{"test@null.postageapp.com":{"status":"completed","created_at":"2013-03-21 17:13:21","failed_at":null,"opened_at":null,"clicked_at":null,"result_code":"SMTP_250","error_message":"2.0.0 OK 1363886006 jt2si6390208obb.44 - gsmtp"}}}`)
	mT := cl.ParseMessageTransmissions(js)
	if mT.Id != 34968902 {
		t.Log(mT.Id)
		t.Fail()
	}

	for k, v := range mT.Transmissions {
		if k != "test@null.postageapp.com" {
			t.Log(k)
			t.Fail()
		}

		if v.Status != "completed" {
			t.Log(mT.Id)
			t.Fail()
		}

		ti, _ := time.Parse(time.RFC3339, "2013-03-21 17:13:21")
		if v.CreatedAt != ti {
			t.Log(v.CreatedAt)
			t.Fail()
		}

		if v.ResultCode != "SMTP_250" {
			t.Log(v.ResultCode)
			t.Fail()
		}

		if v.ResultMessage != "2.0.0 OK 1363886006 jt2si6390208obb.44 - gsmtp" {
			t.Log(v.ResultMessage)
			t.Fail()
		}
		break
	}
}

func TestMetricsParse(t *testing.T) {
	cl, _ := InitMessage()
	js := unmarshal(`{"metrics":{"hour":{"delivered":{"current_percent":98,"previous_percent":100,"diff_percent":-1.4,"current_value":69,"previous_value":91},"opened":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"clicked":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"failed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"rejected":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"spammed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"created":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":70,"previous_value":91},"queued":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":1,"previous_value":0}},"date":{"delivered":{"current_percent":99.37888198757764,"previous_percent":100.38167938931298,"diff_percent":-1.0027974017353358,"current_value":160,"previous_value":263},"opened":{"current_percent":0.0,"previous_percent":1.1450381679389312,"diff_percent":-1.1450381679389312,"current_value":0,"previous_value":3},"clicked":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"failed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"rejected":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"spammed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"created":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":161,"previous_value":262},"queued":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":1,"previous_value":0}},"week":{"delivered":{"current_percent":100.0,"previous_percent":100.0,"diff_percent":0.0,"current_value":477,"previous_value":99},"opened":{"current_percent":1.0482180293501049,"previous_percent":1.0101010101010102,"diff_percent":0.03811701924909472,"current_value":5,"previous_value":1},"clicked":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"failed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"rejected":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"spammed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"created":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":477,"previous_value":99},"queued":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0}},"month":{"delivered":{"current_percent":100.0,"previous_percent":18.970578849958518,"diff_percent":81.02942115004149,"current_value":577,"previous_value":5945},"opened":{"current_percent":1.2131715771230502,"previous_percent":0.02871912693854107,"diff_percent":1.184452450184509,"current_value":7,"previous_value":9},"clicked":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"failed":{"current_percent":0.0,"previous_percent":0.003191014104282341,"diff_percent":-0.003191014104282341,"current_value":0,"previous_value":1},"rejected":{"current_percent":0.0,"previous_percent":0.009573042312847023,"diff_percent":-0.009573042312847023,"current_value":0,"previous_value":3},"spammed":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":0},"created":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":577,"previous_value":31338},"queued":{"current_percent":0.0,"previous_percent":0.0,"diff_percent":0.0,"current_value":0,"previous_value":25389}}}}`)
	mT := cl.ParseMetrics(js)
	if mT.Date == nil {
		t.Log(mT.Date)
		t.Fail()
	}
	if mT.Date == nil {
		t.Log(mT.Date)
		t.Fail()
	}
	if mT.Month == nil {
		t.Log(mT.Month)
		t.Fail()
	}
	if mT.Week == nil {
		t.Log(mT.Week)
		t.Fail()
	}

	if mT.Date.Clicked == nil {
		t.Log(mT.Date.Clicked)
		t.Fail()
	}
	if mT.Date.Created == nil {
		t.Log(mT.Date.Created)
		t.Fail()
	}
	if mT.Date.Delivered == nil {
		t.Log(mT.Date.Delivered)
		t.Fail()
	}
	if mT.Date.Failed == nil {
		t.Log(mT.Date.Failed)
		t.Fail()
	}
	if mT.Date.Opened == nil {
		t.Log(mT.Date.Opened)
		t.Fail()
	}
	if mT.Date.Queued == nil {
		t.Log(mT.Date.Queued)
		t.Fail()
	}
	if mT.Date.Rejected == nil {
		t.Log(mT.Date.Rejected)
		t.Fail()
	}
	if mT.Date.Spammed == nil {
		t.Log(mT.Date.Spammed)
		t.Fail()
	}

	delivered := mT.Date.Delivered
	if delivered.CurrentPercent != 99.37888198757764 {
		t.Log(delivered.CurrentPercent)
		t.Fail()
	}

	if delivered.PreviousPercent != 100.38167938931298 {
		t.Log(delivered.PreviousPercent)
		t.Fail()
	}

	if delivered.DiffPercent != -1.0027974017353358 {
		t.Log(delivered.DiffPercent)
		t.Fail()
	}

	if delivered.CurrentValue != 160 {
		t.Log(delivered.CurrentValue)
		t.Fail()
	}
	if delivered.PreviousValue != 263 {
		t.Log(delivered.PreviousValue)
		t.Fail()
	}
}

func TestPostageResponseParseUid(t *testing.T) {
	cl, _ := InitMessage()
	js := unmarshal(`{"status":"ok","uid":"89067504-0789-47c6-8817-93d3f4e6f8f7"}`)
	resp := cl.ParseResponse(js)

	if resp.Status != "ok" {
		t.Log(resp.Uid)
		t.Fail()
	}

	if resp.Uid != "89067504-0789-47c6-8817-93d3f4e6f8f7" {
		t.Log(resp.Uid)
		t.Fail()
	}

}

func TestPostageResponseParseErrorMessage(t *testing.T) {
	cl, _ := InitMessage()
	js := unmarshal(`{"status":"ok","uid":"89067504-0789-47c6-8817-93d3f4e6f8f7", "message": "Something went wrong!"}`)
	resp := cl.ParseResponse(js)

	if resp.Message != "Something went wrong!" {
		t.Log(resp.Message)
		t.Fail()
	}
}
func TestProjectInfoConverter(t *testing.T) {
	cl, _ := InitMessage()
	js := unmarshal(`{"project":{"name":"Test","url":"https://api.postageapp.com/projects/1212","transmissions":{"today":160,"this_month":577,"overall":16724},"users":{"postage+tester@twg.ca":"Test User","cloudy@mailinator.com":"☁"}}}`)
	proj := cl.ParseProjectInfo(js)

	if proj.Name != "Test" {
		t.Log(proj.Name)
		t.Fail()
	}

	if proj.Url != "https://api.postageapp.com/projects/1212" {
		t.Log(proj.Url)
		t.Fail()
	}

	if proj.Transmissions == nil {
		t.Fail()
	}

	if proj.Transmissions.TodayCount != 160 {
		t.Log(proj.Transmissions.TodayCount)
		t.Fail()
	}

	if proj.Transmissions.ThisMonthCount != 577 {
		t.Log(proj.Transmissions.ThisMonthCount)
		t.Fail()
	}

	if proj.Transmissions.OverallCount != 16724 {
		t.Log(proj.Transmissions.OverallCount)
		t.Fail()
	}

	if proj.Users == nil {
		t.Fail()
	}
	for k, v := range proj.Users {
		if k != "postage+tester@twg.ca" {
			t.Log(k)
			t.Fail()
		}

		if v != "Test User" {
			t.Log(k)
			t.Fail()
		}
		break
	}
}
