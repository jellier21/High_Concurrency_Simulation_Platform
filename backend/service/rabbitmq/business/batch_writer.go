package business

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

var GlobalBatchWriter *BatchWriter

func InitBatchWriter(db *sql.DB) {
	GlobalBatchWriter = NewBatchWriter(db)
}

type BatchWriter struct {
	mu      sync.Mutex
	success []string
	overdue []string
	db      *sql.DB
}

func NewBatchWriter(db *sql.DB) *BatchWriter {
	w := &BatchWriter{db: db}
	// 定时 flush：每 100ms 把攒的批次写出去
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for range ticker.C {
			w.Flush()
		}
	}()
	return w
}

func (w *BatchWriter) Add(uuid, status string) {
	w.mu.Lock()
	if status == "SUCCESS" {
		w.success = append(w.success, uuid)
	} else {
		w.overdue = append(w.overdue, uuid)
	}
	n := len(w.success) + len(w.overdue)
	w.mu.Unlock()
	if n >= 100 {
		w.Flush()
	}
}

func (w *BatchWriter) Flush() {
	w.mu.Lock()
	successBatch := w.success
	overdueBatch := w.overdue
	w.success = nil
	w.overdue = nil
	w.mu.Unlock()

	if len(successBatch) > 0 {
		placeholders := strings.Repeat("?,", len(successBatch))
		placeholders = placeholders[:len(placeholders)-1]
		sqlStr := fmt.Sprintf("UPDATE `order` SET status='SUCCESS' WHERE uuid IN (%s)", placeholders)
		args := make([]interface{}, len(successBatch))
		for i, u := range successBatch {
			args[i] = u
		}
		if _, err := w.db.Exec(sqlStr, args...); err != nil {
			fmt.Printf("批量更新 SUCCESS 失败: %v\n", err)
		}
	}
	if len(overdueBatch) > 0 {
		placeholders := strings.Repeat("?,", len(overdueBatch))
		placeholders = placeholders[:len(placeholders)-1]
		sqlStr := fmt.Sprintf("UPDATE `order` SET status='OVERDUE' WHERE uuid IN (%s)", placeholders)
		args := make([]interface{}, len(overdueBatch))
		for i, u := range overdueBatch {
			args[i] = u
		}
		if _, err := w.db.Exec(sqlStr, args...); err != nil {
			fmt.Printf("批量更新 OVERDUE 失败: %v\n", err)
		}
	}
}
