package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestShouldScaleDown(t *testing.T) {
	var log = logrus.New()
	loc, _ := time.LoadLocation("Europe/Paris")

	errorTestCases := []struct {
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
	}

	for _, tt := range errorTestCases {
		got, _ := shouldScaleDown(tt.upTime, log, tt.currentTime, loc)
		if got != tt.want {
			t.Errorf("got %t want %t", got, tt.want)
		}
	}
}
