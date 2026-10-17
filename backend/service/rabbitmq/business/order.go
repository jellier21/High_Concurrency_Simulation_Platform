package business

import (
	"database/sql"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"service/redisDB"
	"service/sqldb"
	"time"
)

type OrderMessage struct {
	Data       sqldb.OrderDB `json:"data"`
	RetryCount int           `json:"retry_count"`
}

const maxRetries = 3

func handleOrderMessage(
	pubCh *amqp.Channel,
	confirms <-chan amqp.Confirmation,
	msg amqp.Delivery,
	db *sql.DB,
	client *redis.Client,

) error {
	var order OrderMessage

	// 如果连 JSON 都解析不了，也送入死信队列，避免反复消费。
	if err := json.Unmarshal(msg.Body, &order); err != nil {
		log.Printf("消息格式错误，转入死信队列: %v", err)

		if err := publishConfirmed(
			pubCh, confirms, "dead_exchange", "dead", msg.Body,
		); err != nil {
			return err
		}
		return msg.Ack(false)
	}

	// 业务成功：确认原消息。
	if err := processOrderV2(order, db, client); err == nil {
		return msg.Ack(false)
	} else {
		log.Printf("订单 %s 处理失败: %v", order.Data.Uuid, err)
	}

	var exchange, routingKey string

	if order.RetryCount < maxRetries {
		order.RetryCount++
		exchange, routingKey = "retry_exchange", "retry"
	} else {
		exchange, routingKey = "dead_exchange", "dead"
	}

	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// 关键顺序：先确认新消息已被 RabbitMQ 接收，再 ACK 原消息。
	if err := publishConfirmed(
		pubCh, confirms, exchange, routingKey, body,
	); err != nil {
		// 不 ACK 原消息。消费者退出、连接关闭后，
		// RabbitMQ 可以重新投递尚未确认的消息。
		return fmt.Errorf("转发失败，原消息尚未 ACK: %w", err)
	}

	if exchange == "retry_exchange" {
		log.Printf(
			"订单 %s 已进入延迟重试队列，第 %d 次重试将在约 5 秒后进行",
			order.Data.Uuid, order.RetryCount,
		)
	} else {
		log.Printf(
			"订单 %s 已达到重试上限，转入死信队列",
			order.Data.Uuid,
		)
	}

	return msg.Ack(false)
}

func processOrder(order OrderMessage, db *sql.DB, client *redis.Client) error {
	current, err := redisDB.GetOrderStatus(order.Data.Uuid, client)
	if err == nil && current != "" && current != "WAITING" {
		log.Printf("订单 %s 已是 %s 状态，跳过重复消费", order.Data.Uuid, current)
		return nil
	}
	err = CheckWeatherOverdue(order.Data.Uuid, db, client)
	if err != nil {
		Err := SetOrderStatusOverdue(order.Data.Uuid, db, client)
		if Err != nil {
			log.Printf("订单 %s 更新状态失败: %v", order.Data.Uuid, Err)
			return Err
		}
		log.Printf("订单 %s 已超过活动结束时间，标记 OVERDUE", order.Data.Uuid)
		return nil
	}

	err = sqldb.SetOrderStatusDeal(order.Data.Uuid, db)
	if err != nil {
		return err
	}
	Err := redisDB.SetOrderRedisDeal(order.Data.Uuid, client)
	if Err != nil {
		log.Printf("订单 %s 存入 Redis 失败: %v, 此错误不影响业务处理", order.Data.Uuid, Err)
	}

	//log.Printf("订单 %s 处理成功", order.Data.Uuid)
	return nil
}

func CheckWeatherOverdue(uuid string, db *sql.DB, client *redis.Client) error {
	MslNo, err := redisDB.GetMslNORedisByUuidInOrder(uuid, client)
	if err != nil {
		log.Printf("订单 %s 查询Redis的 MslNo 失败: %v， 正在从MySQL查询", uuid, err)
		MslNo, err = sqldb.GetMslNObyUuidInOrder(uuid, db)
		if err != nil {
			log.Printf("订单 %s 查询MySQL的 MslNo 失败: %v", uuid, err)
			return err
		}
	}
	endTime, err := redisDB.GetActivityRedisLastestTime(MslNo, client)
	if err != nil {
		log.Printf("订单 %s 查询Redis的活动时间失败: %v， 正在从MySQL查询", uuid, err)
		endTime, err = sqldb.GetActivitySQLLastestTime(MslNo, db)
		if err != nil {
			log.Printf("订单 %s 查询MySQL的活动时间失败: %v", uuid, err)
			return err
		}
	}

	layout := "2006-01-02 15:04:05"
	targetTime, Error := time.ParseInLocation(layout, endTime, time.Local)
	if Error != nil {
		// 前端传的时间不带秒，兜底再试一次无秒格式
		targetTime, Error = time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
		if Error != nil {
			log.Printf("订单 %s 解析时间失败: %v, 原始值: %q", uuid, Error, endTime)
			return Error
		}
	}

	now := time.Now()
	if now.After(targetTime) {
		log.Printf("订单 %s 超时", uuid)
		return fmt.Errorf("订单 %s 超时", uuid)
	}
	return nil
}

func SetOrderStatusOverdue(uuid string, db *sql.DB, client *redis.Client) error {
	err := sqldb.SetSQLOrderStatusOverdue(uuid, db)
	if err != nil {
		log.Printf("订单 %s 更新数据库状态失败: %v", uuid, err)
		return err
	}
	Err := redisDB.SetOrderOverdue(uuid, client)
	if Err != nil {
		log.Printf("订单 %s 更新 Redis 状态失败: %v", uuid, Err)
	}
	return nil
}

// ==========================------------------第二版处理订单逻辑---------------------===========================
func processOrderV2(order OrderMessage, db *sql.DB, client *redis.Client) error {
	// 一次 HGetAll 拿全：status + msl_no，替代原来的两次 HGET
	data, err := redisDB.GetOrderAllByUuid(order.Data.Uuid, client)
	if err != nil {
		log.Printf("订单 %s 查询 Redis 失败: %v", order.Data.Uuid, err)
		return err
	}

	status := data["status"]
	if status != "" && status != "WAITING" {
		log.Printf("订单 %s 已是 %s，跳过重复消费", order.Data.Uuid, status)
		return nil
	}

	mslNo := data["msl_no"]
	if err := checkOverdue(mslNo, db, client); err != nil {
		GlobalBatchWriter.Add(order.Data.Uuid, "OVERDUE")
		log.Printf("订单 %s 已超过活动结束时间，标记 OVERDUE", order.Data.Uuid)
		return nil
	}

	GlobalBatchWriter.Add(order.Data.Uuid, "SUCCESS")
	if err := redisDB.SetOrderRedisDeal(order.Data.Uuid, client); err != nil {
		log.Printf("订单 %s 更新 Redis 失败: %v", order.Data.Uuid, err)
	}
	return nil
}

func checkOverdue(mslNo string, db *sql.DB, client *redis.Client) error {
	endTime, err := redisDB.GetActivityRedisLastestTime(mslNo, client)
	if err != nil {
		endTime, err = sqldb.GetActivitySQLLastestTime(mslNo, db)
		if err != nil {
			return err
		}
	}

	targetTime, Error := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if Error != nil {
		targetTime, Error = time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
		if Error != nil {
			return Error
		}
	}
	if time.Now().After(targetTime) {
		return fmt.Errorf("订单下单时活动已结束")
	}
	return nil
}
