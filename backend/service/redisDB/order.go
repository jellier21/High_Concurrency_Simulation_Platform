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
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
