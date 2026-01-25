package sqs

import (
	"context"
	"encoding/json"

	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
)

type ContactHandler struct {
	log     *logrus.Logger
	usecase *usecase.ContactUseCase
}

func NewContactHandler(log *logrus.Logger, usecase *usecase.ContactUseCase) *ContactHandler {
	return &ContactHandler{
		log:     log,
		usecase: usecase,
	}
}

func (h *ContactHandler) Handle(ctx context.Context, message types.Message) error {
	if message.Body == nil {
		return nil
	}

	var payload model.CreateContactRequest
	if err := json.Unmarshal([]byte(*message.Body), &payload); err != nil {
		h.log.WithError(err).Error("Failed to unmarshal SQS message")
		return err
	}

	// Example: call usecase here.
	// _, err := h.usecase.Create(ctx, &payload)
	// return err
	h.log.WithField("userId", payload.UserId).Info("Received SQS contact message")
	return nil
}
