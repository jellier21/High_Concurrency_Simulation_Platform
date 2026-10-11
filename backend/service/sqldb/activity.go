package sqldb

import (
	"database/sql"
)

func CreateActivityDB(uid, start, end string, db *sql.DB) error {
	sqlStr := `INSERT INTO activity (msl_no, start_time, end_time) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStr, uid, start, end)
	if err != nil {
		return err
	}
	return nil
}

func GetActivityDB(db *sql.DB) ([]ActivityDB, error) {
	sqlStr := `SELECT msl_no, start_time, end_time FROM activity`
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []ActivityDB
	for rows.Next() {
		var activity ActivityDB
		err := rows.Scan(&activity.MslUid, &activity.StartTime, &activity.EndTime)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
