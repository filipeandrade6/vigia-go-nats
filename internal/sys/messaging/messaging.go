package messaging

import (
	"github.com/nats-io/nats.go"
)

type Config struct {
	Name string
	URL  string
}

func Connect(cfg Config) (*nats.Conn, error) {
	opts := []nats.Option{nats.Name(cfg.Name)}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
