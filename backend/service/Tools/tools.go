package Tools

import (
	"github.com/google/uuid"
	"time"
)

func CreateOnlyUuid() string {
	Uid := uuid.NewString()
	return Uid
}

func CutTimeToMin(timestamp int64) string {
	t := time.UnixMilli(timestamp) // 使用不同的参数名
	// 保留到分钟
	tMinute := t.Truncate(time.Minute)
	return tMinute.Format("2006-01-02 15:04")
}

// DiffToMillisecond  ---------------------时间差计算---------------------
func DiffToMillisecond(timeStr string) (int64, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	layout := "2006-01-02 15:04"

	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return 0, err
	}

	// 截断到分钟，去掉秒和纳秒
	t = t.Truncate(time.Minute)
	now := time.Now().In(loc).Truncate(time.Minute)

	// 差值转毫秒
	return t.Sub(now).Milliseconds(), nil
}
