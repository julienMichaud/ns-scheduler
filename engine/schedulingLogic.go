package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//ConvertHourToDateFormat take a scheduleUptime variable like "1-6 14:45-15:00" and return the first day (1) as int, the last day as int (6),the current date where the resources should start (2022-04-23 14:45:00 +0200 CEST) and the current date after which the resources should be stopped (2022-04-23 15:00:00 +0200 CEST)
func ConvertHourToDateFormat(scheduleUptime string, loc *time.Location) (firstDay int, lastDay int, currentDateWithStartHour time.Time, currentDateWithStopHour time.Time) {

	scheduleDaysandHoursSplitted := strings.Split(scheduleUptime, " ") // splitting 1-6 and 14:45-15:00

	scheduleDays := scheduleDaysandHoursSplitted[0]                          //creating a string 1-6
	scheduleFirstDay, _ := strconv.Atoi(strings.Split(scheduleDays, "-")[0]) //converting 1 into an integer
	scheduleLastDay, _ := strconv.Atoi(strings.Split(scheduleDays, "-")[1])  //converting 6 into an integer

	scheduleHours := scheduleDaysandHoursSplitted[1] // creating a string 14:45-15:00

	scheduleStart := strings.Split(scheduleHours, "-")[0]                        //creating a string 14:45
	scheduleStartHour, _ := strconv.Atoi(strings.Split(scheduleStart, ":")[0])   // converting 14 into an integer
	scheduleStartMinute, _ := strconv.Atoi(strings.Split(scheduleStart, ":")[1]) // converting 45 into an integer

	scheduleEnd := strings.Split(scheduleHours, "-")[1]                      //creating a string 15:00
	scheduleEndHour, _ := strconv.Atoi(strings.Split(scheduleEnd, ":")[0])   // converting 15 into an integer
	scheduleEndMinute, _ := strconv.Atoi(strings.Split(scheduleEnd, ":")[1]) // converting 00 into an integer

	StartScheduling := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), scheduleStartHour, scheduleStartMinute, 0, 0, loc) // creating a time.Time variable with current year,month,day but with hour and minutes 14:45
	EndScheduling := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), scheduleEndHour, scheduleEndMinute, 0, 0, loc)       // creating a time.Time variable with current year,month,day but with hour and minutes 15:00
	return scheduleFirstDay, scheduleLastDay, StartScheduling, EndScheduling
}

func shouldScaleDown(upTime string) bool {
	loc, _ := time.LoadLocation("Europe/Paris")
	firstDay, lastDay, currentDateWithStartHour, currentDateWithStopHour := ConvertHourToDateFormat(upTime, loc)
	now := time.Now().In(loc)
	currentWeekDay := int(now.Weekday())
	if currentWeekDay == 0 {
		currentWeekDay = 7
	}

	if (currentWeekDay >= firstDay) && (currentWeekDay <= lastDay) { //check if current date is in range specified above
		// fmt.Printf("in range\n")

		if (now.After(currentDateWithStartHour)) && (now.Before(currentDateWithStopHour)) { //check if current hour is between start and stop, if not we should scale down
			// fmt.Printf("date is between the start and stop given, resources should be up\n")
			// fmt.Println(currentWeekDay)
			fmt.Println(now)
			fmt.Println(currentDateWithStartHour)
			fmt.Println(currentDateWithStopHour)
			return false
		} else {
			// fmt.Println(currentWeekDay)
			fmt.Println(now)
			fmt.Println(currentDateWithStartHour)
			fmt.Println(currentDateWithStopHour)
			return true
		}

	} else {
		// fmt.Printf("not in range\n")
		return true
	}
}
