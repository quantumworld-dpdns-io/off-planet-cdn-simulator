package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/config"
	"github.com/off-planet-cdn/control-api/internal/db"
	cdnredis "github.com/off-planet-cdn/control-api/internal/redis"
	"github.com/off-planet-cdn/control-api/internal/routes"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlphttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func initTracer(ctx context.Context, endpoint string) (*sdktrace.TracerProvider, error) {
	exp, err := otlphttp.New(ctx,
		otlphttp.WithEndpoint(endpoint),
		otlphttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("control-api"),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tp, err := initTracer(ctx, cfg.OtelEndpoint)
	if err != nil {
		log.Printf("warn: could not init tracer: %v", err)
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(shutdownCtx); err != nil {
				log.Printf("tracer shutdown: %v", err)
			}
		}()
	}

	dbClient, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer dbClient.Close()

	redisAddr := cfg.RedisURL
	if after, found := strings.CutPrefix(redisAddr, "redis://"); found {
		redisAddr = after
	}
	redisClient := cdnredis.New(redisAddr)
	defer redisClient.Close()

	r := chi.NewRouter()
	routes.Register(r, dbClient, redisClient)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("control-api listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down control-api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("control-api stopped")
}
