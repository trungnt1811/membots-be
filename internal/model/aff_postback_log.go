package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
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

func (m *AffPostBackLog) ToDto() dto.AffPostBack {
	return dto.AffPostBack{
		ID:        m.ID,
		OrderId:   m.OrderId,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Data:      string(m.Data),
	}
}
