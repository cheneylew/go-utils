package utils

import (
	"time"
	"fmt"
)

func JKStringToTime(timeString string) time.Time {
	t,_ := time.Parse("2006-01-02 15:04:05", timeString)
	return t
}

func JKStringToDate(timeString string) time.Time {
	t,_ := time.Parse("2006-01-02", timeString)
	return t
}

func JKDateToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func JKTimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func JKTimeNowStamp() int64 {
	return time.Now().Unix()
}

func JKTimeNowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func JKDateNowStr() string {
	return time.Now().Format("2006-01-02")
}

// 时间格式化为 2006-01-02 15:04:05
func JKTimeFormat(stamp int64) string {
	return time.Unix(stamp, 0).Format("2006-01-02 15:04:05")
}

// 时间字符串 "2014-01-08 09:04:41" 转化为 时间戳
func JKTimeStamp(timeStr string) int64 {
	the_time, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err == nil {
		unix_time := the_time.Unix()
		return unix_time
	}else {
		return 0
	}
}

/* 代码运行时间的 Timer*/
type JKTimer struct {
	startTime        time.Time
	lastTime         time.Time
	recordCount      int64
	recordCountTotal int64
}

func (t *JKTimer)Start() {
	t.startTime = time.Now()
	t.lastTime = t.startTime
	t.recordCount = 0
	t.recordCountTotal = 0
}

func (t *JKTimer)Record() {
	t.recordCount ++
	now := time.Now()
	nano := now.UnixNano() - t.lastTime.UnixNano()
	fmt.Printf("per run time %v: %vs\n", t.recordCount, (float64(nano) / float64(1000000000.0)))
	t.lastTime = now
}

func (t *JKTimer)RecordTotal() {
	t.recordCountTotal ++
	now := time.Now()
	nano := now.UnixNano() - t.startTime.UnixNano()
	fmt.Printf("all run time %v: %vs\n", t.recordCountTotal, (float64(nano) / float64(1000000000.0)))
	t.lastTime = now
}

func (t *JKTimer)Reset() {
	t.startTime = time.Now()
	t.lastTime = time.Now()
	t.recordCount = 0
	t.recordCountTotal = 0
}