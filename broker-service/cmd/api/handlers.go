package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unkown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create the JSON we'll send to the auth microservice
	jsonData, _ := json.Marshal(a)
	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	// make sure we get back the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"), res.StatusCode)
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
