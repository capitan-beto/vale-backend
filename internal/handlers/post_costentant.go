package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/capitan-beto/vale-backend/api"
	"github.com/capitan-beto/vale-backend/internal/tools"
	"github.com/capitan-beto/vale-backend/models"
	log "github.com/sirupsen/logrus"
)

func AddContestant(w http.ResponseWriter, r *http.Request) {
	var err error
	var contestantId int64

	var database *sql.DB
	database, err = tools.CreateConnection()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var ContestantData models.ContestantData
	json.NewDecoder(r.Body).Decode(&ContestantData)

	defer r.Body.Close()

	contestantId, err = tools.NewContestant(&ContestantData, database)
	if err != nil {
		log.Error(err)
		return
	}

	res := api.AddContestantResponse{
		Code:    http.StatusOK,
		Id:      int(contestantId),
		Name:    ContestantData.Name,
		Created: time.Now().Format("2006-01-02 15:04:05"),
		ExtRef:  ContestantData.ExtRef,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
