package Redis

type ActivityMsg struct {
	Uid       string `json:"msl_uid"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
