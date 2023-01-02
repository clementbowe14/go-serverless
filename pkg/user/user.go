package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/clementbowe14/go-serverless/pkg/validators"
)

var (
	ErrorFailedToFetchRecord = "Failed to fetch record"
	ErrorFailedToUnmarshalRecord = "Failed to unmarshal record"
	ErrorInvalidUserData = "Invalid user data"
	ErrorInvalidEmail = "Invalid email"
	ErrorCouldNotMarshalItem = "Could not marshal item"
	ErrorCouldNotDeleteItem = "Coud not delete item"
	ErrorCouldNotPutItem = "Could not put item"
	ErrorUserAlreadyExists = "User already exists"
	ErrorUserDoesNotExist = "User does not exist"
)

type User struct {
	Email 	string `json:"email`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchUser(email string, tableName string, dynamoClient dynamodbiface.DynamoDBAPI)(*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, err
}

func FetchUsers(tableName string, dynamoClient dynamodbiface.DynamoDBAPI)(*[]User, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)

	return item, nil

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User
	
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	//check if the user already exists
	currentUser, _ := FetchUser(u.Email, tableName, dynamoClient)
	if currentUser != nil  && len(currentUser.Email) != 0{
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:av,
		TableName:aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &u, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	currentUser, _ := FetchUser(u.Email, tableName, dynamoClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S:aws.String(email),
			},
		},
		TableName:aws.String(tableName),
	}
	
	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}