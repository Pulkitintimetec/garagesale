package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"

	"garagesale/007.errorhandling/internal/platform/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//UsingDb ...
func UsingDb() {
	userData, err := os.Open("jsonUser.json")
	if err != nil {
		log.Fatalln("unable to open json file", err)
	}

	var user user.User
	err = json.NewDecoder(userData).Decode(&user)
	if err != nil {
		log.Fatal("error in encoding", err.Error())
		os.Exit(1)
	}

	svc := database.Open()

	info, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal the Data, %v", err))
	}
	input := &dynamodb.PutItemInput{
		Item:      info,
		TableName: aws.String("NewUser"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// fmt.Printf("We have processed %v records\n", len(user))
}
