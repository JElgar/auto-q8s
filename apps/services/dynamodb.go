package services

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"log"
)

type Dynamo struct {
	Session *dynamodb.DynamoDB
	tableName string
}

type ResultItem struct {
	ID		string
    Status	string
}


func InitDynamo() *Dynamo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

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
