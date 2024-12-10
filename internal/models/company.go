package models

type Company struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
