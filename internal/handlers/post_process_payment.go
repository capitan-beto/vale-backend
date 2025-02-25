package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/capitan-beto/vale-backend/api"
	"github.com/capitan-beto/vale-backend/models"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	log "github.com/sirupsen/logrus"
)

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var pd models.PaymentData
	json.NewDecoder(r.Body).Decode(&pd)

	defer r.Body.Close()

	cfg, err := config.New(os.Getenv("MP_SECRET"))
	if err != nil {
		fmt.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	request := preference.Request{
		ExternalReference: pd.ExternalReference,
		Items: []preference.ItemRequest{
			{
				ID:          pd.Item.ID,
				Title:       pd.Item.Description,
				Quantity:    pd.Item.Quantity,
				Description: pd.Item.Description,
				PictureURL:  pd.Item.PictureURL,
			},
		},
	}

	client := preference.NewClient(cfg)

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		api.PaymentError(w, resource.BackURLs.Failure)
		return
	}

	var res = api.PaymentResponse{
		Message: resource.SandboxInitPoint,
		BackURL: resource.BackURLs.Success,
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
