package function

import (
	"encoding/json"
	"github.com/robbiet480/go.sns"
	"log"
)

// Handle a serverless request
func Handle(req []byte) string {
	var message sns.Payload
	err := json.Unmarshal(req, &message)
	if err != nil {
		log.Fatalf("There was an error:", err)
	}

	verifyErr := message.VerifyPayload()
	if verifyErr != nil {
		log.Fatal(verifyErr)
	}

	var returnValue string

	switch message.Type {
	case "SubscriptionConfirmation":
		message.Subscribe()
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
