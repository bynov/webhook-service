package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/bynov/webhook-service/internal/batch"
	"github.com/bynov/webhook-service/internal/config"
	"github.com/bynov/webhook-service/internal/domain"
	"github.com/bynov/webhook-service/internal/migration"
	"github.com/bynov/webhook-service/internal/storage"
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

	batchProvder := batch.New(repo, time.Second*10, 2000)

	uCase := usecase.NewUsecase(batchProvder, time.Second*5)

	r := chi.NewRouter()

	r.Post("/webhooks", addWebhook(uCase))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go batchProvder.Start()

	go func() {
		for err := range batchProvder.Errors() {
			log.Println(err)
		}
	}()

	server.ListenAndServe()
}

func addWebhook(usecase usecase.Usecase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO:
			w.Write([]byte(err.Error()))
			return
		}

		h := sha1.New()
		h.Write(payload)

		err = usecase.AddWebhook(r.Context(), domain.Webhook{
			Payload:     string(payload),
			PayloadHash: hex.EncodeToString(h.Sum(nil)),
			RecievedAt:  time.Now().UTC(),
		})
		if err != nil {
			// TODO:
			w.Write([]byte(err.Error()))
			return
		}
	}
}
