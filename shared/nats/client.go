package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
)

// InitNATS connects to NATS and returns a connection and JetStream instance
func InitNATS(url string) (*nats.Conn, jetstream.JetStream, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to NATS JetStream successfully.")
	return nc, js, nil
}
