package product

import (
	"fmt"
	"os"

	"garagesale/006.configuration/internal/platform/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// //BookRepo ...
// type BookRepo struct {
// 	dynamodbiface.DynamoDBAPI
// }

//GetAllData ...
func GetAllData() []ProductStructure {
	svc := database.Open()
	// config := &aws.Config{
	// 	Region:      aws.String("us-west-2"),
	// 	Endpoint:    aws.String("http://localhost:8000"),
	// 	Credentials: credentials.NewStaticCredentials("123", "123", ""),
	// }

	// sess := session.Must(session.NewSession(config))
	// svc := dynamodb.New(sess)
	tableName := "NewProducts"
	proj := expression.NamesList(expression.Name("Name"), expression.Name("Cost"), expression.Name("Quantity"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	items := []ProductStructure{}

	for _, i := range result.Items {
		item := ProductStructure{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error in unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		items = append(items, item)
	}
	return items
}
