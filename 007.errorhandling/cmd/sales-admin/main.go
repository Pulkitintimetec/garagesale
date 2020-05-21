package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/platform/conf"
	"garagesale/007.errorhandling/internal/platform/database"
	"garagesale/007.errorhandling/internal/schema"
	"garagesale/007.errorhandling/internal/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:DynamoDB"`
			Password   string `conf:"default:DynamoDB,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:DynamoDB"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "error: parsing config")
	}

	// This is used for multiple commands below.
	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}
	flag.Parse()
	switch flag.Arg(0) {

	case "migrate":
		schema.OpenDb()
		log.Println("Migrations complete")
	case "salestable":
		schema.CreatingSalesTable()
		log.Println("Sales Table Generated")
	case "usertable":
		schema.CreatingUserTable()
		log.Println("Users table Created")
	case "keygen":
		err := keygen(cfg.Args.Num(1))
		fmt.Print(err)
	case "seed":
		schema.UsingDb()
		log.Println("Seed data complete")
	case "useradd":
		err := useradd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2), cfg.Args.Num(3))
		fmt.Print(err)
	}
	return nil
}
func useradd(cfg database.Config, name, email, password string) error {

	if email == "" || password == "" {
		return errors.New("useradd command must be called with three additional arguments for name,email and password")
	}

	fmt.Printf("Admin user will be created with email %q and password %q\n", email, password)
	fmt.Print("Continue? (1/0) ")

	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	ctx := context.Background()

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	u, err := user.Create(ctx, nu, time.Now())
	if err != nil {
		return err
	}
	fmt.Println("User created with hash password:", u.PasswordHash)
	fmt.Println("User created with id:", u.ID)
	return nil
}

// keygen creates an x509 private key for signing auth tokens.
func keygen(path string) error {
	if path == "" {
		return errors.New("keygen missing argument for key path")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "generating keys")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer file.Close()

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	if err := file.Close(); err != nil {
		return errors.Wrap(err, "closing private file")
	}

	return nil
}
