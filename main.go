package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/bynov/webhook-service/internal/batch"

	"github.com/bynov/webhook-service/internal/config"
	"github.com/bynov/webhook-service/internal/handler"
	"github.com/bynov/webhook-service/internal/migration"
	"github.com/bynov/webhook-service/internal/storage"
	"github.com/bynov/webhook-service/internal/synchronizer"
	"github.com/bynov/webhook-service/internal/usecase"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.Connect(context.Background(), cfg.DatabaseAddr)
	if err != nil {
		panic(err)
	}

	err = migration.MigratePostgres("file://./assets/migrations", cfg.DatabaseAddr)
	if err != nil {
		panic(err)
	}

	repo := storage.NewRepository(pool)

	batchProvder := batch.New(repo, 2000)

	uCase := usecase.NewUsecase(batchProvder, repo, time.Second*5)

	r := chi.NewRouter()

	r.Post("/webhooks", handler.AddWebhook(uCase))

	if !cfg.IsMaster() {
		r.Get("/webhooks", handler.GetWebhooksByIDs(uCase))
		r.Get("/webhooks/lite", handler.GetLiteWebhooks(uCase))
	}

	if cfg.IsMaster() {
		slaveWebhookProvider := synchronizer.NewWebhookProviderFromSlave(cfg.SlaveAddr, time.Second*5)

		sync := synchronizer.New(repo, slaveWebhookProvider, batchProvder)

		go sync.Run(time.Minute * 2)

		go func() {
			for err := range sync.Errors() {
				log.Printf("sync provider error: %v", err)
			}
		}()
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go batchProvder.Start(time.Second * 10)

	go func() {
		for err := range batchProvder.Errors() {
			log.Printf("batch provider error: %v", err)
		}
	}()

	idleConnsClosed := make(chan struct{})
	go func() {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
		<-gracefulStop

		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		close(idleConnsClosed)
	}()

	go server.ListenAndServe()

	log.Println("HTTP server is started")

	<-idleConnsClosed
	log.Println("Service gracefully stopped")
}
