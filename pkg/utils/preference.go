package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/capitan-beto/vale-backend/models"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	log "github.com/sirupsen/logrus"
)

func Preference(pd models.ContestantData) (string, string, string, error) {
	cfg, err := config.New(os.Getenv("MP_SECRET"))
	if err != nil {
		log.Error(err)
		return "", "", "", err
	}

	request := preference.Request{
		ExternalReference: pd.ExtRef,
		BackURLs: &preference.BackURLsRequest{
			Success: "https://www.instagram.com/carlosnana1/",
			Failure: "https://www.instagram.com/milo.aleuzirb/",
		},
		Items: []preference.ItemRequest{
			{
				ID:        pd.Name,
				Title:     "Número de rifa solidaria",
				Quantity:  1,
				UnitPrice: 3000,
			},
		},
	}

	client := preference.NewClient(cfg)

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return "", "", resource.BackURLs.Failure, err
	}

	fmt.Println(resource.ID)
	return resource.InitPoint, resource.ID, resource.BackURLs.Success, nil
}
