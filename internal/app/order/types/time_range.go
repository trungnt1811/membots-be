package types

import "time"

type TimeRange struct {
	Since *time.Time `json:"since" form:"since"`
	Until *time.Time `json:"until" form:"until"`
}
