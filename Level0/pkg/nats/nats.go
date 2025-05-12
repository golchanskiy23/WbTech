package nats

import (
	"Level0/config"
	"github.com/nats-io/stan.go"
	"os"
)

type NatsSource struct {
	Conn stan.Conn
}

func NewNatsSource(cfg *config.NATS) (*NatsSource, error) {
	conn, err := stan.Connect(os.Getenv("CLUSTER_ID"), os.Getenv("SECOND_CLIENT_ID"), stan.NatsURL(cfg.URL))
	if err != nil {
		return nil, err
	}
	return &NatsSource{Conn: conn}, nil
}

func (src *NatsSource) Close() error {
	if err := src.Conn.Close(); err != nil {
		return err
	}
	return nil
}
