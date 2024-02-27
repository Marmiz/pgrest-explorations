package main

import (
	"encoding/json"
	"net/http"

	"github.com/marmiz/pgrest-explorations/internal/data"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "available"}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error("failed to write JSON", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
	}
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	dbUser, err := app.queries.GetUser(r.Context(), "olm_dev@jll.com")
	if err != nil {
		app.logger.Error("failed to get user", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
		return
	}

	// map the DB model to the API model
	u := data.NewFromDB(&dbUser)
	err = app.writeJSON(w, http.StatusOK, u, nil)
	if err != nil {
		app.logger.Error("failed to write JSON", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
	}
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}
