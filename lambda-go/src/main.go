package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"whatsapp-webhook-golang/src/domain/dto"
	"whatsapp-webhook-golang/src/domain/entity"
	"whatsapp-webhook-golang/src/services"
	"whatsapp-webhook-golang/src/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var wg = sync.WaitGroup{}

	ctx := context.Background()
	badRequestStr := strconv.Itoa(http.StatusBadRequest)
	okRequestStr := strconv.Itoa(http.StatusOK)
	token := "supersecreto"

	if request.HTTPMethod == "GET" &&
		request.QueryStringParameters["hub.mode"] == "subscribe" &&
		request.QueryStringParameters["hub.verify_token"] == token {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       request.QueryStringParameters["hub.challenge"],
		}

		return response, nil
	}

	if request.HTTPMethod == "POST" {
		var whatsappEvent dto.WhatsAppEvent

		logger := map[string]string{
			"user_agent": request.Headers["User-Agent"],
			"ip":         request.RequestContext.Identity.SourceIP,
			"request_id": request.RequestContext.RequestID,
		}
		err := json.Unmarshal([]byte(request.Body), &whatsappEvent)

		if err != nil {
			err := "request body could not be serialized as whatsapp event"
			return utils.HandleError(logger, http.StatusBadRequest, err)
		}

		if !utils.SliceHasElements(whatsappEvent.Entry) {
			return utils.HandleError(logger, http.StatusBadRequest, "no entry found in the request")
		}

		entry := whatsappEvent.Entry[0]

		if !utils.SliceHasElements(entry.Changes) {
			return utils.HandleError(logger, http.StatusBadRequest, "no changes found in the request")
		}

		change := entry.Changes[0]
		field := change.Field
		logger["field"] = field

		if change.Field != "messages" {
			return utils.HandleError(logger, http.StatusBadRequest, "field not supported")
		}

		value := change.Value

		if !utils.SliceHasElements(value.Contacts) && !utils.SliceHasElements(value.Statuses) {
			return utils.HandleError(logger, http.StatusBadRequest, "No contacts and no Statuses in the request")
		} else if !utils.SliceHasElements(value.Contacts) && utils.SliceHasElements(value.Statuses) {
			status := value.Statuses[0]
			logger["status_code"] = strconv.Itoa(http.StatusOK)
			logger["event_name"] = status.Status
			logger["from"] = status.RecipientID
			logger["event_id"] = status.ID
			logger["list_name"] = entity.ListOfPhonesNames[status.RecipientID]
			services.Logger(logger)

			response := events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       "Successful event received",
			}

			return response, nil
		}

		userName := value.Contacts[0].Profile.Name
		logger["user"] = userName

		if !utils.SliceHasElements(value.Messages) {
			return utils.HandleError(logger, http.StatusBadRequest, "no Messages found in the request")
		}

		userCellPhone := value.Messages[0].From
		commonName := entity.ListOfPhonesNames[userCellPhone]
		logger["from"] = userCellPhone
		logger["list_name"] = commonName
		logger["message_id"] = value.Messages[0].ID

		if field != "messages" {
			return utils.HandleError(logger, http.StatusBadRequest, fmt.Sprintf("field not supported: %s", field))
		}

		message := value.Messages[0].Text.Body
		typeMessage := value.Messages[0].Type
		logger["event_name"] = "received message"
		logger["message_received"] = message

		if typeMessage == "text" && strings.ToLower(message) == "si" {
			familyName := entity.ListOfGuestsFamilyNames[userCellPhone]
			err := services.ResponseWhatsapp(userCellPhone, familyName, "invitation", ctx)

			if err != nil {
				err := fmt.Sprintf(
					"error sending whatsapp: %s",
					err,
				)

				return utils.HandleError(logger, http.StatusBadRequest, err)
			}

			// imageId := entity.ListOfGuestsImage[userCellPhone]
			// err = services.SendImage(userCellPhone, imageId, ctx)

			// if err != nil {
			// 	err := fmt.Sprintf(
			// 		"error sending whatsapp image: %s",
			// 		err,
			// 	)

			// 	return utils.HandleError(logger, http.StatusBadRequest, err)
			// }

			logger["event_name"] = "send whatsapp image and text"
		}

		if typeMessage == "button" {
			logger["event_name"] = "fast button"
			message = value.Messages[0].Button.Text
			logger["answer"] = message
			familyName := entity.ListOfGuestsFamilyNames[userCellPhone]
			numberOfRutines := 2
			wg.Add(numberOfRutines)

			go func() {
				err := services.ResponseWhatsapp(userCellPhone, familyName, message, ctx)

				if err != nil {
					err := fmt.Sprintf(
						"error responding whatsapp button: %s",
						err,
					)
					logger["error"] = err
					logger["status_code"] = badRequestStr
					loggerData, _ := json.Marshal(logger)
					fmt.Print(string(loggerData))
				}

				fmt.Println("whatsapp sent in go rutine")
				wg.Done()
			}()

			dynamoDb := services.GetDynamoClient()
			item := map[string]string{
				"phoneNumber": userCellPhone,
				"name":        commonName,
				"familyName":  familyName,
				"status":      message,
				"createdAt":   utils.GetCurrentTime(),
			}

			go func() {
				err := services.PutItem(ctx, item, dynamoDb)

				if err != nil {
					err := fmt.Sprint(err)
					logger["error"] = err
					logger["status_code"] = "400"
					loggerData, _ := json.Marshal(logger)
					fmt.Print(string(loggerData))
				}

				fmt.Println("response saved in dynamo")
				wg.Done()
			}()

			wg.Wait()

			logger["event_name"] = "send fast button and save survey"
		}

		logger["status_code"] = okRequestStr
		services.Logger(logger)

		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "Successful event received",
		}

		return response, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       "request nor post nor get",
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
