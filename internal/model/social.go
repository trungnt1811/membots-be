package model

import "github.com/flexstack.ai/membots-be/internal/dto"

type Social struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	MemeID      uint   `json:"meme_id"`
	Provider    string `json:"provider"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	PhotoURL    string `json:"photo_url"`
	URL         string `json:"url"`
}

func (m *Social) TableName() string {
	return "social"
}

func (m *Social) ToMapDto() map[string]dto.Social {
	socialDto := dto.Social{
		Provider:    m.Provider,
		Username:    m.Username,
		DisplayName: m.DisplayName,
		PhotoURL:    m.PhotoURL,
		URL:         m.URL,
	}
	return map[string]dto.Social{m.Provider: socialDto}
}
