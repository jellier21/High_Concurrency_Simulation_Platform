package business

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishConfirmed(
	ch *amqp.Channel,
	confirms <-chan amqp.Confirmation,
	exchange string,
	routingKey string,
	body []byte,
) error {
	if err := ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("发布消息失败: %w", err)
	}

	confirmation, ok := <-confirms
	if !ok {
		return errors.New("发布确认通道已关闭")
	}
	if !confirmation.Ack {
		return errors.New("RabbitMQ 未确认该消息")
	}

	return nil
}
