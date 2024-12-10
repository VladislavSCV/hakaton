package models

type Game struct {
	ID          string `json:"id"`
	CompanyID   string `json:"company_id"`
	Name        string `json:"name"`
	Data        string `json:"data"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CompanyName string `json:"company_name"`
}
