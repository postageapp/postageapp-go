package postage_app

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

func (client *Client) ParseResponse(json map[string]interface{}) *Response {
	response := new(Response)
	response.Status = json["status"].(string)
	if json["uid"] != nil {
		uidS := json["uid"].(string)
		response.Uid = uidS
	}

	if json["message"] != nil {
		response.Message = json["message"].(string)
	}

	return response
}

func (client *Client) MarshalMessage(message *Message) ([]byte, error) {
	hash := map[string]interface{}{
		"api_key":   client.ApiKey,
		"uid":       message.Uid,
		"arguments": map[string]interface{}{},
	}

	arguments := hash["arguments"].(map[string]interface{})

	if message.Html != "" || message.Text != "" {
		content := map[string]interface{}{}
		arguments["content"] = content
		if message.Html != "" {
			content["text/html"] = message.Html
		}

		if message.Text != "" {
			content["text/plain"] = message.Text
		}
	}

	if len(message.Recipients) != 0 {
		recipients := map[string]interface{}{}
		arguments["recipients"] = recipients
		for _, recipient := range message.Recipients {
			variables := map[string]interface{}{}
			recipients[recipient.Email] = variables

			for key, variable := range recipient.Variables {
				variables[key] = variable
			}
		}
	}

	if message.RecipientOverride != "" {
		arguments["recipient_override"] = message.RecipientOverride
	}

	if message.Template != "" {
		arguments["template"] = message.Template
	}

	if len(message.Variables) != 0 {
		variables := map[string]interface{}{}
		arguments["variables"] = variables
		for key, variable := range message.Variables {
			variables[key] = variable
		}
	}

	if message.Subject != "" || message.From != "" || message.ReplyTo != "" || len(message.Headers) != 0 {
		headers := map[string]interface{}{}
		arguments["headers"] = headers
		if message.Subject != "" {
			headers["subject"] = message.Subject
		}

		if message.From != "" {
			headers["from"] = message.From
		}

		if message.ReplyTo != "" {
			headers["reply-to"] = message.ReplyTo
		}

		if len(message.Headers) != 0 {
			for key, header := range message.Headers {
				headers[key] = header
			}
		}
	}

	if len(message.Attachments) != 0 {
		attachments := map[string]interface{}{}
		arguments["attachments"] = attachments
		for _, attachment := range message.Attachments {
			attachmentJ := map[string]interface{}{}
			attachmentJ["content_type"] = attachment.ContentType
			attachmentJ["content"] = base64.StdEncoding.EncodeToString(attachment.ContentBytes)
			attachments[attachment.FileName] = attachmentJ
		}
	}

	return json.Marshal(hash)
}

func (client *Client) ParseMessages(json map[string]interface{}) map[string]*MessageInfo {
	data := make(map[string]*MessageInfo)
	for messageUid, messageInf := range json {
		messageInfo := new(MessageInfo)
		miJ := messageInf.(map[string]interface{})

		if miJ["project_id"] != nil {
			messageInfo.ProjectId = int(miJ["project_id"].(float64))
		}

		if miJ["template"] != nil {
			messageInfo.Template = miJ["template"].(string)
		}

		if miJ["transmissions_total"] != nil {
			messageInfo.TotalTransmissionsCount = int(miJ["transmissions_total"].(float64))
		}

		if miJ["transmissions_failed"] != nil {
			messageInfo.FailedTransmissionsCount = int(miJ["transmissions_failed"].(float64))
		}

		if miJ["transmissions_completed"] != nil {
			messageInfo.CompletedTransmissionsCount = int(miJ["transmissions_completed"].(float64))
		}

		if miJ["created_at"] != nil {
			createdAt, _ := time.Parse(time.RFC3339, miJ["created_at"].(string))
			messageInfo.CreatedAt = createdAt
		}

		if miJ["will_purge_at"] != nil {
			willPurgeAt, _ := time.Parse(time.RFC3339, miJ["will_purge_at"].(string))
			messageInfo.WillPurgeAt = willPurgeAt
		}

		data[messageUid] = messageInfo
	}

	return data
}

func (client *Client) ParseProjectInfo(json map[string]interface{}) *ProjectInfo {
	projectInfo := new(ProjectInfo)
	project := json["project"].(map[string]interface{})
	projectInfo.Name = project["name"].(string)
	projectInfo.Url = project["url"].(string)
	tr := project["transmissions"].(map[string]interface{})

	statistic := new(TransmissionsStatistic)
	statistic.TodayCount = int(tr["today"].(float64))
	statistic.ThisMonthCount = int(tr["this_month"].(float64))
	statistic.OverallCount = int(tr["overall"].(float64))
	projectInfo.Transmissions = statistic

	projectInfo.Users = make(map[string]string)
	users := project["users"].(map[string]interface{})
	for email, name := range users {
		projectInfo.Users[email] = name.(string)
	}
	return projectInfo
}

func (client *Client) ParseAccountInfo(json map[string]interface{}) *AccountInfo {
	accountInfo := new(AccountInfo)

	account := json["account"].(map[string]interface{})
	accountInfo.Name = account["name"].(string)
	accountInfo.Url = account["url"].(string)
	tr := account["transmissions"].(map[string]interface{})

	statistic := new(TransmissionsStatistic)
	statistic.TodayCount = int(tr["today"].(float64))
	statistic.ThisMonthCount = int(tr["this_month"].(float64))
	statistic.OverallCount = int(tr["overall"].(float64))
	accountInfo.Transmissions = statistic

	accountInfo.Users = make(map[string]string)
	users := account["users"].(map[string]interface{})
	for email, name := range users {
		accountInfo.Users[email] = name.(string)
	}
	return accountInfo
}

func (client *Client) ParseMessageReceipt(json map[string]interface{}) *MessageReceipt {
	message := json["message"].(map[string]interface{})
	messageReceipt := new(MessageReceipt)
	messageReceipt.Id = int(message["id"].(float64))
	messageReceipt.Url = message["url"].(string)
	return messageReceipt
}

func (client *Client) ParseMessageTransmissions(json map[string]interface{}) *MessageTransmissions {
	message := json["message"].(map[string]interface{})

	messageTransmissions := new(MessageTransmissions)
	messageTransmissions.Id = int(message["id"].(float64))
	messageTransmissions.Transmissions = make(map[string]*MessageTransmission)
	for email, tr := range json["transmissions"].(map[string]interface{}) {
		transmissionJson := tr.(map[string]interface{})
		messageTransmission := new(MessageTransmission)
		messageTransmission.Status = transmissionJson["status"].(string)
		if transmissionJson["error_message"] != nil {
			messageTransmission.ResultMessage = transmissionJson["error_message"].(string)
		}
		if transmissionJson["result_code"] != nil {
			messageTransmission.ResultCode = transmissionJson["result_code"].(string)
		}

		if transmissionJson["created_at"] != nil {
			createdAt, _ := time.Parse(time.RFC3339, transmissionJson["created_at"].(string))
			messageTransmission.CreatedAt = createdAt
		}

		if transmissionJson["failed_at"] != nil {
			failedAt, _ := time.Parse(time.RFC3339, transmissionJson["failed_at"].(string))
			messageTransmission.FailedAt = failedAt
		}

		if transmissionJson["opened_at"] != nil {
			openedAt, _ := time.Parse(time.RFC3339, transmissionJson["opened_at"].(string))
			messageTransmission.OpenedAt = openedAt
		}

		messageTransmissions.Transmissions[email] = messageTransmission
	}

	return messageTransmissions
}

func (client *Client) ParseMetrics(json map[string]interface{}) *Metrics {
	metrics := new(Metrics)
	metricsJson := json["metrics"].(map[string]interface{})

	for metricKey, metricJson := range metricsJson {
		metric := new(Metric)
		for metricStatisticKey, msj := range metricJson.(map[string]interface{}) {
			metricStatisticJson := msj.(map[string]interface{})
			metricStatistic := new(MetricStatistic)
			metricStatistic.CurrentPercent = metricStatisticJson["current_percent"].(float64)
			metricStatistic.PreviousPercent = metricStatisticJson["previous_percent"].(float64)
			metricStatistic.DiffPercent = metricStatisticJson["diff_percent"].(float64)
			metricStatistic.CurrentValue = int(metricStatisticJson["current_value"].(float64))
			metricStatistic.PreviousValue = int(metricStatisticJson["previous_value"].(float64))

			switch metricStatisticKey {
			case "delivered":
				metric.Delivered = metricStatistic
			case "opened":
				metric.Opened = metricStatistic
			case "failed":
				metric.Failed = metricStatistic
			case "rejected":
				metric.Rejected = metricStatistic
			case "created":
				metric.Created = metricStatistic
			case "queued":
				metric.Queued = metricStatistic
			case "clicked":
				metric.Clicked = metricStatistic
			case "spammed":
				metric.Spammed = metricStatistic
			}
		}

		switch metricKey {
		case "hour":
			metrics.Hour = metric
		case "date":
			metrics.Date = metric
		case "week":
			metrics.Week = metric
		case "month":
			metrics.Month = metric
		}
	}

	return metrics
}
