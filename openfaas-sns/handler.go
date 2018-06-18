package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SnsMessage struct {
	Type             string    `json:"Type"`
	MessageID        string    `json:"MessageId"`
	Token            string    `json:"Token"`
	TopicArn         string    `json:"TopicArn"`
	Subject          string    `json:"Subject"`
	Message          string    `json:"Message"`
	SubscribeURL     string    `json:"SubscribeURL"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
}

// Handle a serverless request
func Handle(req []byte) string {
	var message SnsMessage
	err := json.Unmarshal(req, &message)
	if err != nil {
		log.Fatalf("There was an error:", err)
	}

	var returnValue string

	switch message.Type {
	case "SubscriptionConfirmation":
		HandleSubscribe(message.SubscribeURL)
		break
	case "Notification":
		returnValue = HandleNotification(message.Message)
		break
	default:
		break
	}

	return returnValue
}

func HandleNotification(message string) string {
	return message
}

func HandleSubscribe(SubscribeURL string) string {
	resp, err := http.Get(SubscribeURL)
	if err != nil {
		log.Fatalf("There was an error:", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
