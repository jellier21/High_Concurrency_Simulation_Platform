package sqldb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func InitMySQL() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/simulation")
	if err != nil {
		panic(err)
	}

	// ★ 关键：连接池配置
	db.SetMaxOpenConns(50)                 // 最多 50 个并发连接
	db.SetMaxIdleConns(10)                 // 保持 10 个空闲连接随时复用
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最多活 5 分钟，定期换新

	return db
}
