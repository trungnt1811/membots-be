package shipping

import "time"

type ExchangeErrorResponse struct {
	Errors []string `json:"errors"`
}

type ReqSendItem struct {
	OrderCode     string `form:"orderCode" json:"orderCode,omitempty"`         // Customer order code
	Amount        string `form:"amount" json:"amount" binding:"required"`      // Amount of reward will be send
	CustomerId    int    `form:"customerId" json:"customerId,omitempty"`       // Optional, prefer customer id over wallet address and other information
	Name          string `form:"name" json:"name,omitempty"`                   // Optional, if provided
	Email         string `form:"email" json:"email,omitempty"`                 // Optional, either send via email/sms or wallet
	WalletAddress string `form:"walletAddress" json:"walletAddress,omitempty"` // Optional, either send via email/sms or wallet
	PhoneNumber   string `form:"phoneNumber" json:"phoneNumber,omitempty"`     // Optional, either send via email/sms or wallet
}

// ReqSendPayload is a request to send predefined amount of reward (ASA)
type ReqSendPayload struct {
	SellerId       uint          `form:"sellerId" json:"sellerId" binding:"required"`             // Required seller id
	RequestId      string        `form:"requestId" json:"requestId" binding:"required"`           // Required unique request id
	ProgramAddress string        `form:"programAddress" json:"programAddress" binding:"required"` // Program to interact, required
	Items          []ReqSendItem `form:"items" json:"items"`                                      // List of reward info to be sent, must not empty
}

type RewardEntity struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ProgramId     int       `json:"programId"`
	Amount        float64   `json:"amount"`
	AmountAastra  string    `json:"amountAastra"`
	TierAmount    float64   `json:"tierAmount"`
	Status        string    `json:"status"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	WalletAddress string    `json:"walletAddress"`
	TxHash        string    `json:"txHash"`
	SellerId      int       `json:"sellerId"`
	CustomerId    int       `json:"customerId"`
	TxError       string    `json:"txError"`
	ImportId      int       `json:"importId"`
	OrderId       *int      `json:"orderId"`
	RequestId     string    `json:"requestId"`
	DeliveryId    int       `json:"deliveryId"`
	QtyTraceback  string    `json:"qtyTraceback"`
	Type          string    `json:"type"`
}

type ReqSendResponse struct {
	Rewards []RewardEntity `json:"rewards"`
}
