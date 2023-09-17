package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"main/internal/domain/dto"
	"main/internal/domain/entity"
	"main/internal/services"
	"main/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	token := "supersecreto"

	if request.HTTPMethod == "GET" &&
		request.QueryStringParameters["hub.mode"] == "subscribe" &&
		request.QueryStringParameters["hub.verify_token"] == token {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       request.QueryStringParameters["hub.challenge"],
		}
		return response, nil
	}

	if request.HTTPMethod == "POST" {
		var whatsappEvent dto.WhatsAppEvent
		logger := map[string]string{
			"user_agent": request.Headers["User-Agent"],
			"ips":        request.Headers["X-Forwarded-For"],
			"ip":         request.RequestContext.Identity.SourceIP,
		}
		err := json.Unmarshal([]byte(request.Body), &whatsappEvent)
		if err != nil {
			logger["error"] = "Request body could not be serialized as whatsapp event"
			logger["status_code"] = "400"
			logger["body"] = err.Error()
			services.Logger(logger)
			response := events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Json error",
			}
			return response, nil
		}

		if !utils.SliceHasElements(whatsappEvent.Entry) {
			return utils.HandleError(logger, 400, "No entry found in the request")
		}

		entry := whatsappEvent.Entry[0]

		if !utils.SliceHasElements(entry.Changes) {
			return utils.HandleError(logger, 400, "No changes found in the request")
		}

		change := entry.Changes[0]
		field := change.Field
		logger["field"] = field
		if change.Field != "messages" {
			return utils.HandleError(logger, 400, "Field not supported")
		}

		value := change.Value

		if !utils.SliceHasElements(value.Contacts) && !utils.SliceHasElements(value.Statuses) {
			fmt.Print(request.Body)
			return utils.HandleError(logger, 400, "No contacts and no Statuses in the request")
		} else if !utils.SliceHasElements(value.Contacts) && utils.SliceHasElements(value.Statuses) {
			status := value.Statuses[0]
			logger["status_code"] = "200"
			logger["event_name"] = status.Status
			logger["from"] = status.RecipientID
			logger["event_id"] = status.ID
			logger["list_name"] = entity.ListOfPhonesNames[status.RecipientID]
			services.Logger(logger)
			response := events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Successful event received",
			}
			return response, nil
		}

		userName := value.Contacts[0].Profile.Name
		logger["user"] = userName

		if !utils.SliceHasElements(value.Messages) {
			return utils.HandleError(logger, 400, "No Messages found in the request")
		}
		userCellPhone := value.Messages[0].From
		logger["from"] = userCellPhone
		logger["list_name"] = entity.ListOfPhonesNames[userCellPhone]
		logger["message_id"] = value.Messages[0].ID
		if field != "messages" {
			logger["error"] = "Field not supported"
			logger["status_code"] = "400"
			services.Logger(logger)
			response := events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Field not supported",
			}
			return response, nil
		}

		message := value.Messages[0].Text.Body
		typeMessage := value.Messages[0].Type
		logger["event_name"] = "received message"
		logger["message_received"] = message
		if typeMessage == "text" && strings.ToLower(message) == "si" {
			familyName := entity.ListOfGuestsFamilyNames[userCellPhone]
			err := services.SendInvitation(userCellPhone, familyName)
			if err != nil {
				logger["reason"] = err.Error()
				logger["status_code"] = "400"
				logger["error"] = "Error sending whatsapp"
				services.Logger(logger)
				response := events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Error sending whatsapp",
				}
				return response, nil
			}
			imageId := entity.ListOfGuestsImage[userCellPhone]
			err = services.SendImage(userCellPhone, imageId)
			if err != nil {
				logger["reason"] = err.Error()
				logger["status_code"] = "400"
				logger["error"] = "Error sending whatsapp/image"
				services.Logger(logger)
				response := events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Error sending whatsapp",
				}
				return response, nil
			}
			logger["event_name"] = "send whatsapp image and text"
		}
		logger["status_code"] = "200"
		services.Logger(logger)
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Successful event received",
		}
		return response, nil

	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Bad Request",
	}
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
