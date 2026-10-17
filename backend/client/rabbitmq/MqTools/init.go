package MqTools

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(ch *amqp091.Channel) error {
	Exchanges := []struct {
		name string
		kind string
	}{
		{"order_exchange", "direct"},
		{"retry_exchange", "direct"},
		{"dead_exchange", "direct"},
	}

	for _, ex := range Exchanges {
		if err := ch.ExchangeDeclare(
			ex.name, ex.kind, true, false, false, false, nil,
		); err != nil {
			return fmt.Errorf("声明交换机 %s 失败: %w", ex.name, err)
		}
	}

	//-------------业务订单队列----------------
	if _, err := ch.QueueDeclare(
		"order_queue", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("声明订单队列失败: %w", err)
	}

	if err := ch.QueueBind(
		"order_queue", "order", "order_exchange", false, nil,
	); err != nil {
		return fmt.Errorf("绑定订单队列失败: %w", err)
	}

	//-------------重试队列----------------

	retryArgs := amqp091.Table{
		"x-message-ttl":             int32(5000),
		"x-dead-letter-exchange":    "order_exchange",
		"x-dead-letter-routing-key": "order",
	}

	if _, err := ch.QueueDeclare(
		"retry_queue", true, false, false, false, retryArgs,
	); err != nil {
		return fmt.Errorf("声明重试队列失败: %w", err)
	}
	if err := ch.QueueBind(
		"retry_queue", "retry", "retry_exchange", false, nil,
	); err != nil {
		return fmt.Errorf("绑定重试队列失败: %w", err)
	}

	//-------------死信队列----------------
	if _, err := ch.QueueDeclare(
		"dead_queue", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("声明死信队列失败: %w", err)
	}
	if err := ch.QueueBind(
		"dead_queue", "dead", "dead_exchange", false, nil,
	); err != nil {
		return fmt.Errorf("绑定死信队列失败: %w", err)
	}

	return nil
}

func InitMQChannel() (*amqp091.Channel, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("连接RabbitMQ失败: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("创建通道失败: %w", err)
	}

	if err = InitRabbitMQ(ch); err != nil {
		return nil, fmt.Errorf("初始化RabbitMQ失败: %w", err)
	}
	return ch, nil
}
