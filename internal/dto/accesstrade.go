package dto

const (
	REQ_STATUS_NEW      = "0"
	REQ_STATUS_APPROVED = "1"
	REQ_STATUS_REJECTED = "2"
)

var AtOrderStatusMap map[string]string = map[string]string{
	REQ_STATUS_NEW:      "new",
	REQ_STATUS_APPROVED: "approved",
	REQ_STATUS_REJECTED: "rejected",
}

type ATPostBackRequest struct {
	TransactionId      string `json:"transaction_id"`      // Mã unique trên hệ thống AccessTrade
	OrderId            string `json:"order_id"`            // Mã đơn hàng hiển thị trên trang pub
	CampaignId         string `json:"campaign_id"`         // ID của campaign trên hệ thống
	ProductId          string `json:"product_id"`          // Mã sản phẩm
	Quantity           int    `json:"quantity"`            // Số lượng sản phẩm
	ProductCategory    string `json:"product_category"`    // Group commission của sản phẩm
	ProductPrice       string `json:"product_price"`       // Giá của một sản phẩm
	Reward             string `json:"reward"`              // Hoa hồng nhận được
	SalesTime          string `json:"sales_time"`          // Thời gian phát sinh của đơn hàng
	Browser            string `json:"browser"`             // Trình duyệt sử dụng
	ConversionPlatform string `json:"conversion_platform"` // Platform sử dụng
	Status             string `json:"status"`              // Status của đơn hàng gồm 3 giá trị: 0: new, 1: approved, 2: rejected
	IP                 string `json:"ip"`                  // IP phát sinh đơn hàng
	Referrer           string `json:"referrer"`            // click_referrer
	ClickTime          string `json:"click_time"`          // Thời gian phát sinh click
	IsConfirmed        string `json:"is_confirmed"`        // Đơn hàng khóa data và được thanh toán: 0: chưa đối soát, 1: đã đối soát
	UTMSource          string `json:"utm_source"`          // Thông tin tùy biến pub truyền vào url trong param utm_source
	UTMCampaign        string `json:"utm_campaign"`        // Thông tin tùy biến pub truyền vào url trong param utm_campaign
	UTMContent         string `json:"utm_content"`         // Thông tin tùy biến pub truyền vào url trong param utm_content
	UTMMedium          string `json:"utm_medium"`          // Thông tin tùy biến pub truyền vào url trong param utm_medium
	CustomerType       string `json:"customer_type"`       // Thuộc tính của khách hàng phụ thuộc theo campaigns
	PublisherLoginName string `json:"publisher_login_name"`
}

type ATPostBackResponse struct {
	Success bool   `json:"success"`
	OrderId string `json:"order_id"`
}
