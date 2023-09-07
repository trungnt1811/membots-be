package redeem

import (
	"gorm.io/gorm"
)

type RedeemRepository struct {
	db *gorm.DB
}

func NewRedeemRepository(db *gorm.DB) *RedeemRepository {
	return &RedeemRepository{
		db: db,
	}
}

func (r *RedeemRepository) TopUpCashBack(userId int, amount float32) error {
	return nil
}
