package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	grab_instagram "github.com/vladimirgolovanov/grab-proto/gen/instagram"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vladimirgolovanov/checker-backend/handlers"
	"github.com/vladimirgolovanov/checker-backend/namespaces"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		fmt.Printf("SENTRY_DSN: %s\n", sentryDsn)
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		log.Fatalf("GRPC_ADDR is not set")
	}
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %s", err)
	}
	defer conn.Close()

	CheckerRegistry[0] = single(namespaces.NewInstagramChecker(grab_instagram.NewInstagramClient(conn)))

	appHost := os.Getenv("APP_HOST")
	if appHost == "" {
		appHost = "0.0.0.0:8080"
	}

	http.HandleFunc("/api/check_name", handlers.CheckNameHandler(CheckerRegistry))
	log.Fatal(http.ListenAndServe(appHost, nil))
}
