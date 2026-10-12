package sqldb

import (
	"database/sql"
	"service/Tools"
)

func CreateAndInsertOrderToSQL(MslNo string, db *sql.DB) (OrderDB, error) {
	sqlSelect := "SELECT start_time FROM activity WHERE msl_no = ?"
	var startTime string
	err := db.QueryRow(sqlSelect, MslNo).Scan(&startTime)
	if err != nil {
		return OrderDB{}, err
	}
	sqlStr := "INSERT INTO `order` (uuid, msl_no, start_time, status) VALUES (?, ?, ?, ?)"
	uuid := Tools.CreateOnlyUuid()
	_, err = db.Exec(sqlStr, uuid, MslNo, startTime, "WAITING")
	if err != nil {
		return OrderDB{}, err
	}
	return OrderDB{Uuid: uuid, MslNo: MslNo, StartTime: startTime, Status: "WAITING"}, nil

}

func GetOrderMessageInSQL(Uuid string, db *sql.DB) (OrderDB, error) {
	var orderData OrderDB
	sqlStr := "SELECT uuid, msl_no, start_time, status FROM `order` WHERE uuid = ?"
	err := db.QueryRow(sqlStr, Uuid).Scan(&orderData.Uuid, &orderData.MslNo, &orderData.StartTime, &orderData.Status)
	if err != nil {
		return OrderDB{}, err
	}
	return orderData, nil
}
