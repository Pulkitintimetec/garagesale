package main

import (
	"context"
	"encoding/json"
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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Start Database
	usingDb()
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
func usingDb() {
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
			TableName: aws.String("Products"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}
	fmt.Printf("We have processed %v records\n", len(product))
}

// ProductStructure is an item we sell.
type ProductStructure struct {
	Name     string `json:"Name"`
	Cost     string `json:"Cost"`
	Quantity string `json:"Quantity"`
	// DateCreated time.Time `db:"date_created" json:"date_created"`
	// DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
