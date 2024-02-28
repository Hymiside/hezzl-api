package queue

import "github.com/nats-io/nats.go"

type Queue struct {
	Nats  *nats.Conn
}

func NewQueue(nats *nats.Conn) *Queue {
	return &Queue{Nats: nats}
}