package user

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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")
)

// Create inserts a new user into the database.
func Create(ctx context.Context, n NewUser, now time.Time) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}
	fmt.Print(hash)
	u := User{
		ID:           uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}
	svc := database.Open()
	av, err := dynamodbattribute.MarshalMap(u)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("NewUser"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return &u, err
	}

	// fmt.Printf("We have inserted a new item!\ns")
	return &u, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims value representing this user. The claims can be
// used to generate a token for future authentication.
func Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	var u User
	svc := database.Open()
	tableName := "NewUser"
	filt := expression.Name("Email").Equal(expression.Value(email))
	proj := expression.NamesList(expression.Name("Name"), expression.Name("ID"), expression.Name("Email"), expression.Name("DateCreated"), expression.Name("DateUpdated"), expression.Name("Roles"), expression.Name("PasswordHash"))
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
	u = User{}
	for _, i := range result.Items {
		// item := ProductStructure{}

		err = dynamodbattribute.UnmarshalMap(i, &u)

		if err != nil {
			fmt.Println("Got error in unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// items = append(items, item)
	}
	// return items, err

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(u.ID, u.Roles, now, time.Hour)
	return claims, nil
}
