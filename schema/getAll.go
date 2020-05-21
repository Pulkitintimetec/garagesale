package schema

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
	config := &aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
	}

	sess := session.Must(session.NewSession(config))
	svc := dynamodb.New(sess)
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

// List ...
func List(w http.ResponseWriter, req *http.Request) {
	// var br BookRepo
	getdata := GetAllData()
	data, err := json.Marshal(getdata)
	if err != nil {
		fmt.Print("error in getting data from database", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", data)

}
