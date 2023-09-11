package util

import (
	"fmt"
	"strings"
	"time"
)

func ParsePostBackTime(s string) (time.Time, error) {
	parts := strings.Split(s, " ")
	formattedStr := fmt.Sprintf("%sT%s+07:00", parts[0], parts[1])
	t, err := time.Parse(time.RFC3339, formattedStr)
	return t, err
}

func GetSinceUntilTime(t time.Time, nDate int) (time.Time, time.Time) {
	since := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	nHours := time.Duration(nDate * 24)
	until := since.Add(time.Duration(nHours * time.Hour))
	return since, until
}
