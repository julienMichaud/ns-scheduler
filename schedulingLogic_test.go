package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestShouldScaleDown(t *testing.T) {
	var log = logrus.New()
	loc, _ := time.LoadLocation("Europe/Paris")

	shouldScaleDownTestCases := []struct {
		upTime      string
		currentTime time.Time
		want        bool
	}{
		{
			upTime:      "1-7 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc), // 2022-10-20  10:30
			want:        true,
		},
		{
			upTime:      "1-7 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 2, 50, 0, 0, loc), // 2022-10-20  02:50
			want:        true,
		},
		{
			upTime:      "1-7 11:00-13:00",
			currentTime: time.Date(2022, 11, 30, 11, 50, 0, 0, loc), // 2022-10-20  11:50
			want:        false,
		},
		{
			upTime:      "1-7 11:00-13:00",
			currentTime: time.Date(2022, 12, 50, 12, 24, 0, 0, loc), // 2022-10-20  12:24
			want:        false,
		},
		{
			upTime:      "1-7 11:00-13:00",
			currentTime: time.Date(2022, 12, 50, 12, 24, 0, 0, loc), // 2022-10-20  12:24
			want:        false,
		},
	}

	for _, tt := range shouldScaleDownTestCases {
		got, _ := shouldScaleDown(tt.upTime, log, tt.currentTime, loc)
		if got != tt.want {
			t.Errorf("got %t want %t", got, tt.want)
		}

	}
}

func TestShouldScaleDownError(t *testing.T) {
	var log = logrus.New()
	loc, _ := time.LoadLocation("Europe/Paris")

	shouldScaleDownTestCases := []struct {
		upTime      string
		currentTime time.Time
		want        string
	}{
		{
			upTime:      "a-7 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert first day from string to integer: strconv.Atoi: parsing \"a\": invalid syntax",
		},

		{
			upTime:      "7 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: days not in good format, want 3 got 1",
		},

		{
			upTime:      "12-7 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: days not in good format, want 3 got 4",
		},

		{
			upTime:      "2-14 11:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: days not in good format, want 3 got 4",
		},

		{
			upTime:      "1-5 11:0-13:000",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: scheduleStart hour not in good format, want 5 got 4",
		},

		{
			upTime:      "1-5 9:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: scheduleHours hour not in good format, want 11 got 10",
		},
		{
			upTime:      "1-5 1@:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert start hour day from string to integer: strconv.Atoi: parsing \"1@\": invalid syntax",
		},

		{
			upTime:      "1-5 aa:00-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert start hour day from string to integer: strconv.Atoi: parsing \"aa\": invalid syntax",
		},

		{
			upTime:      "1-5 12:ae-13:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert start minute day from string to integer: strconv.Atoi: parsing \"ae\": invalid syntax",
		},

		{
			upTime:      "1-5 12:12-0e:00",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert last day hour from string to integer: strconv.Atoi: parsing \"0e\": invalid syntax",
		},

		{
			upTime:      "1-5 12:12-00:e0",
			currentTime: time.Date(2022, 10, 20, 10, 30, 0, 0, loc),
			want:        "could not convert Hour to Date format: could not convert last day minute from string to integer: strconv.Atoi: parsing \"e0\": invalid syntax",
		},
	}

	for _, tt := range shouldScaleDownTestCases {
		_, err := shouldScaleDown(tt.upTime, log, tt.currentTime, loc)

		if err.Error() != tt.want {
			t.Errorf("got \"%s\", want \"%s\"", err, tt.want)

		}

	}
}
