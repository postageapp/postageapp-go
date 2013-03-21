package postage_app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	Url = "https://api.postageapp.com/v.1.0/"
)

type Client struct {
	ApiKey  string
	BaseUrl string
}

type Attachment struct {
	FileName     string
	ContentType  string
	ContentBytes []byte
}

type Recipient struct {
	Email     string
	Variables map[string]string
}

type Message struct {
	Uid               string
	Template          string
	Attachments       []*Attachment
	Recipients        []*Recipient
	Variables         map[string]string
	Headers           map[string]string
	RecipientOverride string
	Subject           string
	From              string
	ReplyTo           string
	Text              string
	Html              string
}

type MessageInfo struct {
	ProjectId                   int
	Template                    string
	TotalTransmissionsCount     int
	FailedTransmissionsCount    int
	CompletedTransmissionsCount int
	CreatedAt                   time.Time
	WillPurgeAt                 time.Time
}

type MessageReceipt struct {
	Id  int
	Url string
}

type MessageTransmission struct {
	Status        string
	ResultCode    string
	ResultMessage string
	CreatedAt     time.Time
	FailedAt      time.Time
	OpenedAt      time.Time
}

type MessageTransmissions struct {
	Id            int
	Transmissions map[string]*MessageTransmission
}

type TransmissionsStatistic struct {
	TodayCount     int
	ThisMonthCount int
	OverallCount   int
}

type ProjectInfo struct {
	Name          string
	Url           string
	Transmissions *TransmissionsStatistic
	Users         map[string]string
}

type AccountInfo struct {
	Name          string
	Url           string
	Transmissions *TransmissionsStatistic
	Users         map[string]string
}

type MetricStatistic struct {
	CurrentPercent  float64
	PreviousPercent float64
	DiffPercent     float64
	CurrentValue    int
	PreviousValue   int
}

type Metric struct {
	Delivered *MetricStatistic
	Opened    *MetricStatistic
	Failed    *MetricStatistic
	Rejected  *MetricStatistic
	Created   *MetricStatistic
	Queued    *MetricStatistic
	Clicked   *MetricStatistic
	Spammed   *MetricStatistic
}

type Metrics struct {
	Hour  *Metric
	Date  *Metric
	Week  *Metric
	Month *Metric
}

type Response struct {
	Status  string
	Uid     string
	Message string
}

type MessagesResponse struct {
	Response *Response
	Data     map[string]*MessageInfo
}

type ProjectResponse struct {
	Response *Response
	Data     *ProjectInfo
}

type AccountResponse struct {
	Response *Response
	Data     *AccountInfo
}

type MetricsResponse struct {
	Response *Response
	Data     *Metrics
}

type MessageResponse struct {
	Response *Response
	Data     *MessageReceipt
}

type MessageTransmissionsResponse struct {
	Response *Response
	Data     *MessageTransmissions
}

type PostageError struct {
	Message    string
	InnerError error
}

type PostageResponseError PostageError

type ResponseParseError PostageError

func (e *PostageError) Error() string {
	return e.Message
}

func (e *ResponseParseError) Error() string {
	return e.Message
}

func (e *PostageResponseError) Error() string {
	return e.Message
}

func (client *Client) post(path string, params string) (map[string]interface{}, error) {
	b := bytes.NewBufferString(params)
	return client.postBuffer(path, b)
}

func (client *Client) postBuffer(path string, b *bytes.Buffer) (map[string]interface{}, error) {
	if client.BaseUrl == "" {
		client.BaseUrl = Url
	}
	url := Url + path
	response, err := http.Post(url, "application/json", b)
	if err != nil {
		return nil, &PostageResponseError{err.Error(), err}
	}

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f interface{}
	parseError := json.Unmarshal(bs, &f)
	if parseError != nil {
		return nil, &PostageResponseError{err.Error(), err}
	}

	return f.(map[string]interface{}), nil
}

func (client *Client) SendMessage(message *Message) (*MessageResponse, error) {
	bts, _ := client.MarshalMessage(message)
	b := bytes.NewBuffer(bts)

	m, err := client.postBuffer("send_message.json", b)
	if err != nil {
		return nil, err
	}
	messageResponse := new(MessageResponse)
	messageResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if messageResponse.Response.Status == "ok" {
		data := m["data"].(map[string]interface{})
		message := data["message"].(map[string]interface{})
		messageReceipt := new(MessageReceipt)
		messageReceipt.Id = int(message["id"].(float64))
		messageReceipt.Url = message["url"].(string)
		messageResponse.Data = messageReceipt
	} else {
		return nil, &ResponseParseError{messageResponse.Response.Status, nil}
	}

	return messageResponse, nil
}

func (client *Client) GetMessages() (*MessagesResponse, error) {
	m, err := client.post("get_messages.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey))
	if err != nil {
		return nil, err
	}

	messagesResponse := new(MessagesResponse)

	messagesResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if messagesResponse.Response.Status == "ok" {
		messagesResponse.Data = client.ParseMessages(m["data"].(map[string]interface{}))
	} else {
		return nil, &ResponseParseError{messagesResponse.Response.Status, nil}
	}

	return messagesResponse, nil
}

func (client *Client) GetProjectInfo() (*ProjectResponse, error) {
	m, err := client.post("get_project_info.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey))
	if err != nil {
		return nil, err
	}

	projectResponse := new(ProjectResponse)

	projectResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if projectResponse.Response.Status == "ok" {
		projectResponse.Data = client.ParseProjectInfo(m["data"].(map[string]interface{}))
	} else {
		return nil, &ResponseParseError{projectResponse.Response.Status, nil}
	}

	return projectResponse, nil
}

func (client *Client) GetAccountInfo() (*AccountResponse, error) {
	m, err := client.post("get_account_info.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey))
	if err != nil {
		return nil, err
	}

	accountResponse := new(AccountResponse)

	accountResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))
	if accountResponse.Response.Status == "ok" {
		accountResponse.Data = client.ParseAccountInfo(m["data"].(map[string]interface{}))
	} else {
		return nil, &ResponseParseError{accountResponse.Response.Status, nil}
	}

	return accountResponse, nil
}

func (client *Client) GetMetrics() (*MetricsResponse, error) {
	m, err := client.post("get_metrics.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey))
	if err != nil {
		return nil, err
	}

	metricsResponse := new(MetricsResponse)

	metricsResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if metricsResponse.Response.Status == "ok" {
		metricsResponse.Data = client.ParseMetrics(m["data"].(map[string]interface{}))
	} else {
		return nil, &ResponseParseError{metricsResponse.Response.Status, nil}
	}

	return metricsResponse, nil
}

func (client *Client) GetMessageReceipt(uid string) (*MessageResponse, error) {
	m, err := client.post("get_message_receipt.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey, uid))
	if err != nil {
		return nil, err
	}

	messageReceiptResponse := new(MessageResponse)
	messageReceiptResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if messageReceiptResponse.Response.Status == "ok" {
		messageReceiptResponse.Data = client.ParseMessageReceipt(m["data"].(map[string]interface{}))
	} else {
		return nil, &ResponseParseError{messageReceiptResponse.Response.Status, nil}
	}
	return messageReceiptResponse, nil
}

func (client *Client) GetMessageTransmissions(uid string) (*MessageTransmissionsResponse, error) {
	m, err := client.post("get_message_transmissions.json", fmt.Sprintf(`{"api_key":"%s", "uid":"%s"}`, client.ApiKey, uid))
	if err != nil {
		return nil, err
	}

	messageTransmissionsResponse := new(MessageTransmissionsResponse)
	messageTransmissionsResponse.Response = client.ParseResponse(m["response"].(map[string]interface{}))

	if messageTransmissionsResponse.Response.Status == "ok" {
		messageTransmissionsResponse.Data = client.ParseMessageTransmissions(m["data"].(map[string]interface{}))

	} else {
		return nil, &ResponseParseError{messageTransmissionsResponse.Response.Status, nil}
	}
	return messageTransmissionsResponse, nil
}
