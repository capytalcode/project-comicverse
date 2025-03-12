package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	comicverse "forge.capytal.company/capytalcode/project-comicverse"
	"forge.capytal.company/loreddev/x/tinyssert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/tursodatabase/go-libsql"
)

var (
	hostname     = flag.String("hostname", "localhost", "Host to listen to")
	port         = flag.Uint("port", 8080, "Port to be used for the server.")
	templatesDir = flag.String("templates", "", "Templates directory to be used instead of built-in ones.")
	verbose      = flag.Bool("verbose", false, "Print debug information on logs")
	dev          = flag.Bool("dev", false, "Run the server in debug mode.")
)

var (
	databaseURL = getEnv("DATABASE_URL", "file:./libsql.db")

	awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsDefaultRegion   = os.Getenv("AWS_DEFAULT_REGION")
	awsEndpointURL     = os.Getenv("AWS_ENDPOINT_URL")
	s3Bucket           = os.Getenv("S3_BUCKET")
)

func getEnv(key string, d string) string {
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	return v
}

func init() {
	flag.Parse()

	switch {
	case databaseURL == "":
		log.Fatal("DATABASE_URL should not be a empty value")
	case awsAccessKeyID == "":
		log.Fatal("AWS_ACCESS_KEY_ID should not be a empty value")
	case awsDefaultRegion == "":
		log.Fatal("AWS_DEFAULT_REGION should not be a empty value")
	case awsEndpointURL == "":
		log.Fatal("AWS_ENDPOINT_URL should not be a empty value")
	case s3Bucket == "":
		log.Fatal("S3_BUCKET should not be a empty value")
	}
}

func main() {
	ctx := context.Background()

	assertions := tinyssert.NewDisabledAssertions()
	if *dev {
		assertions = tinyssert.NewAssertions(tinyssert.Opts{
			Panic: true,
		})
	}

	level := slog.LevelError
	if *dev {
		level = slog.LevelDebug
	} else if *verbose {
		level = slog.LevelInfo
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	db, err := sql.Open("libsql", databaseURL)
	if err != nil {
		log.Error("Failed open connection to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	credentials := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     awsAccessKeyID,
			SecretAccessKey: awsSecretAccessKey,
			CanExpire:       false,
		}, nil
	})
	storage := s3.New(s3.Options{
		AppID:        "comicverse-pre-alpha",
		BaseEndpoint: &awsEndpointURL,
		Region:       awsDefaultRegion,
		Credentials:  &credentials,
	})

	opts := []comicverse.Option{
		comicverse.WithContext(ctx),
		comicverse.WithAssertions(assertions),
		comicverse.WithLogger(log),
	}

	if *dev {
		d := os.DirFS("./static")
		opts = append(opts, comicverse.WithStaticFiles(d))

		opts = append(opts, comicverse.WithDevelopmentMode())
	}

	app, err := comicverse.New(comicverse.Config{
		DB:     db,
		S3:     storage,
		Bucket: s3Bucket,
	}, opts...)
	if err != nil {
		log.Error("Failed to initiate comicverse app", slog.String("error", err.Error()))
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *hostname, *port),
		Handler: app,
	}

	c, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Starting application",
			slog.String("host", *hostname),
			slog.Uint64("port", uint64(*port)),
			slog.Bool("verbose", *verbose),
			slog.Bool("development", *dev))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start application server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-c.Done()

	log.Info("Stopping application gracefully")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to stop application server gracefully", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("FINAL")
	os.Exit(0)
}
