package redisDB

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InsertActivityRedisMessage(uuid, start, end string, client *redis.Client) error {
	//overdue, _ := Tools.DiffToMillisecond(end)
	ctx := context.Background()
	key := fmt.Sprintf("activity:%s", uuid)

	// HSet 存 hash 字段
	err := client.HSet(ctx, key, map[string]interface{}{
		"uid":        uuid,
		"start_time": start,
		"end_time":   end,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetActivityRedisLastestTime(msl_no string, client *redis.Client) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("activity:%s", msl_no)
	return client.HGet(ctx, key, "end_time").Result()
}
