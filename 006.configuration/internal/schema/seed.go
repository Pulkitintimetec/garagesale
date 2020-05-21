package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"garagesale/006.configuration/internal/platform/database"
	"garagesale/006.configuration/internal/product"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//UsingDb ...
func UsingDb() {
	productData, err := os.Open("jsonProduct.json")
	if err != nil {
		log.Fatalln("unable to open json file", err)
	}

	var product []product.ProductStructure
	err = json.NewDecoder(productData).Decode(&product)
	if err != nil {
		log.Fatal("error in encoding", err.Error())
		os.Exit(1)
	}

	svc := database.Open()
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
