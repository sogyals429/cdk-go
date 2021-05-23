package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {

	eventJson, _ := json.MarshalIndent(event, "", " ")
	log.Printf("EVENT: %s", eventJson)

	resp, err := http.Get("http://google.com/")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	lambda.Start(HandleRequest)
}
