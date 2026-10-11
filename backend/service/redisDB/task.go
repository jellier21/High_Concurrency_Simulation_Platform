package redisDB

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"service/sqldb"
)

func InsertTaskRedisMessage(task sqldb.TaskDB, client *redis.Client) error {
	ctx := context.Background()
	key := fmt.Sprintf("task:%s", task.Uuid)

	// HSet 存 hash 字段
	err := client.HSet(ctx, key, map[string]interface{}{
		"uid":         task.Uuid,
		"status":      task.Status,
		"total":       task.Total,
		"success_int": task.SuccessInt,
		"fail_int":    task.FailInt,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskRedisMessage(Uuid string, client *redis.Client, db *sql.DB) error {
	sqlStr := `SELECT uuid,status,total,success_int,fail_int FROM task WHERE uuid = ?`
	row := db.QueryRow(sqlStr, Uuid)
	var successInt, failInt int64
	var Status string
	var Total int
	err := row.Scan(&Uuid, &Status, &Total, &successInt, &failInt)
	if err != nil {
		return err
	}
	ctx := context.Background()
	key := fmt.Sprintf("task:%s", Uuid)
	return client.HSet(ctx, key, map[string]interface{}{
		"status":      Status,
		"total":       Total,
		"success_int": successInt,
		"fail_int":    failInt,
	}).Err()
}
