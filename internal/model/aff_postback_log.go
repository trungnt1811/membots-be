package model

import "time"

type AffPostBackLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	OrderId   string    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      any       `json:"data"`
}

func (m *AffPostBackLog) TableName() string {
	return "aff_postback_log"
}
