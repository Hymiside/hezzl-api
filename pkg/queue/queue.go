package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type clickhouse interface {
	CreateLogs(ctx context.Context, logs []models.Good) error
}

type Queue struct {
	nats       *nats.Conn
	clickhouse clickhouse

	logs      []models.Good
	numOfLogs int
}

func NewQueue(
	nc *nats.Conn,
	clickhouse clickhouse,
	numOfLogs int,
) *Queue {
	return &Queue{
		nats:       nc,
		clickhouse: clickhouse,
		logs:       []models.Good{},
		numOfLogs:  numOfLogs,
	}
}

func (q *Queue) Subscribe() error {
	if _, err := q.nats.Subscribe("logs", func(m *nats.Msg) {
		if err := q.read(m.Data); err != nil {
			log.Errorf("error to read: %v", err)
		}
	}); err != nil {
		return fmt.Errorf("error to subscribe: %v", err)
	}

	return nil
}

func (q *Queue) Publish(b []byte) error {
	if err := q.nats.Publish("logs", b); err != nil {
		return fmt.Errorf("error to publish: %v", err)
	}
	return nil
}

func (q *Queue) read(b []byte) error {
	if len(q.logs) == q.numOfLogs {
		if err := q.clickhouse.CreateLogs(context.Background(), q.logs); err != nil {
			return fmt.Errorf("error to create logs: %v", err)
		}
		q.logs = []models.Good{}
		return nil
	}

	log := models.Good{}
	if err := json.Unmarshal(b, &log); err != nil {
		return fmt.Errorf("error to unmarshal: %v", err)
	}
	q.logs = append(q.logs, log)
	return nil
}
