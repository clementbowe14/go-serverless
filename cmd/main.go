package main

import (
	"os"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/clementbowe14/go-serverless/pkg/handlers"
)

var (
	dynamoClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	
	log.Printf("Now starting new session in aws region %s", region)
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return
	}

	dynamoClient = dynamodb.New(awsSession)

	log.Printf("dynamodb client successfully created.")

	lambda.Start(handler)
}

const tableName = "go-serverless-users"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	log.Printf("Current method request=%s", req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynamoClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dynamoClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynamoClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynamoClient)
	default:
		return handlers.UnhandledMethod()
	}
}
