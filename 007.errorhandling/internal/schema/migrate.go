package schema

import (
	"fmt"

	"garagesale/007.errorhandling/internal/platform/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//OpenDb ...
func OpenDb() {
	svc := database.Open()
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

// CreatingSalesTable is used for crating table for Sales part
func CreatingSalesTable() {
	svc := database.Open()
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("ProductID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("ProductID"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("NewProductsSale"),
	}
	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}

// CreatingUserTable is used for crating table for User part
func CreatingUserTable() {
	svc := database.Open()
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("NewUser"),
	}
	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
