package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetClaimableReward(t *testing.T) {
	layout := "2006-01-02T15:04:05.000Z"

	createdTime := "2014-11-12T11:45:26.371Z"
	createdAt, err := time.Parse(layout, createdTime)
	assert.Nil(t, err)

	endedAt := createdAt.Add(60 * time.Hour * 24)
	fmt.Println("endedAt", endedAt)

	currentTime := "2014-12-03T07:45:26.371Z"
	now, err := time.Parse(layout, currentTime)
	assert.Nil(t, err)

	daysPassed := int(now.Sub(createdAt) / OneDay)
	totalDays := int(endedAt.Sub(createdAt) / OneDay) // total lock days
	withdrawablePercent := float64(daysPassed) / float64(totalDays)

	fmt.Println("daysPassed", daysPassed)
	fmt.Println("totalDays", totalDays)
	fmt.Println("withdrawablePercent", withdrawablePercent)
}
