package main

import (
	"fmt"

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
		fmt.Println(string(request.Body))
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
