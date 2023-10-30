package model

import (
	"time"

	"gorm.io/datatypes"
)

type AffPostBackLog struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	OrderId      string         `json:"order_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Data         datatypes.JSON `json:"data"`
	ErrorMessage string         `json:"error_message"`
}

func (m *AffPostBackLog) TableName() string {
	return "aff_postback_log"
}
