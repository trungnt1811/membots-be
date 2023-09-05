package util

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/segmentio/kafka-go"
)

const (
	SUCCESS                = 1
	FAILED                 = 0
	PENDING                = 2
	TXTYPE_REWARD_SHIPPING = 3
	TXTYPE_SELLER_WITHDRAW = 4
	TIMEOUT                = 5
	SHIPPINGTYPE_EMAIL     = "EMAIL"
	SHIPPINGTYPE_SMS       = "SMS"
	WEBHOOK                = "WEBHOOK"
	KAFKA                  = "KAFKA"
	WORKER_ADDR            = "WORKER"
	COUPON_OPERATORS_ADDR  = "COUPON OPERATORS"
	DEPLOYER_ADDR          = "DEPLOYER"
	REWARD_OPERATORS_ADDR  = "REWARD OPERATORS"
)

type Channel struct {
	Seller    chan *TxReceiptSt
	TxReceipt chan *TxReceiptSt
}

type TxInfo struct {
	TxHash  string
	Type    int
	Holders []common.Address
}

type TxReceiptSt struct {
	Code    int
	Receipt *types.Receipt
	Message *kafka.Message
	Hodlers []common.Address
}

func NewChannel() *Channel {
	sellerChn := make(chan *TxReceiptSt)
	receiptChn := make(chan *TxReceiptSt)

	return &Channel{
		sellerChn,
		receiptChn,
	}
}

//type RequestStatusMsg struct {
//	Timestamp time.Time `json:"Timestamp"`
//	CustomID  string    `json:"CustomID"`
//	Status    string    `json:"Status"` //"NOT FOUND", "FAILED", "REJECTED", "DELIVERED", "QUEUE"
//}

type EmailSmsResponseMsg struct {
	Timestamp time.Time `json:"Timestamp"`
	MsgID     string    `json:"MsgID"`
	CustomID  string    `json:"CustomID"`
	Status    string    `json:"Status"` //"NOT FOUND", "FAILED", "REJECTED", "DELIVERED", "QUEUE"
}

type DeliveryMsg struct {
	SellerId uint `json:"sellerId"`
	// Receipt transaction hash
	TxHash string `json:"txHash"`
	// Calling program contract
	ProgramAddress string `json:"programAddress"`
	// Type of shipping batch, WALLET or EMAIL-SMS
	ShippingType string `json:"shippingType"`
	// Transaction Status: 0 - Failed, 1 - Success
	TxStatus  uint64             `json:"txStatus"`
	RequestId string             `json:"requestId"`
	Customers []DeliveryCustomer `json:"customers"`
}

type DeliveryCustomer struct {
	CustomerId      int    `json:"customerId"`
	HolderAddress   string `json:"holderAddress"`
	CustomerAddress string `json:"customerAddress"`
	RedeemCode      string `json:"redeemCode"`
	RedeemExpiredAt int64  `json:"redeemExpiredAt"`
	TokenAddress    string `json:"tokenAddress"`
	Amount          string `json:"amount"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phoneNumber"`
}
