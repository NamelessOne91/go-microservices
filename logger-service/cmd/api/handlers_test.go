package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test_write_log(t *testing.T) {
	tests := []struct {
		name           string
		bodyError      bool
		postBody       map[string]string
		jsonToReturn   string
		expectedStatus int
	}{
		{
			name:      "valid log request",
			bodyError: false,
			postBody: map[string]string{
				"name": "test",
				"data": "some log data",
			},
			jsonToReturn: `{
				"error": false,
				"message": "logged"
			}`,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:      "missing name in post body",
			bodyError: false,
			postBody: map[string]string{
				"data": "some log data",
			},
			jsonToReturn:   `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "missing data in post body",
			bodyError: false,
			postBody: map[string]string{
				"name": "test",
			},
			jsonToReturn:   `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error reading request body",
			bodyError:      true,
			postBody:       nil,
			jsonToReturn:   `{}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			if tt.bodyError {
				req, _ = http.NewRequest("POST", "/log", errReader{})
			} else {
				body, _ := json.Marshal(tt.postBody)
				req, _ = http.NewRequest("POST", "/log", bytes.NewReader(body))
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.WriteLog)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected HTTP %d but got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
