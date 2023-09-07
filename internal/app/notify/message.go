package notify

const (
	NotiCategoryWallet   = "wallet"
	NotiDataTypeTxDetail = "tx-detail"
	NotiDataKeyType      = "type"
	NotiDataKeyId        = "id"
)

type NotiMsg struct {
	// Category is the category of the message.
	Category string `json:"category,omitempty"`

	// Topic is the topic of the message.
	Topic string `json:"topic,omitempty"`

	// Data is the message custom data.
	Data map[string]string `json:"data,omitempty"`

	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
	UserId uint   `json:"userId" binding:"required"`
}

func GetRewardNotiCustomData(txHash string) map[string]string {
	data := make(map[string]string)
	data[NotiDataKeyType] = NotiDataTypeTxDetail
	data[NotiDataKeyId] = txHash
	return data
}
