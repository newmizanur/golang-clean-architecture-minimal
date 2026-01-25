package sqs

import (
	"golang-clean-architecture/internal/model"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

type ContactProducer struct {
	Producer[model.CreateContactRequest]
}

func NewContactProducer(client *sqs.Client, queueURL string, log *logrus.Logger) *ContactProducer {
	return &ContactProducer{
		Producer: Producer[model.CreateContactRequest]{
			Client:   client,
			QueueURL: queueURL,
			Log:      log,
		},
	}
}
