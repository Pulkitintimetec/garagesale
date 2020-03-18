package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Start Database
	openDb()
	// =========================================================================
	// Start API Service

	api := http.Server{
		Addr:         "localhost:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

func openDb() {
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
		TableName: aws.String("Products"),
	}
	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
