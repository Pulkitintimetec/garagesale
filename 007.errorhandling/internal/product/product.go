package product

import (
	"context"
	"fmt"
	"os"
	"time"

	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/platform/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// //BookRepo ...
// type BookRepo struct {
// 	dynamodbiface.DynamoDBAPI
// }
// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("product not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")

	ErrForbidden = errors.New("Attempted action is not allowed")
)

//GetAllData ...
func GetAllData(ctx context.Context) ([]ProductStructure, error) {
	ctx, span := trace.StartSpan(ctx, "internal.product.List")
	defer span.End()
	svc := database.Open()
	// config := &aws.Config{
	// 	Region:      aws.String("us-west-2"),
	// 	Endpoint:    aws.String("http://localhost:8000"),
	// 	Credentials: credentials.NewStaticCredentials("123", "123", ""),
	// }

	// sess := session.Must(session.NewSession(config))
	// svc := dynamodb.New(sess)
	tableName := "NewProducts"
	proj := expression.NamesList(expression.Name("Name"), expression.Name("Cost"), expression.Name("Quantity"), expression.Name("DateCreated"), expression.Name("DateUpdated"), expression.Name("id"))
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
	return items, err
}

//GetDatabyName to fetch particular data by there name
func GetDatabyName(ctx context.Context, Name, Cost string) (ProductStructure, error) {
	svc := database.Open()
	productData := ProductStructure{}
	productData.Name = Name
	productData.Cost = Cost
	//	av, err := dynamodbattribute.MarshalMap(b)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("NewProducts"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {

				S: aws.String(productData.Name),
			},
			"Cost": {
				S: aws.String(productData.Cost),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return productData, err

	}

	item := ProductStructure{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	if item.Name == "" {
		fmt.Println("Not Get any Data", err.Error())
	}
	return item, nil
}

//InsertProduct ...
func InsertProduct(ctx context.Context, user auth.Claims, np NewProduct, now time.Time) (ProductStructure, error) {
	p := ProductStructure{
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
		ID:          np.ID,
		UserID:      user.Subject,
	}
	svc := database.Open()
	av, err := dynamodbattribute.MarshalMap(p)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("NewProducts"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return p, err
	}

	// fmt.Printf("We have inserted a new item!\ns")
	return p, nil
}

// GetProductByID used to get the product by using Id
func GetProductByID(ctx context.Context, productID string) (ProductStructure, error) {
	svc := database.Open()
	tableName := "NewProducts"
	filt := expression.Name("id").Equal(expression.Value(productID))
	proj := expression.NamesList(expression.Name("Name"), expression.Name("Cost"), expression.Name("Quantity"), expression.Name("DateCreated"), expression.Name("DateUpdated"), expression.Name("id"), expression.Name("UserID"))
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
	items := ProductStructure{}
	for _, i := range result.Items {
		// item := ProductStructure{}

		err = dynamodbattribute.UnmarshalMap(i, &items)

		if err != nil {
			fmt.Println("Got error in unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// items = append(items, item)
	}
	return items, err
}

// UpdateProduct used to update the data of Products
func UpdateProduct(ctx context.Context, user auth.Claims, productID string, update UpdateProductStructure, now time.Time) error {
	p, err := GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	if !user.HasRole(auth.RoleAdmin) && p.UserID != user.Subject {
		return ErrForbidden
	}
	svc := database.Open()
	p.DateUpdated = now
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			// ":c": {
			// 	S: aws.String(*update.Cost),
			// },
			":q": {
				S: aws.String(*update.Quantity),
			},
		},
		TableName: aws.String("NewProducts"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(p.Name),
			},
			"Cost": {
				S: aws.String(p.Cost),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set  Quantity = :q"),
	}

	_, er := svc.UpdateItem(input)
	if er != nil {
		fmt.Println(er.Error())
		return errors.Wrap(err, "updating product")
	}
	return nil
}

// DeleteProduct used to delete the product when user provide the ProductID
func DeleteProduct(ctx context.Context, productID string) error {
	tableName := "NewProducts"
	svc := database.Open()
	p, err := GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Cost": {
				S: aws.String(p.Cost),
			},
			"Name": {
				S: aws.String(p.Name),
			},
		},
		TableName: aws.String(tableName),
	}

	data, erro := svc.DeleteItem(input)
	fmt.Print(data)
	if err != nil {
		return errors.Wrap(erro, "Got error calling DeleteItem")
	}
	return nil

}
