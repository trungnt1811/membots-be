package dto

type BrandDto struct {
	ID         uint32  `json:"id"`
	Name       string  `json:"name"`
	Logo       string  `json:"logo"`
	CoverPhoto *string `json:"cover_photo"`
}
