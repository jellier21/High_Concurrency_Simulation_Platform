package sqldb

type ActivityDB struct {
	MslUid    string `json:"msl_uid"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type OrderDB struct {
	Uuid      string `json:"uuid"`
	MslNo     string `json:"msl_no"`
	StartTime string `json:"start_time"`
}
