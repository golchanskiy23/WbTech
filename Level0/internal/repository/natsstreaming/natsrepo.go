package natsstreaming

import (
	natsPkg "Level0/pkg/nats"
)

type NatsStreamingRepository struct {
	NatsSrc *natsPkg.NatsSource
}

func CreateNewNatsStreamingRepository(nats *natsPkg.NatsSource) NatsStreamingRepository {
	return NatsStreamingRepository{NatsSrc: nats}
}
