package Redis

import (
	"client/model"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func GetRedisOrderByUuid(uuid string, client *redis.Client) (model.OrderDB, error) {
	ctx := context.Background()

	key := fmt.Sprintf("order:%s", uuid)
	hashMap, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		return model.OrderDB{}, err
	}
	if len(hashMap) == 0 {
		return model.OrderDB{}, fmt.Errorf("order not found")
	}
	return model.OrderDB{
		Uuid:      hashMap["uuid"],
		MslNo:     hashMap["msl_no"],
		StartTime: hashMap["start_time"],
	}, nil
}
