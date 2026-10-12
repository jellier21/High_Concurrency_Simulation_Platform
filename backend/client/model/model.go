package model

type TaskDB struct {
	Uuid       string `json:"uuid"`
	MslNo      string `json:"msl_no"`
	Status     string `json:"status"`
	Total      int64  `json:"total"`
	SuccessInt int64  `json:"success_int"`
	FailInt    int64  `json:"fail_int"`
}

type ActivityMsg struct {
	Uid       string `json:"msl_uid"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type OrderDB struct {
	Uuid      string `json:"uuid"`
	MslNo     string `json:"msl_no"`
	StartTime string `json:"start_time"`
	Status    string `json:"status"`
}
