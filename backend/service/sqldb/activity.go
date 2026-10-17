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

func GetActivitySQLLastestTime(msl_no string, db *sql.DB) (string, error) {
	sqlStr := `SELECT end_time FROM activity WHERE msl_no = ? ORDER BY end_time DESC LIMIT 1`
	var latestTime string
	err := db.QueryRow(sqlStr, msl_no).Scan(&latestTime)
	if err != nil {
		return "", err
	}
	return latestTime, nil
}
