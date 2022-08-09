package util

import "time"

const TimeFormat = "2006-01-02 15:04:05"

//时间格式化
func FormatTimeToString(t uint64) (ts string) {
	ts = time.Unix(int64(t)/1e3, 0).Format(TimeFormat)
	return
}

func FromatStringToTimeStamp(s string) (stamp uint64) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return 0
	}
	return uint64(t.UnixNano() / 1e6)
}
