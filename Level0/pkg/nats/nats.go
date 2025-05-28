package nats

import (
	"github.com/nats-io/nats.go"
)

type NatsSource struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

func NewNatsSource(conn *nats.Conn, js nats.JetStreamContext) (*NatsSource, error) {
	return &NatsSource{
		Conn:      conn,
		JetStream: js,
	}, nil
}

func (src *NatsSource) Close() error {
	if src.Conn != nil && !src.Conn.IsClosed() {
		src.Conn.Close()
	}
	return nil
}
