package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var invitationStr = `
Querido/a  %s.

Les hacemos una cordial invitación a nuestra boda 🕊 La boda se llevará a cabo el día *16 de agosto del 2024*, en la ciudad de *Olomouc, República Checa* 🇨🇿 y será una boda muy pequeña. Estaremos muy contentos si podrán acompañarnos en este gran día. En la parte baja del texto encontrarán el enlace a la pagina web con más información sobre el evento. En enero de 2024 vamos a enviar un nuevo mensaje para saber si podrán acompañarnos. Les mandamos un abrazo sincero. Kami y Julián. 

Nota: no podemos ver sus mensajes en esta cuenta o responderlos 💻

Todo lo que se hace por amor, se hace más allá del bien y del mal.
Nietzsche

Para más información por favor visita nuestra pagina web 👇🏾
https://kamilaandree.wixsite.com/kamiaandjulian
`

func SendInvitation(from string, name string) error {
	inivitation := fmt.Sprintf(invitationStr, name)
	whatsappToken := os.Getenv("WHATSAPP_TOKEN")
	messageData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                from,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": false,
			"body":        inivitation,
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

func SendImage(from string, imageId string) error {
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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
