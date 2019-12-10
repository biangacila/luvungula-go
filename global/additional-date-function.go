package global

import (
	"fmt"
	"strings"
	"time"
)

func CompareDateToNow(dateStore string) string {
	//dateStore:="2019-11-27"
	targetDate := StringToTime(dateStore, "00:01")
	state := "none"
	now := time.Now()
	if now.After(targetDate) {
		state = "after"
	}
	if now.Before(targetDate) {
		state = "before"
	}
	dt, _ := GetDateAndTimeString()
	if dt == dateStore {
		state = "equal"
	}
	return state
}
func AddToMyDate(storeDate string, unit string, val float64) string {
	targetDate := StringToTime(storeDate, "00:01")
	dur := time.Duration(val)
	end := targetDate.Add(dur * time.Hour)
	if unit == "day" {
		v := int(val)
		fmt.Println("****> ", v)
		end = targetDate.AddDate(0, 0, v)
	}
	if unit == "week" {
		v := int(val + 7)
		end = targetDate.AddDate(0, 0, v)
	}
	if unit == "month" {
		v := int(val)
		end = targetDate.AddDate(0, v, 0)
	}
	if unit == "year" {
		v := int(val)
		end = targetDate.AddDate(v, 0, 0)
	}
	arr := strings.Split(end.String(), " ")
	return arr[0]
}
