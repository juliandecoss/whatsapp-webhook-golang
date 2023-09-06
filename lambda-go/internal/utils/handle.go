package utils

import (
	"main/internal/services"
	"reflect"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

func SliceHasElements(slice interface{}) bool {
	v := reflect.ValueOf(slice)
	return v.Kind() == reflect.Slice && v.Len() > 0
}

func HandleError(
	logger map[string]string,
	statusCode int,
	message string,
) (events.APIGatewayProxyResponse, error) {
	logger["error"] = message
	logger["status_code"] = strconv.Itoa(statusCode)
	services.Logger(logger)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}, nil
}
