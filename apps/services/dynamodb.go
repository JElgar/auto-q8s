package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"log"
	"time"
)

type Dynamo struct {
	Session *dynamodb.DynamoDB
	tableName string
}

type ResultItem struct {
	ID		string `json:"id"`
	Status	string `json:"status"`
	CompletedAt time.Time `json:"completedAt"`
}


func InitDynamo() *Dynamo {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("eu-west-2"),
			Credentials: credentials.NewEnvCredentials(),
		},
	)
	if err != nil {
		log.Fatalf("Failed to connect to aws")
	} 

	dbSession := dynamodb.New(sess)
	return &Dynamo{
		Session: dbSession,
		tableName: os.Getenv("DYNAMO_TABLE"),
	}
} 

func (dynamo *Dynamo) PutItem(item *ResultItem) {
	attr, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Failed to marshall item %s", err)
	}

	fmt.Println("Putting item: ")
	fmt.Println(item)
	_, err = dynamo.Session.PutItem(
		&dynamodb.PutItemInput{
			Item:      attr,
			TableName: aws.String(dynamo.tableName),
		},
	)

	if err != nil {
		log.Fatalf("Failed to put item in table %s", err)
	}
}
