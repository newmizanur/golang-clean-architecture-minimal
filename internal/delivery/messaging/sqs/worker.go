package sqs

import (
	"context"
	"time"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Worker struct {
	log *logrus.Logger
	cfg *viper.Viper
	queueURL string
	client   *sqs.Client
	handler  Handler
	wg  sync.WaitGroup
}

func NewWorker(log *logrus.Logger, cfg *viper.Viper, handler Handler) *Worker {
	return &Worker{
		log:     log,
		cfg:     cfg,
		handler: handler,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	queueURL := w.cfg.GetString("aws.sqs.queue_url")
	if queueURL == "" {
		w.log.Warn("SQS queue URL is empty; worker will exit")
		return
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(w.cfg.GetString("aws.region")))
	if err != nil {
		w.log.WithError(err).Error("Failed to load AWS config")
		return
	}

	w.queueURL = queueURL
	w.client = sqs.NewFromConfig(awsCfg)

	w.log.Info("SQS worker running")
	for {
		select {
		case <-ctx.Done():
			w.log.Info("SQS worker stopped")
			return
		default:
		}

		output, err := w.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &w.queueURL,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			w.log.WithError(err).Error("Failed to receive messages")
			time.Sleep(2 * time.Second)
			continue
		}

		for _, msg := range output.Messages {
			if err := w.handler.Handle(ctx, msg); err != nil {
				w.log.WithError(err).Error("Failed to handle message")
				continue
			}

			_, err = w.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &w.queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				w.log.WithError(err).Error("Failed to delete message")
			}
		}
	}
}

func (w *Worker) Wait() {
	w.wg.Wait()
}
