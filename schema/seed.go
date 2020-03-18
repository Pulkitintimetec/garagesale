package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//UsingDb ...
func UsingDb() {
	productData, err := os.Open("jsonProduct.json")
	if err != nil {
		log.Fatalln("unable to open json file", err)
	}

	var product []ProductStructure
	err = json.NewDecoder(productData).Decode(&product)
	if err != nil {
		log.Fatal("error in encoding", err.Error())
		os.Exit(1)
	}

	config := &aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
	}

	sess := session.Must(session.NewSession(config))
	svc := dynamodb.New(sess)

	for _, productIteration := range product {
		info, err := dynamodbattribute.MarshalMap(productIteration)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal the movie, %v", err))
		}
		input := &dynamodb.PutItemInput{
			Item:      info,
			TableName: aws.String("NewProducts"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}
	fmt.Printf("We have processed %v records\n", len(product))
}
