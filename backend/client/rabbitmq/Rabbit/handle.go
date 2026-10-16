package Rabbit

import (
	"client/model"
	"client/rabbitmq/MqTools"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func SendOrderToQueueAndProcessOrder(data model.OrderDB, ch *amqp091.Channel) error {
	// =======·实现将订单发送到队列并处理订单的逻辑·=========
	// 开启发布确认。下面一次只发布一条并等待确认。
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("启动发布前确认失败: %w", err)
	}
	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := MqTools.PublishConfirmed(
		ch, confirms, "order_exchange", "order", body,
	); err != nil {
		return err
	}

	log.Printf("订单消息已获 RabbitMQ 发布确认: %+v", data)
	return nil
}
