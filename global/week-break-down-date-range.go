package global

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DateWeekDay struct {
	Date       string
	Day        string
	WeekNumber int
	WeekStart  string
	WeekEnd    string
}

type WeekBreakDownDateRange struct {
	StartDate string
	EndDate   string

	targetYear  int
	targetMonth int

	allYearDates map[string]DateWeekDay
	WeekDates    map[int]WeekDate
}

func (obj *WeekBreakDownDateRange) Run() {
	obj.allYearDates = make(map[string]DateWeekDay)
	obj.WeekDates = make(map[int]WeekDate)

	//todo let find our target year and month
	obj.findTargetMonthYear()

	//todo generate dates
	obj.generateAllYearDates()

	//todo generate week number with hist days
	obj.generateWeekNumberDays()
}
func (obj *WeekBreakDownDateRange) findTargetMonthYear() {
	arrDate := strings.Split(obj.StartDate, "-")
	obj.targetYear, _ = strconv.Atoi(arrDate[0])
	obj.targetMonth, _ = strconv.Atoi(arrDate[1])
}
func (obj *WeekBreakDownDateRange) generateAllYearDates() {
	date := time.Date(obj.targetYear, time.Month(obj.targetMonth), 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 2000; i++ {
		var o DateWeekDay
		weekStart := WeekStartDate(date)
		o.WeekStart = strings.Split(weekStart.String(), "T")[0]
		o.Date = fmt.Sprintf("%v", date.Format("2006-01-02")) // strings.Split(date.String(),"T")[0]
		o.Day = fmt.Sprintf("%v",
			date.Format("Mon"))
		_, o.WeekNumber = date.ISOWeek()

		//let find start and end of the week date
		wStart, wEnd := WeekRange(obj.targetYear, o.WeekNumber)
		o.WeekStart = fmt.Sprintf("%v", wStart.Format("2006-01-02"))
		o.WeekEnd = fmt.Sprintf("%v", wEnd.Format("2006-01-02"))

		obj.allYearDates[o.Date] = o

		//let now increment our date to next
		date = date.Add(24 * time.Hour)

		//let check to see if your still in the current year
		if date.Year() != obj.targetYear {
			//break
		}

		nextDate := fmt.Sprintf("%v", date.Format("2006-01-02"))

		if !IsDateBetween(nextDate, obj.StartDate, obj.EndDate) {
			break
		}
	}
}
func (obj *WeekBreakDownDateRange) generateWeekNumberDays() {
	var lsWeeks []int
	var tmpWeeks = make(map[int]int)
	for _, row := range obj.allYearDates {
		tmpWeeks[row.WeekNumber] = row.WeekNumber
	}
	for num, _ := range tmpWeeks {
		lsWeeks = append(lsWeeks, num)
	}
	sort.Ints(lsWeeks)

	for _, targetWeek := range lsWeeks {
		var weekInfo WeekDate
		weekInfo.Week = targetWeek
		for _, row := range obj.allYearDates {
			if row.WeekNumber != targetWeek {
				continue
			}
			weekInfo.StartDate = row.WeekStart
			weekInfo.EndDate = row.WeekEnd

			if row.Day == "Mon" {
				weekInfo.Mon.Date = row.Date
			}
			if row.Day == "Tue" {
				weekInfo.Tue.Date = row.Date
			}
			if row.Day == "Wed" {
				weekInfo.Wed.Date = row.Date
			}
			if row.Day == "Thu" {
				weekInfo.Thu.Date = row.Date
			}
			if row.Day == "Fri" {
				weekInfo.Fri.Date = row.Date
			}
			if row.Day == "Sat" {
				weekInfo.Sat.Date = row.Date
			}
			if row.Day == "Sun" {
				weekInfo.Sun.Date = row.Date
			}
		}
		obj.WeekDates[targetWeek] = weekInfo
	}
}
