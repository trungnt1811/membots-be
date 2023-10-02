package model

import "time"

const (
	IMPORT_STATUS_PENDING   = "pending"
	IMPORT_STATUS_IMPORTED  = "imported"
	IMPORT_STATUS_SUBMITTED = "submited"
	IMPORT_STATUS_COMPLETED = "completed"
	IMPORT_STATUS_FAILED    = "failed"
	IMPORT_STATUS_SUCCESS   = "success"
)

type ImportEntity struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	Status           string    `json:"status"`
	SellerId         int       `json:"sellerId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	FileName         string    `json:"fileName"`
	OriginalFileName string    `json:"originalFileName"`
}

type OrderEntity struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	ImportId      int       `json:"import_id"`
	Amount        int       `json:"amount"`
	Tier          string    `json:"tier"`
	IsFirst       int8      `json:"isFirst"`
	SellerId      int       `json:"sellerId"`
	CustomerId    int       `json:"customerId"`
	RewardId      int       `json:"rewardId"` // Be careful for NULL or 0 values
	InitializedAt time.Time `json:"initializedAt"`
}

type OrderByTime []OrderEntity

func (a OrderByTime) Len() int           { return len(a) }
func (a OrderByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderByTime) Less(i, j int) bool { return a[i].InitializedAt.Before(a[j].InitializedAt) }
