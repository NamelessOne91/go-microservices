package main

import (
	"testing"

	"github.com/NamelessOne91/logger/data"
)

func Test_log_info(t *testing.T) {
	server := RPCServer{
		repo: testApp.Repo,
	}

	tests := []struct {
		name    string
		payload RPCPayload
		resp    string
		err     error
	}{
		{
			"successfull request",
			RPCPayload{Name: "RPC unit test", Data: "test"},
			"Processed payload via RPC: RPC unit test",
			nil,
		},
		{
			"missing name field in payload",
			RPCPayload{Data: "test"},
			"",
			data.ErrNoName,
		},
		{
			"missing data field in payload",
			RPCPayload{Name: "RPC unit test"},
			"",
			data.ErrNoData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ""
			err := server.LogInfo(tt.payload, &r)
			if tt.err != nil && err != tt.err {
				t.Errorf("expected error to be %v - got %v", tt.err, err)
			} else if tt.resp != r {
				t.Errorf("expect res to be %s, got %s", tt.resp, r)
			}
		})
	}
}
