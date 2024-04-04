package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ResponseWhatsapp(from, name, typeOfMessage string, ctx context.Context) error {
	var message string

	switch typeOfMessage {
	case "invitation":
		message = fmt.Sprintf(invitationStr, name)
	case "Si":
		message = invitationSi
	case "No":
		message = invitationNo
	case "Tal vez":
		message = invitationMaybe
	default:
		message = defaultMessage
	}

	whatsappToken := os.Getenv("WHATSAPP_TOKEN")
	messageData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                from,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": false,
			"body":        message,
		},
	}

	url := "https://graph.facebook.com/v17.0/117929121407949/messages"

	requestBody, err := json.Marshal(messageData)
	if err != nil {
		err = errors.New("whatsapp/message error decoding messageData in whatsapp")
		return err
	}

	authorization := fmt.Sprintf("Bearer %s", whatsappToken)
	headers := map[string]string{
		"Authorization": authorization,
		"Content-Type":  "application/json",
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))

	if err != nil {
		finalError := fmt.Sprintf(
			"whatsapp/message error creating http request in whatsapp: %s",
			err,
		)

		return errors.New(finalError)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		finalError := fmt.Sprintf("whatsapp/message error doing the http request: %s", err)
		return errors.New(finalError)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		responseBody, err := io.ReadAll(response.Body)
		finalError := ""

		if err != nil {
			finalError = fmt.Sprintf("whatsapp/message http error and reading the response: %s", err)
		} else {
			finalError = fmt.Sprintf("whatsapp/message http error response: %s", responseBody)
		}

		return errors.New(finalError)
	}
}

func SendImage(from, imageId string, ctx context.Context) error {
	whatsappToken := os.Getenv("WHATSAPP_TOKEN")
	messageData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                from,
		"type":              "image",
		"image": map[string]string{
			"id": imageId,
		},
	}
	requestBody, err := json.Marshal(messageData)
	url := "https://graph.facebook.com/v17.0/117929121407949/messages"

	if err != nil {
		err = errors.New("whatsapp/image error decoding messageData in whatsapp")
		return err
	}

	authorization := fmt.Sprintf("Bearer %s", whatsappToken)
	headers := map[string]string{
		"Authorization": authorization,
		"Content-Type":  "application/json",
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))

	if err != nil {
		finalError := fmt.Sprintf("whatsapp/image error creating http request in whatsapp: %s", err)
		return errors.New(finalError)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		finalError := fmt.Sprintf("whatsapp/image error doing the http request: %s", err)
		return errors.New(finalError)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	responseBody, err := io.ReadAll(response.Body)
	finalError := ""

	if err != nil {
		finalError = fmt.Sprintf("whatsapp/image http error and reading the response: %s", err)
	} else {
		finalError = fmt.Sprintf("whatsapp/image http error response: %s", responseBody)
	}

	return errors.New(finalError)
}
