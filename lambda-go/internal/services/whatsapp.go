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
Querido  %s.

Con alegría en nuestro corazón queremos informarles que hemos decidido unir nuestras vidas en matrimonio 🕊 Por tal motivo les hacemos una cordial invitación a nuestra boda ✉️ La boda se llevará a cabo el día *16 de agosto del 2024*, en la ciudad de *Olomouc, República Checa* 🇨🇿. Estaremos muy contentos si podrán acompañarnos en este gran día. En la parte baja del texto encontrarán el enlace a la pagina web para más información sobre el evento y tips cómo llegar, dónde hospedarse y qué visitar. Para cualquier duda, pueden escribirnos por acá, con gusto les responderemos. En enero de 2024 vamos a enviar un nuevo mensaje para recibir su confirmación de asistencia. Les mandamos un abrazo sincero. Kami y Julián.

Todo lo que se hace por amor, se hace más allá del bien y de mal.
Nietzsche

Para más información por favor visita nuestra pagina web 👇🏾
https://www.theknot.com/kamilaandjulian
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

	url := "https://graph.facebook.com/v17.0/126502743873273/messages"

	requestBody, err := json.Marshal(messageData)
	if err != nil {
		err = errors.New("error decoding messageData in whatsapp")
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
		finalError := fmt.Sprintf("error creating http request in whatsapp: %s", err)
		return errors.New(finalError)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		finalError := fmt.Sprintf("error doing the http request: %s", err)
		return errors.New(finalError)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		responseBody, err := io.ReadAll(response.Body)
		finalError := ""
		if err != nil {
			finalError = fmt.Sprintf("whatsapp http error and reading the response: %s", err)
		} else {
			finalError = fmt.Sprintf("whatsapp http error response: %s", responseBody)
		}
		return errors.New(finalError)
	}
}
