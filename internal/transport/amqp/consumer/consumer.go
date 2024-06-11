package amqp

import (
	"context"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/service"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
)

func StartConsumers(amqp *amqpclient.Amqp, services *service.Services) error {
	_, err := NewMessageConsumer(amqp, func(ctx context.Context, message model.AmqpMessage) error {
		return services.AmqpMessage.SaveFromAmqp(ctx, message)
	})
	return err
}
