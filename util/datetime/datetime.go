package datetime

import (
	"strings"
	"time"
)

const (
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

func FormatDate(date time.Time) string {
	return date.Format(YYYY_MM_DD)
}

func FormarTime(date time.Time) string {
	return date.Format(YYYY_MM_DD_HH_MM_SS)
}

func ToDate(str string) (time.Time, error) {
	return time.ParseInLocation(YYYY_MM_DD, str, time.Local)
}

func ToDateTime(str string) time.Time {
	t, _ := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, str, time.Local)
	return t
}

func ToTimestamp(str string) int64 {
	now, _ := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, str, time.Local)
	return now.Unix()
}

func ToDayTimestamp(str string) int64 {
	date := strings.Split(str, " ")
	now, _ := time.ParseInLocation(YYYY_MM_DD, date[0], time.Local)
	return now.Unix()
}

func UtilDayEnd() int64 {
	remainSecond := EndOfToday().Unix() - time.Now().Local().Unix()
	return remainSecond
}

// Beyond2DayUtilNow 时间距离现在超过2天,即中间至少断了1天
func Beyond2DayUtilNow(t time.Time) bool {
	return time.Now().Local().Add(-48 * time.Hour).After(t)
}

// In1Day 时间在当天
func In1Day(t time.Time) bool {
	if t.After(time.Now()) {
		return false
	}
	return EndOfToday().Unix()-t.Unix() < 24*3600
}

func EndOfToday() time.Time {
	return EndOfDate(time.Now())
}

func EndOfDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
}

func EndOfThisWeek() time.Time {
	return EndOfWeek(time.Now())
}

func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t.AddDate(0, 0, 7)).Add(time.Second * -1)
}

func EndOfThisMonth() time.Time {
	return EndOfMonth(time.Now())
}

func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t.AddDate(0, 1, 0)).Add(time.Second * -1)
}

func EndOfThisYear() time.Time {
	return EndOfYear(time.Now())
}

func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t.AddDate(1, 0, 0)).Add(time.Second * -1)
}

func StartOfToday() time.Time {
	return StartOfDate(time.Now())
}

func StartOfDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func StartOfThisWeek() time.Time {
	return StartOfWeek(time.Now())
}

func StartOfWeek(t time.Time) time.Time {
	var (
		add int
		w   = t.Weekday()
	)
	switch w {
	case time.Sunday:
		add = -6
	default:
		add = 1 - int(w)
	}
	return StartOfDate(t.AddDate(0, 0, add))
}

func StartOfThisMonth() time.Time {
	return StartOfMonth(time.Now())
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func StartOfThisYear() time.Time {
	return StartOfYear(time.Now())
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
}

func NDaysAfter(n int) time.Time {
	return EndOfToday().AddDate(0, 0, n)
}

// 当前月一号
// Deprecated: use datetime.StartOfThisMonth instead
func CurrentMonth() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, time.Local)
}

// 上个月
// Deprecated: use datetime.StartOfMonth instead
func PrevMonth(n time.Time) time.Time {
	m := n.Month()
	y := n.Year()
	m = m - 1
	if m == 0 {
		y = y - 1
		m = 12
	}
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

func TimeAfterIgnoreWeekend(start time.Time, delta time.Duration) time.Time {
	end := EndOfDate(start)
	curDayLeft := end.Sub(start)
	if delta <= curDayLeft { // same day
		return start.Add(delta)
	}

	res := start
	res = JumpWeekend(res)
	for {
		res = JumpWeekend(res)
		if delta >= 0 && delta <= time.Hour*24 {
			return res.Add(delta)
		}
		delta -= time.Minute * 60 * 24
		res = res.Add(time.Hour * 24)
	}
	return res
}

func JumpWeekend(r time.Time) time.Time {
	if r.Weekday() == time.Friday {
		r = r.Add(time.Hour * 48)
	} else if r.Weekday() == time.Saturday {
		r = r.Add(time.Hour * 24)
	} else if r.Weekday() == time.Sunday {
		r = r.Add(time.Hour * 0)
	}
	return r
}
