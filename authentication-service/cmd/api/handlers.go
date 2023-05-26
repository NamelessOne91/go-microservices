package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type requestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type logEntry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate the user against dB
	user, err := app.Repo.GetByEmail(payload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := app.Repo.PasswordMatches(payload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

	resPayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJSON(w, http.StatusAccepted, resPayload)
}

func (app *Config) logRequest(name, data string) error {
	entry := logEntry{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.Marshal(entry)
	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = app.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
