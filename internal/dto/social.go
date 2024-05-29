package dto

type Social struct {
	Provider                     string         `json:"provider"`
	Username                     string         `json:"username"`
	DisplayName                  string    		`json:"displayName"`
	PhotoURL                     string    		`json:"photoUrl"`
	URL                          string       	`json:"url"`
}