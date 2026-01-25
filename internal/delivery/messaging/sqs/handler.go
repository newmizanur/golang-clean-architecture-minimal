package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	Handle(ctx context.Context, message types.Message) error
}
