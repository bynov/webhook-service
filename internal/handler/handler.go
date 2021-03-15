package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bynov/webhook-service/internal/domain"
	"github.com/bynov/webhook-service/internal/usecase"
)

func AddWebhook(usecase usecase.Usecase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
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
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GetLiteWebhooks(usecase usecase.Usecase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		fr := query.Get("from")
		if fr == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("query param 'from' is missing in request"))
			return
		}

		from, err := strconv.ParseInt(fr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		t := query.Get("to")
		if t == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("query param 'to' is missing in request"))
			return
		}

		to, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		webhooks, err := usecase.GetLiteWebhooks(
			r.Context(),
			time.Unix(from, 0),
			time.Unix(to, 0),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		if len(webhooks) == 0 {
			_, _ = w.Write([]byte("[]"))
			return
		}

		hooks, err := json.Marshal(webhooks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(hooks)
	}
}

func GetWebhooksByIDs(usecase usecase.Usecase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["ids[]"]
		if len(ids) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("'ids[]' is empty or not provided"))
			return
		}

		webhooks, err := usecase.GetWebhooksByIDs(
			r.Context(),
			ids,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		hooks, err := json.Marshal(webhooks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(hooks)
	}
}
