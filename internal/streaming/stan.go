package streaming

import (
	"wbtech/level0/internal/config"

	"github.com/nats-io/stan.go"
)

func Connect(cfg *config.Config) (stan.Conn, error) {

	sc, err := stan.Connect(cfg.StanCluster, cfg.StanClient, stan.NatsURL("nats://stan:4222"))

	return sc, err
}
