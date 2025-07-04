package main

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
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
	"forge.capytal.company/capytalcode/project-comicverse/templates"
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

	privateKeyEnv = os.Getenv("PRIVATE_KEY")
	publicKeyEnv  = os.Getenv("PUBLIC_KEY")
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
	case privateKeyEnv == "":
		log.Fatal("PRIVATE_KEY not be a empty value")
	case publicKeyEnv == "":
		log.Fatal("PUBLIC_KEY not be a empty value")
	}
}

func main() {
	ctx := context.Background()

	level := slog.LevelError
	if *dev {
		level = slog.LevelDebug
	} else if *verbose {
		level = slog.LevelInfo
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	assertions := tinyssert.NewDisabled()
	if *dev {
		assertions = tinyssert.New(
			tinyssert.WithPanic(),
			tinyssert.WithLogger(log),
		)
	}

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
		d := os.DirFS("./assets")
		opts = append(opts, comicverse.WithAssets(d))

		t := templates.NewHotTemplates(os.DirFS("./templates"))
		opts = append(opts, comicverse.WithTemplates(t))

		opts = append(opts, comicverse.WithDevelopmentMode())
	}

	// TODO: Move this to dedicated function
	privateKeyStr, err := base64.URLEncoding.DecodeString(privateKeyEnv)
	if err != nil {
		log.Error("Failed to decode PRIVATE_KEY from base64", slog.String("error", err.Error()))
		os.Exit(1)
	}
	publicKeyStr, err := base64.URLEncoding.DecodeString(publicKeyEnv)
	if err != nil {
		log.Error("Failed to decode PUBLIC_KEY from base64", slog.String("error", err.Error()))
		os.Exit(1)
	}

	edPrivKey := ed25519.PrivateKey(privateKeyStr)
	edPubKey := ed25519.PublicKey(publicKeyStr)

	if len(edPrivKey) != ed25519.PrivateKeySize {
		log.Error("PRIVATE_KEY is not of valid size", slog.Int("size", len(edPrivKey)))
		os.Exit(1)
	}
	if len(edPubKey) != ed25519.PublicKeySize {
		log.Error("PUBLIC_KEY is not of valid size", slog.Int("size", len(edPubKey)))
		os.Exit(1)
	}

	if !edPubKey.Equal(edPrivKey.Public()) {
		log.Error("PUBLIC_KEY is not equal from extracted public key",
			slog.String("extracted", fmt.Sprintf("%x", edPrivKey.Public())),
			slog.String("key", fmt.Sprintf("%x", edPubKey)),
		)
		os.Exit(1)
	}

	app, err := comicverse.New(comicverse.Config{
		DB:         db,
		S3:         storage,
		PrivateKey: edPrivKey,
		PublicKey:  edPubKey,
		Bucket:     s3Bucket,
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
