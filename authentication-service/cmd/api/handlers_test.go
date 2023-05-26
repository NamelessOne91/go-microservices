package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc struct {
	f func(*http.Request) *http.Response
	e error
}

// implement RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.f(req), f.e
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

var errLogService = errors.New("some log-service error")

func Test_authenticate(t *testing.T) {
	tests := []struct {
		name           string
		postBody       map[string]any
		jsonToReturn   string
		clientError    error
		expectedStatus int
	}{
		{
			name: "Accepted request",
			postBody: map[string]any{
				"email":    "me@here.com",
				"password": "verysecret",
			},
			jsonToReturn: `{
				"error": false,
				"message": "some message"
			}`,
			clientError:    nil,
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "Missing password in POST body",
			postBody: map[string]any{
				"email": "me@here.com",
			},
			jsonToReturn:   `{}`,
			clientError:    nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing email in POST body",
			postBody: map[string]any{
				"password": "verysecret",
			},
			jsonToReturn:   `{}`,
			clientError:    nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Client error",
			postBody: map[string]any{
				"email":    "me@here.com",
				"password": "verysecret",
			},
			jsonToReturn: `{
				"error": true,
				"message": "service unavalaible"
			}`,
			clientError:    errLogService,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mocks internal call to logger-service
			client := newTestClient(RoundTripFunc{
				f: func(r *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(tt.jsonToReturn)),
						Header:     make(http.Header),
					}
				},
				e: tt.clientError,
			})
			testApp.Client = client

			body, _ := json.Marshal(tt.postBody)
			req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.Authenticate)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected HTTP %d but got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func Test_log_request(t *testing.T) {
	tests := []struct {
		name        string
		entryName   string
		entryData   string
		clientError error
	}{
		{"successful log request", "authentication", "test@gmail.com logged in", nil},
		{"successful log request", "authentication", "test@gmail.com logged in", errLogService},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(RoundTripFunc{
				f: func(r *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
					}
				},
				e: tt.clientError,
			})
			testApp.Client = client

			err := testApp.logRequest(tt.entryName, tt.entryData)
			if err == nil && tt.clientError != nil {
				t.Errorf("expected logRequest to return an error - got %v", err)
			}
		})

	}
}
