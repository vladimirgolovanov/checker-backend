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

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	grab_instagram "github.com/vladimirgolovanov/grab-proto/gen/instagram"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vladimirgolovanov/checker-backend/handlers"
	"github.com/vladimirgolovanov/checker-backend/namespaces"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	registry := NewCheckerRegistry(Deps{
		InstagramClient: grab_instagram.NewInstagramClient(conn),
		Fetcher:         namespaces.NewHTTPFetcher(&http.Client{Timeout: 5 * time.Second}),
	})

	appHost := os.Getenv("APP_HOST")
	if appHost == "" {
		appHost = "0.0.0.0:8080"
	}

	http.HandleFunc("/api/check_name", handlers.CheckNameHandler(ctx, registry))
	srv := &http.Server{Addr: appHost}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %s", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %s", err)
	}
}
