package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/capitan-beto/vale-backend/api"
	"github.com/capitan-beto/vale-backend/internal/tools"
	"github.com/capitan-beto/vale-backend/models"
	"github.com/capitan-beto/vale-backend/pkg/utils"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	log "github.com/sirupsen/logrus"
)

func Payment(w http.ResponseWriter, r *http.Request) {
	xSignature := r.Header.Get("x-signature")
	xRequestId := r.Header.Get("x-request-id")
	queryParams := r.URL.Query()
	dataID := queryParams.Get("data.id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	parts := strings.Split(xSignature, ",")
	var ts, hash string

	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			if key == "ts" {
				ts = value
			} else if key == "v1" {
				hash = value
			}
		}
	}

	manifest := fmt.Sprintf("id:%v;request-id:%v;ts:%v;", dataID, xRequestId, ts)

	if err = utils.VerifySignature(hash, manifest); err != nil {
		log.Error(err)
		api.UnauthorizedErrorHandler(w)
		return
	}

	var payload models.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Request body: %s", string(body))
		api.InternalErrorHandler(w)
		return
	}

	id, err := strconv.Atoi(payload.Id["id"])
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println(id)

	accessToken := os.Getenv("MP_SECRET")
	cfg, err := config.New(accessToken)
	if err != nil {
		log.Error(err)
		return
	}

	client := payment.NewClient(cfg)

	resources, err := client.Get(context.Background(), id)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println(resources.ExternalReference)
	fmt.Println(resources.ID)

	var database *sql.DB
	database, err = tools.CreateConnection()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = tools.CheckPayment(resources.ExternalReference, database)
	if err != nil {
		log.Error(err)
	}

	var res = api.PaymentResponse{
		Message: "Webhook received!",
		Code:    http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
