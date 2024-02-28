package queue

import (
	"context"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/nats-io/nats.go"
)

func NewNats(ctx context.Context, c models.ConfigNats) (*nats.Conn, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", c.Host, c.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect nats: %w", err)
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		nc.Close()
	}(ctx)

	return nc, err
}