package app

import (
	"Level0/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"math/rand"
	"time"
)

const (
	MinTime   = 10 * time.Millisecond
	MaxTime   = 5000 * time.Millisecond
	ClusterID = "CLUSTER_ID"
	ClientID  = "CLIENT_ID"
)

func getConnection(host, port string) string {
	return fmt.Sprintf("nats://%s:%s", host, port)
}

/*func ExecutePublisher() error {
	// localhost -> nats
	sc, err := stan.Connect(os.Getenv(ClusterID), os.Getenv(ClientID), stan.NatsURL(getConnection(os.Getenv("NATS_HOST"), "4222")))
	if err != nil {
		return fmt.Errorf("can't connect to nats server: %v", err)
	}
	lastOrder, err := utils.GetGivenOrder()
	if err != nil {
		return fmt.Errorf("can't get given order: %v", err)
	}

	for counter := 1; ; counter++ {
		marshalled, err := json.Marshal(lastOrder)
		if err != nil {
			return fmt.Errorf("can't marshal order: %v", err)
		}

		if err = sc.Publish(Channel, marshalled); err != nil {
			return fmt.Errorf("can't publish order: %v", err)
		}

		lastOrder = utils.RandomOrder()
		jitter := func(min, max time.Duration) time.Duration {
			if min >= max {
				return min
			}
			delta := max - min
			return min + time.Duration(rand.Int63n(int64(delta)))
		}(MinTime, MaxTime)
		time.Sleep(jitter)
	}
}*/

func ExecutePublisher(js nats.JetStreamContext, conn *nats.Conn) error {
	defer conn.Close()
	lastOrder, err := utils.GetGivenOrder()
	if err != nil {
		return fmt.Errorf("can't get given order: %v", err)
	}

	for {
		marshalled, err := json.Marshal(lastOrder)
		if err != nil {
			return fmt.Errorf("can't marshal order: %v", err)
		}

		_, err = js.Publish(Channel, marshalled)
		if err != nil {
			return fmt.Errorf("can't publish order: %v", err)
		}

		lastOrder = utils.RandomOrder()
		jitter := func(min, max time.Duration) time.Duration {
			if min >= max {
				return min
			}
			delta := max - min
			return min + time.Duration(rand.Int63n(int64(delta)))
		}(MinTime, MaxTime)
		time.Sleep(jitter)
	}
}
