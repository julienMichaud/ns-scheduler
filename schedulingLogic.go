package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//ConvertHourToDateFormat take a scheduleUptime variable like "1-6 14:45-15:00" and return the first day (1) as int, the last day as int (6),the current date where the resources should start (2022-04-23 14:45:00 +0200 CEST) and the current date after which the resources should be stopped (2022-04-23 15:00:00 +0200 CEST)
func ConvertHourToDateFormat(scheduleUptime string, timeNow time.Time, loc *time.Location) (firstDay int, lastDay int, currentDateWithStartHour time.Time, currentDateWithStopHour time.Time, err error) {

	scheduleDaysandHoursSplitted := strings.Split(scheduleUptime, " ") // splitting 1-6 and 14:45-15:00

	scheduleDays := scheduleDaysandHoursSplitted[0] //creating a string 1-6
	if len(scheduleDays) != 3 {
		return 0, 0, timeNow, timeNow, fmt.Errorf("days not in good format, want 3 got %v", len(scheduleDays))
	}

	scheduleFirstDay, err := strconv.Atoi(strings.Split(scheduleDays, "-")[0]) //converting 1 into an integer
	if err != nil {
		return 0, 0, timeNow, timeNow, fmt.Errorf("could not convert first day from string to integer: %v", err)
	}
	if (scheduleFirstDay >= 1) && (scheduleFirstDay <= 7) {
	} else {
		return 0, 0, timeNow, timeNow, fmt.Errorf("first day not between 1 and 7, monday is 1 and sunday is 7, got %d", scheduleFirstDay)
	}

	scheduleLastDay, err := strconv.Atoi(strings.Split(scheduleDays, "-")[1]) //converting 6 into an integer
	if err != nil {
		return 0, 0, timeNow, timeNow, fmt.Errorf("could not convert last day from string to integer: %v", err)
	}
	if (scheduleLastDay >= 1) && (scheduleLastDay <= 7) {
	} else {
		return 0, 0, timeNow, timeNow, fmt.Errorf("last day not between 1 and 7, monday is 1 and sunday is 7, got %d", scheduleLastDay)
	}

	if scheduleFirstDay > scheduleLastDay {
		return 0, 0, timeNow, timeNow, fmt.Errorf("first day cannot be higher than last day, got %d but lastday is %d", scheduleFirstDay, scheduleLastDay)
	}
	scheduleHours := scheduleDaysandHoursSplitted[1] // creating a string 14:45-15:00
	if len(scheduleHours) != 11 {
		return 0, 0, timeNow, timeNow, fmt.Errorf("scheduleHours hour not in good format, want 11 got %v", len(scheduleHours))
	}

	scheduleStart := strings.Split(scheduleHours, "-")[0] //creating a string 14:45
	if len(scheduleStart) != 5 {
		return 0, 0, timeNow, timeNow, fmt.Errorf("scheduleStart hour not in good format, want 5 got %v", len(scheduleStart))
	}

	scheduleStartHour, err := strconv.Atoi(strings.Split(scheduleStart, ":")[0]) // converting 14 into an integer
	if err != nil {
		return 0, 0, time.Now(), time.Now(), fmt.Errorf("could not convert start hour day from string to integer: %v", err)
	}

	scheduleStartMinute, err := strconv.Atoi(strings.Split(scheduleStart, ":")[1]) // converting 45 into an integer
	if err != nil {
		return 0, 0, time.Now(), time.Now(), fmt.Errorf("could not convert start minute day from string to integer: %v", err)
	}

	scheduleEnd := strings.Split(scheduleHours, "-")[1] //creating a string 15:00
	if len(scheduleEnd) != 5 {
		return 0, 0, timeNow, timeNow, fmt.Errorf("scheduleEnd hour not in good format, want 5 got %v", len(scheduleStart))
	}

	scheduleEndHour, err := strconv.Atoi(strings.Split(scheduleEnd, ":")[0]) // converting 15 into an integer
	if err != nil {
		return 0, 0, time.Now(), time.Now(), fmt.Errorf("could not convert last day hour from string to integer: %v", err)
	}

	scheduleEndMinute, err := strconv.Atoi(strings.Split(scheduleEnd, ":")[1]) // converting 00 into an integer
	if err != nil {
		return 0, 0, time.Now(), time.Now(), fmt.Errorf("could not convert last day minute from string to integer: %v", err)
	}

	StartScheduling := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), scheduleStartHour, scheduleStartMinute, 0, 0, loc) // creating a time.Time variable with current year,month,day but with hour and minutes 14:45
	EndScheduling := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), scheduleEndHour, scheduleEndMinute, 0, 0, loc)       // creating a time.Time variable with current year,month,day but with hour and minutes 15:00
	return scheduleFirstDay, scheduleLastDay, StartScheduling, EndScheduling, nil
}

func shouldScaleDown(upTime string, logger *logrus.Logger, nowTime time.Time, loc *time.Location) (bool, error) {

	firstDay, lastDay, currentDateWithStartHour, currentDateWithStopHour, err := ConvertHourToDateFormat(upTime, nowTime, loc)
	if err != nil {
		return false, fmt.Errorf("could not convert Hour to Date format: %v", err)
	}
	currentWeekDay := int(nowTime.Weekday())

	if currentWeekDay == 0 {
		currentWeekDay = 7
	}
	contextLogger := logger.WithFields(log.Fields{
		"currentWeekDay":           currentWeekDay,
		"currentTime":              nowTime,
		"currentDateWithStartHour": currentDateWithStartHour,
		"currentDateWithStopHour":  currentDateWithStopHour,
	})

	if (currentWeekDay >= firstDay) && (currentWeekDay <= lastDay) { //check if current date is in range specified above
		// fmt.Printf("in range\n")
		contextLogger.Debug("resources in between days given")

		if (nowTime.After(currentDateWithStartHour)) && (nowTime.Before(currentDateWithStopHour)) { //check if current hour is between start and stop, if not we should scale down
			contextLogger.Debug("date is between the start and stop given, resources should be up")
			return false, nil
		} else {
			contextLogger.Debug("date is NOT between the start and stop given, resources should be down")
			return true, nil
		}

	} else {
		contextLogger.Debug("resources not in between days given")
		return true, nil
	}
}
