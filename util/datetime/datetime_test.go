package datetime_test

import (
	"fmt"
	"testing"
	"time"

	th "github.com/imdm/go-util/util/datetime"
)

func TestStartAndEnd(t *testing.T) {
	t.Log(th.StartOfToday())
	t.Log(th.EndOfToday())
	t.Log(th.StartOfThisMonth())
	t.Log(th.EndOfThisMonth())
	t.Log(th.StartOfThisYear())
	t.Log(th.EndOfThisYear())
}

func TestCurrentMonth(t *testing.T) {
	ti := th.CurrentMonth()
	fmt.Println("current month", ti)
}

func TestPrevMonth(t *testing.T) {
	ti := th.PrevMonth(th.PrevMonth(time.Now()))
	fmt.Println("prev month", ti)
}

func TestStartOfThisWeek(t *testing.T) {
	t.Log(th.StartOfThisWeek())
}

func TestEndOfThisWeek(t *testing.T) {
	t.Log(th.EndOfThisWeek())
}

func TestHoursInDayIgnoreWeekend(t *testing.T) {
	tests := []struct {
		name   string
		minute int64
		start  time.Time
		want   time.Time
	}{
		{
			"same day",
			5,
			time.Date(2020, 8, 18, 13, 14, 15, 16, time.Local),
			time.Date(2020, 8, 18, 13, 19, 15, 16, time.Local),
		},
		{
			"cross single day, both are weekdays",
			720,
			time.Date(2020, 8, 18, 13, 14, 15, 16, time.Local),
			time.Date(2020, 8, 19, 1, 14, 15, 16, time.Local),
		},
		{
			"start from weekday, end at weekday",
			2880, // 2 days
			time.Date(2020, 8, 18, 13, 14, 15, 16, time.Local), // Tue
			time.Date(2020, 8, 20, 13, 14, 15, 16, time.Local), // Thu
		},
		{
			"start from weekday, end at weekend",
			2880, // 2 days
			time.Date(2020, 8, 20, 13, 14, 15, 16, time.Local),
			time.Date(2020, 8, 24, 13, 14, 15, 16, time.Local),
		},
		{
			"start from weekday, cross weekend, end at weekday",
			10080, // 7 days
			time.Date(2020, 8, 18, 13, 14, 15, 16, time.Local),
			time.Date(2020, 8, 27, 13, 14, 15, 16, time.Local),
		},
		{
			"start from weekday, cross double weekends, end at weekday",
			20160, // 14 days
			time.Date(2020, 8, 18, 13, 14, 15, 16, time.Local),
			time.Date(2020, 9, 7, 13, 14, 15, 16, time.Local),
		},
		{
			"start from weekend, end at weekday",
			2880, // 2 days
			time.Date(2020, 8, 22, 13, 14, 15, 16, time.Local),
			time.Date(2020, 8, 25, 13, 14, 15, 16, time.Local),
		},
		{
			"start from weekend, cross weekday, end at weekend",
			10080, // 7 days
			time.Date(2020, 8, 22, 13, 14, 15, 16, time.Local),
			time.Date(2020, 9, 1, 13, 14, 15, 16, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := th.TimeAfterIgnoreWeekend(tt.start, time.Duration(tt.minute)*time.Minute)
			if !got.Equal(tt.want) {
				t.Errorf("HoursInDayIgnoreWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}
