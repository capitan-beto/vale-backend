package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/capitan-beto/vale-backend/api"
	log "github.com/sirupsen/logrus"
)

func WelcomeMessage(w http.ResponseWriter, r *http.Request) {
	var err error

	res := struct {
		Code int
		Msg  string
	}{
		Code: http.StatusOK,
		Msg:  "Everything's ok champ!",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
