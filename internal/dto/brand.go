package dto

type BrandDto struct {
	ID         uint32  `json:"id"`
	Name       string  `json:"name"`
	Logo       string  `json:"logo"`
	CoverPhoto *string `json:"cover_photo"`
}

type UserViewBrandDto struct {
	UserId  uint32 `json:"user_id"`
	BrandId uint64 `json:"brand_id"`
}
