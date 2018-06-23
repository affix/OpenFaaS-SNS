package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robbiet480/go.sns"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type s3Notification struct {
	Records []struct {
		EventVersion string    `json:"eventVersion"`
		EventSource  string    `json:"eventSource"`
		AwsRegion    string    `json:"awsRegion"`
		EventTime    time.Time `json:"eventTime"`
		EventName    string    `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				ETag      string `json:"eTag"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

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
	var notification s3Notification
	json.Unmarshal([]byte(message), &notification)
	url := fmt.Sprintf("http://s3-%s.amazonaws.com/%s/%s",
		notification.Records[0].AwsRegion,
		notification.Records[0].S3.Bucket.Name,
		notification.Records[0].S3.Object.Key)

	return CallFunction("http://ocr:8080", url)
}

func CallFunction(function string, imageUrl string) string {
	imageUrlBytes := []byte(imageUrl)
	req, err := http.NewRequest("POST", function, bytes.NewBuffer(imageUrlBytes))
	if err != nil {
		return fmt.Sprintf("http.NewRequest() error: %v\n", err)
	}

	c := &http.Client{}
	resp, _ := c.Do(req)
	if err != nil {
		return fmt.Sprintf("http.Do() error: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("ioutil.ReadAll() error: %v\n", err)
	}

	return fmt.Sprintf("Parsed %s :\n%v\n", imageUrl, string(data))
}
