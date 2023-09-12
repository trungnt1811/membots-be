package dto

type CategoryDto struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	TotalCoupon uint32 `json:"total_coupon"`
}
