package nats

import (
	"Level0/config"
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
)

const (
	ClusterID      = "CLUSTER_ID"
	SecondClientID = "SECOND_CLIENT_ID"
)

type NatsSource struct {
	Conn stan.Conn
}

func NewNatsSource(cfg *config.NATS) (*NatsSource, error) {
	conn, err := stan.Connect(os.Getenv(ClusterID), os.Getenv(SecondClientID), stan.NatsURL(fmt.Sprintf("%s:%s", os.Getenv("NATS_URL"), cfg.Port)))
	if err != nil {
		return nil, fmt.Errorf("error during connection with nats: %v", err)
	}
	return &NatsSource{Conn: conn}, nil
}

func (src *NatsSource) Close() error {
	if err := src.Conn.Close(); err != nil {
		return fmt.Errorf("error during close nats connection: %v", err)
	}
	return nil
}
