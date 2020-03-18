package schema

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//OpenDb ...
func OpenDb() {
	config := &aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
	}
	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Cost"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Cost"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Name"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("NewProducts"),
	}
	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
