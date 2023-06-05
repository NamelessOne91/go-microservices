package main

import (
	"context"
	"testing"

	"github.com/NamelessOne91/logger/data"
	"github.com/NamelessOne91/logger/logs"
)

func Test_grpc_write_log(t *testing.T) {
	ctx := context.Background()
	logServer := LogServer{repo: testApp.Repo}

	type expected struct {
		res *logs.LogResponse
		err error
	}

	tests := []struct {
		name string
		req  *logs.LogRequest
		exp  expected
	}{
		{
			"successfull request",
			&logs.LogRequest{
				LogEntry: &logs.Log{
					Name: "gRPC unit test",
					Data: "gRPC unit test",
				},
			},
			expected{
				&logs.LogResponse{
					Result: "logged",
				},
				nil,
			},
		},
		{
			"missing name fied in entry",
			&logs.LogRequest{
				LogEntry: &logs.Log{
					Data: "gRPC unit test",
				},
			},
			expected{
				&logs.LogResponse{
					Result: "failed",
				},
				data.ErrNoName,
			},
		},
		{
			"missing data field in entry",
			&logs.LogRequest{
				LogEntry: &logs.Log{
					Name: "gRPC unit test",
				},
			},
			expected{
				&logs.LogResponse{
					Result: "failed",
				},
				data.ErrNoData,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := logServer.WriteLog(ctx, tt.req)
			if err != tt.exp.err {
				t.Errorf("expected error to be %v - got %v", tt.exp.err, err)
			}
			if res.Result != tt.exp.res.Result {
				t.Errorf("expected LogResponse's Result to be %s, got %s", tt.exp.res.Result, res.Result)
			}
		})
	}
}
