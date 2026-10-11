package order

import (
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"service/redisDB"
	"service/sqldb"
	"sync"
)

type Job struct {
	Start int
	Count int
}

const (
	workPoolNum = 5
	batchSize   = 100
)

func CreateOrderBatch(Uuid, mslNo string, job Job, db *sql.DB, red *redis.Client) error {
	// TODO：在这里实现当前批次的订单生成和 MySQL 写入
	for i := job.Start; i < job.Start+job.Count; i++ {
		var orders sqldb.OrderDB
		orders, err := sqldb.CreateAndInsertOrderToSQL(mslNo, db)
		if err != nil {
			log.Printf("此条订单生成失败: %v", err)
			return err
		}
		Err := redisDB.InsertOrderRedisMessage(orders, red)
		if Err != nil {
			log.Printf("此条订单插入失败: %v", Err)
			continue
		}
	}

	fmt.Printf(
		"处理订单编号范围：%d ~ %d\n",
		job.Start,
		job.Start+job.Count-1,
	)

	err := redisDB.UpdateTaskRedisMessage(Uuid, red, db)
	if err != nil {
		log.Println("更新Redis任务状态失败")
	}
	return nil
}

func Worker(mslNo, Uuid string, id int, jobs <-chan Job, wg *sync.WaitGroup, db *sql.DB, red *redis.Client) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d 开始处理批次\n", id)

		if err := CreateOrderBatch(Uuid, mslNo, job, db, red); err != nil {
			fmt.Printf("Worker %d 处理失败：%v\n", id, err)

			err := sqldb.UpdataTaskForOneFail(Uuid, db)
			if err != nil {
				fmt.Printf("Worker %d 更新status为fail任务失败：%v\n", id, err)
			}
			continue
		}

		err := sqldb.UpdataTaskForOneSuccess(Uuid, db)
		if err != nil {
			fmt.Printf("Worker %d 更新status为success任务失败：%v\n", id, err)
		}

		fmt.Printf("Worker %d 处理完成\n", id)
	}
}

func StartOrderTasking(mslNo, Uuid string, total int, db *sql.DB, red *redis.Client) {
	jobs := make(chan Job, workPoolNum)
	var wg sync.WaitGroup

	for i := 0; i < workPoolNum; i++ {
		wg.Add(1)
		go Worker(mslNo, Uuid, i, jobs, &wg, db, red)
	}

	// 按批次生成任务
	for start := 0; start < total; start += batchSize {
		count := batchSize

		// 处理最后一个不足 batchSize 的批次
		if start+count > total {
			count = total - start
		}

		jobs <- Job{
			Start: start,
			Count: count,
		}
	}
	// 告诉 Worker：不会再有新任务
	close(jobs)

	// 等待全部 Worker 完成
	wg.Wait()

	_ = sqldb.UpdateAllTaskSuccess(Uuid, db)

	// 所有批次处理完，MySQL 已是最新，同步最终状态到 Redis
	if uerr := redisDB.UpdateTaskRedisMessage(Uuid, red, db); uerr != nil {
		log.Printf("更新Redis任务状态失败: %v", uerr)
	}
	fmt.Println("全部批次处理完成")
}
