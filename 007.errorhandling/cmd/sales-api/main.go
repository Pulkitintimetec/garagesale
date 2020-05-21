package main

import (
	"context"
	"crypto/rsa"
	_ "expvar" //Register the debug var handlers
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" //Register the pprof handler
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"contrib.go.opencensus.io/exporter/zipkin"
	"garagesale/007.errorhandling/cmd/sales-api/internal/handlers"
	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/platform/conf"
	"garagesale/007.errorhandling/internal/platform/database"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	log := log.New(os.Stdout, "SALES:", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address string `conf:"default:localhost:9000"`
			Debug   string `conf:"default:localhost:6060"`

			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:DynamoDB"`
			Password   string `conf:"default:DynamoDB,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:DynamoDB"`
			DisableTLS bool   `conf:"default:false"`
		}
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api"`
			Service     string  `conf:"default:sales-api"`
			Probability float64 `conf:"default:1"`
		}
	}

	// flag.BoolVar(&noVerifySsl, "no-verify-ssl", false, "Don't verify SSL certificates")

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")
	// =========================================================================
	// Get Configuration

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config: %s")
	}
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Initialize authentication support

	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}
	shutdown := make(chan os.Signal, 1)
	// =========================================================================
	// Start Database

	db := database.OpenConfiguration(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	fmt.Println(db)

	// =========================================================================
	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()
	// =========================================================================
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	//
	// Not concerned with shutting this down when the application is shutdown.
	go func() {
		log.Println("debug service listening on", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Println("debug service closed", err)
	}()
	// =========================================================================
	// Start API Service
	// ps := handlers.Product{Log: log}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(shutdown, log, authenticator),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.ShutdownTimeout,
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
	shutdowns := make(chan os.Signal, 1)
	signal.Notify(shutdowns, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case sigh := <-shutdowns:
		log.Println("main : Start shutdown", sigh)

		// Give outstanding requests a deadline for completion.
		// const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "shutdown gracefully ")
		}
		if sigh == syscall.SIGHUP {
			return errors.New("integity error detected,asking for self shutdown")
		}
	}
	return nil
}
func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}
func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
