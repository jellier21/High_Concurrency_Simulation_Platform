package Redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func GetAllActivity(client *redis.Client) ([]ActivityMsg, error) {
	ctx := context.Background()
	var allKeys []string
	// 游标，初始0
	cursor := uint64(0)
	for {
		// Scan：游标，匹配表达式，每次最多返回数量
		keys, nextCursor, err := client.Scan(ctx, cursor, "activity:*", 100).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break // 游标回到0，遍历结束
		}
	}

	var list []ActivityMsg
	// 循环读取每个activity的hash数据
	for _, k := range allKeys {
		hashData, err := client.HGetAll(ctx, k).Result()
		if err != nil {
			continue
		}
		// 跳过空壳数据，避免误导回源判断
		if hashData["uid"] == "" && hashData["start_time"] == "" && hashData["end_time"] == "" {
			continue
		}
		list = append(list, ActivityMsg{Uid: hashData["uid"], StartTime: hashData["start_time"], EndTime: hashData["end_time"]})
	}

	return list, nil
}

// 全量刷新活动缓存：先清空旧的 activity:*，再按数据库结果写回
func RefreshActivity(list []ActivityMsg, client *redis.Client) error {
	ctx := context.Background()
	// 1. 删除所有旧活动 key
	cursor := uint64(0)
	for {
		keys, next, err := client.Scan(ctx, cursor, "activity:*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	// 2. 全量写回
	for _, item := range list {
		key := fmt.Sprintf("activity:%s", item.Uid)
		if err := client.HSet(ctx, key, map[string]interface{}{
			"uid":        item.Uid,
			"start_time": item.StartTime,
			"end_time":   item.EndTime,
		}).Err(); err != nil {
			return err
		}
	}
	return nil
}
