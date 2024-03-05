package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) readJSON(r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
