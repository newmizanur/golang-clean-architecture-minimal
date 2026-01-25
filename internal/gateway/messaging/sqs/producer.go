package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

type Producer[T any] struct {
	Client   *sqs.Client
	QueueURL string
	Log      *logrus.Logger
}

func (p *Producer[T]) Send(ctx context.Context, payload T) error {
	body, err := json.Marshal(payload)
	if err != nil {
		p.Log.WithError(err).Error("failed to marshal SQS payload")
		return err
	}

	_, err = p.Client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.QueueURL,
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		p.Log.WithError(err).Error("failed to send SQS message")
		return err
	}

	return nil
}
