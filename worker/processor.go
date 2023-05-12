package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/julysNICK/simplebank/db/sqlc"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (process *RedisTaskProcessor) Start() error {

	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, process.ProcessTaskSendVerifyEmail)

	return process.server.Start(mux)
}
