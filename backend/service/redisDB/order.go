package redisDB

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"service/sqldb"
)

func InsertOrderRedisMessage(order sqldb.OrderDB, client *redis.Client) error {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", order.Uuid)

	// HSet 存 hash 字段
	err := client.HSet(ctx, key, map[string]interface{}{
		"uid":        order.Uuid,
		"msl_no":     order.MslNo,
		"start_time": order.StartTime,
		"status":     order.Status,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetOrderRedisDeal(uuid string, client *redis.Client) error {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", uuid)
	return client.HSet(ctx, key, "status", "SUCCESS").Err()
}

func GetMslNORedisByUuidInOrder(uuid string, client *redis.Client) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", uuid)
	return client.HGet(ctx, key, "msl_no").Result()
}

func SetOrderOverdue(uuid string, client *redis.Client) error {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", uuid)
	return client.HSet(ctx, key, "status", "OVERDUE").Err()
}

func GetOrderStatus(uuid string, client *redis.Client) (string, error) {
	ctx := context.Background()
	return client.HGet(ctx, fmt.Sprintf("order:%s", uuid), "status").Result()
}
