package natsstreaming

import natsPkg "Level0/pkg/nats"

/*type NatsStreamingRepository struct {
	NatsSrc *natsPkg.NatsSource
}

func CreateNewNatsStreamingRepository(nats *natsPkg.NatsSource) NatsStreamingRepository {
	return NatsStreamingRepository{NatsSrc: nats}
}*/

type NatsJetStreamRepository struct {
	NatsSrc *natsPkg.NatsSource
}

func NewNatsJetStreamRepository(nats *natsPkg.NatsSource) *NatsJetStreamRepository {
	return &NatsJetStreamRepository{NatsSrc: nats}
}
