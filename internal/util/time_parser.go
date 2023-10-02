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

func GetSinceUntilTime(middle time.Time, nDate int) (time.Time, time.Time) {
	nHours := time.Duration(nDate * 24)
	since := middle.Add(time.Duration(-nHours * time.Hour))
	until := middle.Add(time.Duration(nHours * time.Hour))
	return since, until
}
