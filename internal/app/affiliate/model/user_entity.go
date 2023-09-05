package model

import "time"

type UserEntity struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePicture"`
	Password       string    `json:"password"`
	RegisterType   string    `json:"registerType"`
	EmailActivated string    `json:"emailActivated"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Website        string    `json:"website"`
	Phone          string    `json:"phone"`
	WalletAddress  string    `json:"walletAddress"`
}

func (e *UserEntity) TableName() string {
	return "users"
}
