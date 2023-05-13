package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/mail"
	"github.com/rs/zerolog/log"
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
	mail   mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("failed to process task")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mail:   mailer,
	}
}

func (process *RedisTaskProcessor) Start() error {

	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, process.ProcessTaskSendVerifyEmail)

	return process.server.Start(mux)
}
