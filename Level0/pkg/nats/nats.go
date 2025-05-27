package nats

import (
	"github.com/nats-io/nats.go"
)

const (
	ClusterID      = "CLUSTER_ID"
	SecondClientID = "SECOND_CLIENT_ID"
)

type NatsSource struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

/*
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
}*/

const (
	NatsURLEnv  = "NATS_URL"
	NatsPortEnv = "NATS_PORT"
	// Для JetStream ClusterID и ClientID не нужны
	ClientIDEnv = "CLIENT_ID"
)

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
