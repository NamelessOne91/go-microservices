package main

import (
	"log"
	"time"

	"github.com/NamelessOne91/logger/data"
)

// RPCServer is the type for our RPC Server. Methods having this as a receiver
// are available over RPC, as long as they are exported
type RPCServer struct {
	repo data.Repository
}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	err := r.repo.Insert(data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to Mongo")
		return err
	}

	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}
