package sqldb

import (
	"database/sql"
	"service/Tools"
)

type TaskDB struct {
	ID         int
	Uuid       string `json:"uuid"`
	MslNo      string `json:"msl_no"`
	Status     string `json:"status"` //Success|Interrupt|Running
	Total      int64  `json:"total"`
	SuccessInt int64  `json:"success_int"`
	FailInt    int64  `json:"fail_int"`
}

func CreateTaskDB(mslNo string, num int, db *sql.DB) (TaskDB, error) {
	uuid := Tools.CreateOnlyUuid()
	sqlStr := `INSERT INTO task (uuid, msl_no, status, total, success_int, fail_int) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, uuid, mslNo, "Running", num, 0, 0)
	if err != nil {
		return TaskDB{}, err
	}
	data := TaskDB{
		Uuid:       uuid,
		MslNo:      mslNo,
		Status:     "Running",
		Total:      int64(num),
		SuccessInt: 0,
		FailInt:    0,
	}
	return data, nil
}

func GetTaskDB(db *sql.DB) ([]TaskDB, error) {
	sqlStr := `SELECT id, uuid, msl_no, status, total, success_int, fail_int FROM task`
	rows, err := db.Query(sqlStr)
	if err != nil {
		return []TaskDB{}, err
	}
	defer rows.Close()
	var data []TaskDB
	for rows.Next() {
		var item TaskDB
		err := rows.Scan(&item.ID, &item.Uuid, &item.MslNo, &item.Status, &item.Total, &item.SuccessInt, &item.FailInt)
		if err != nil {
			return []TaskDB{}, err
		}
		data = append(data, item)
	}
	return data, nil
}

func UpdataTaskForOneSuccess(uuid string, db *sql.DB) error {
	sqlStr := `UPDATE task SET success_int = success_int + 100 WHERE uuid = ?`
	_, err := db.Exec(sqlStr, uuid)
	return err
}

func UpdataTaskForOneFail(uuid string, db *sql.DB) error {
	sqlStr := `UPDATE task SET fail_int = fail_int + 100 WHERE uuid = ?`
	_, err := db.Exec(sqlStr, uuid)
	return err
}

func UpdateAllTaskSuccess(uuid string, db *sql.DB) error {
	sqlStr := `UPDATE task SET status = 'Success' WHERE uuid = ?`
	_, err := db.Exec(sqlStr, uuid)
	return err
}
