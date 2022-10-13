package global

import (
	"testing"
)

func TestWeekBreakDownDateRange_Run(t *testing.T) {
	hub := WeekBreakDownDateRange{
		StartDate: "2021-07-01",
		EndDate:   "2021-07-30",
	}
	hub.Run()

	DisplayObject("AllDates", hub.allYearDates)
	DisplayObject("WeekDates", hub.WeekDates)
}
