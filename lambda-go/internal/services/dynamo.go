package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client
var once sync.Once

func GetDynamoClient() *dynamodb.Client {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			newErr := fmt.Errorf("unable to load SDK config, %w", err)
			panic(newErr)
		}

		if os.Getenv("env") == "local" {
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile("personaljulian"),
			)
			if err != nil {
				panic("unable to load SDK config, " + err.Error())
			}

			cfg.Region = "us-west-2"
		}

		dynamoClient = dynamodb.NewFromConfig(cfg)
	})

	return dynamoClient
}

func PutItem(ctx context.Context, item map[string]string, client *dynamodb.Client) error {
	data, err := attributevalue.MarshalMap(item)

	if err != nil {
		return fmt.Errorf("marshalMap: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("survey-invitations"),
		Item:      data,
	}
	_, err = client.PutItem(ctx, input)

	if err != nil {
		return fmt.Errorf("error inserting item: %w", err)
	}

	return nil
}
