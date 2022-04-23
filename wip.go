package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//ConvertHourToDateFormat take a scheduleUptime variable like "1-6 14:45-15:00" and return the first day (1) as int, the last day as int (6),the current date where the resources should start (2022-04-23 14:45:00 +0200 CEST) and the current date after which the resources should be stopped (2022-04-23 15:00:00 +0200 CEST)
func ConvertHourToDateFormat(scheduleUptime string) (firstDay int, lastDay int, currentDateWithStartHour time.Time, currentDateWithStopHour time.Time) {

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

	StartScheduling := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), scheduleStartHour, scheduleStartMinute, 0, 0, time.Local) // creating a time.Time variable with current year,month,day but with hour and minutes 14:45
	EndScheduling := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), scheduleEndHour, scheduleEndMinute, 0, 0, time.Local)       // creating a time.Time variable with current year,month,day but with hour and minutes 15:00
	return scheduleFirstDay, scheduleLastDay, StartScheduling, EndScheduling
}

func main() {
	scheduleUptime := "1-6 14:45-15:30"

	now := time.Now()
	currentWeekDay := int(now.Weekday())

	startDay, EndDay, startScheduling, stopScheduling := ConvertHourToDateFormat(scheduleUptime)
	fmt.Println(startScheduling)

	if (currentWeekDay >= startDay) && (currentWeekDay <= EndDay) { //check if current date is in range specified above
		fmt.Printf("in range\n")

		if (now.After(startScheduling)) && (now.Before(stopScheduling)) { //check if current hour is between start and stop, if not we should scale down
			fmt.Printf("hour is between the start and stop given, resources should be up\n")
		} else {
			fmt.Printf("not in scheduling range, we should scale down")
		}

	} else {
		fmt.Printf("not in range")
	}

}
