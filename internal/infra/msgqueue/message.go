package msgqueue

const (
	NotiCategoryCommerce     = "commerce"
	NotiCategorySystem       = "system"
	NotiDataTypeCouponDetail = "coupon-detail"
	NotiDataKeyType          = "type"
	NotiDataKeyId            = "id"
	NotiDataKeyOrderId       = "orderID"
)

type MsgOrderApproved struct {
	AtOrderID string `json:"accesstrade_order_id"`
}

type MsgOrderUpdated struct {
	AtOrderID   string `json:"accesstrade_order_id"`
	OrderStatus string `json:"orderStatus"`
	IsConfirmed uint8  `json:"is_confirmed"`
}

type AppNotiMsg struct {
	// Category is the category of the message.
	Category string `json:"Category,omitempty"`

	// Topic is the topic of the message.
	Topic string `json:"topic,omitempty"`

	// Data is the message custom data.
	Data map[string]string `json:"data,omitempty"`

	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
	UserId uint   `json:"userId" binding:"required"`
}

func GetOrderApprovedNotiData() map[string]string {
	data := make(map[string]string)
	// data[NotiDataKeyType] = NotiDataTypeCouponDetail
	// data[NotiDataKeyId] = strconv.Itoa(int(couponId))
	// data[NotiDataKeyOrderId] = strconv.Itoa(int(orderId))
	return data
}

// Reward-shipping delivery receipt
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
	OrderCode       string `json:"orderCode"`
	RedeemExpiredAt int64  `json:"redeemExpiredAt"`
	TokenAddress    string `json:"tokenAddress"`
	Amount          string `json:"amount"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phoneNumber"`
}
