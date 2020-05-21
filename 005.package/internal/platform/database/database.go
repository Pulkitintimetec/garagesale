package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//Open knows how to open the database connection
func Open() *dynamodb.DynamoDB {
	config := &aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
	}

	sess := session.Must(session.NewSession(config))
	svc := dynamodb.New(sess)
	return svc
}
