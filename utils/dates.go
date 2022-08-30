package utils

import (
	"log"
	"strconv"
	"time"
)

type StartStop struct {
	Start string
	Stop  string
}

type TimeSeries []StartStop

// CreateTimeSeries fills all days between start and stop with valid dates
func CreateTimeSeries(start, end string) TimeSeries {
	var ts TimeSeries
	t1, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatalln(err)
	}
	t2, err := time.Parse("2006-01-02", end)
	if err != nil {
		log.Fatalln(err)
	}
	for d := t1; d.After(t2) == false; d = d.AddDate(0, 0, 1) {
		var t StartStop
		t.Start = d.Format("2006-01-02")
		t.Stop = d.Format("2006-01-02")
		ts = append(ts, t)
	}
	return ts
}

func convStrToI(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln(err)
	}
	return i
}
