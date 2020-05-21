package product

import (
	"context"
	"fmt"
	"os"
	"time"

	"garagesale/007.errorhandling/internal/platform/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

// AddSales records a sales transaction for a single Product.
func AddSales(ctx context.Context, ns NewSale, productID string, now time.Time) (*Sale, error) {
	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now,
	}
	svc := database.Open()
	av, err := dynamodbattribute.MarshalMap(s)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("NewProductsSale"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return &s, err
	}

	// fmt.Printf("We have inserted a new item!\ns")
	return &s, nil
}

// ListSales gives all Sales for a Product.
func ListSales(ctx context.Context, productID string) ([]Sale, error) {
	svc := database.Open()
	tableName := "NewProductsSale"
	proj := expression.NamesList(expression.Name("ID"), expression.Name("ProductID"), expression.Name("Quantity"), expression.Name("Paid"), expression.Name("DateCreated"))
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

	items := []Sale{}

	for _, i := range result.Items {
		item := Sale{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error in unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		items = append(items, item)
	}
	return items, err
}

// ListS ...
func ListS(ctx context.Context, productID string) ([]Sale, error) {
	svc := database.Open()
	tableName := "NewProductsSale"
	filt := expression.Name("ProductID").Equal(expression.Value(productID))
	proj := expression.NamesList(expression.Name("ID"), expression.Name("ProductID"), expression.Name("Quantity"), expression.Name("Paid"), expression.Name("DateCreated"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	items := []Sale{}

	for _, i := range result.Items {
		item := Sale{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error in unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		items = append(items, item)
	}
	return items, err
}
