package sqldb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // 添加这行
)

func InitMySQL() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/simulation")
	if err != nil {
		panic(err)
	}
	return db
}
