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
