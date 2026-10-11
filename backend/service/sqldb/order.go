package sqldb

import (
	"database/sql"
	"service/Tools"
)

type OrderDB struct {
	Uuid      string `json:"uuid"`
	MslNo     string `json:"msl_no"`
	StartTime string `json:"start_time"`
}

func CreateAndInsertOrderToSQL(MslNo string, db *sql.DB) (OrderDB, error) {
	sqlSelect := "SELECT start_time FROM activity WHERE msl_no = ?"
	var startTime string
	err := db.QueryRow(sqlSelect, MslNo).Scan(&startTime)
	if err != nil {
		return OrderDB{}, err
	}
	sqlStr := "INSERT INTO `order` (uuid, msl_no, start_time) VALUES (?, ?, ?)"
	uuid := Tools.CreateOnlyUuid()
	_, err = db.Exec(sqlStr, uuid, MslNo, startTime)
	if err != nil {
		return OrderDB{}, err
	}
	return OrderDB{Uuid: uuid, MslNo: MslNo, StartTime: startTime}, nil
}
