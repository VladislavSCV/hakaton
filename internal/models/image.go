package models

type Image struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	GameID    string `json:"game_id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}
