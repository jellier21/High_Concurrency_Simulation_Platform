package Redis

import (
	"client/model"
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func GetTaskMessageInRedis(client *redis.Client) ([]model.TaskDB, error) {
	ctx := context.Background()
	var allKeys []string
	// 游标，初始0
	cursor := uint64(0)
	for {
		// Scan：游标，匹配表达式，每次最多返回数量
		keys, nextCursor, err := client.Scan(ctx, cursor, "task:*", 100).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break // 游标回到0，遍历结束
		}
	}

	var list []model.TaskDB
	// 循环读取每个activity的hash数据
	for _, k := range allKeys {
		hashData, err := client.HGetAll(ctx, k).Result()
		if err != nil {
			continue
		}
		// 跳过空壳数据，避免误导回源判断
		if hashData["uuid"] == "" && hashData["status"] == "" && hashData["msl_no"] == "" && hashData["total"] == "" && hashData["success_int"] == "" && hashData["fail_int"] == "" {
			continue
		}
		Total_s, _ := strconv.ParseInt(hashData["total"], 10, 64)
		SuccessInt_s, _ := strconv.ParseInt(hashData["success_int"], 10, 64)
		FailInt_s, _ := strconv.ParseInt(hashData["fail_int"], 10, 64)
		list = append(list, model.TaskDB{Uuid: hashData["uuid"], Status: hashData["status"], MslNo: hashData["msl_no"], Total: Total_s, SuccessInt: SuccessInt_s, FailInt: FailInt_s})
	}

	return list, nil
}
