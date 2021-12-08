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
	Session dynamodb.Session
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

	session := dynamodb.New(sess)
	return &Dynamo{
		Session: session,
		tableName: os.Getenv("DYNAMO_TABLE"),
	}
} 

func (dynamo *Dynamo) PutItem(item *ResultItem) {
	attr, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Failed to marshall item %s", err)
	}

	_, err = dynamo.Session.PutItem(attr)
	if err != nil {
		log.Fatalf("Failed to put item in table %s", err)
	}
}
