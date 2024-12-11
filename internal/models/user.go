package models

import "encoding/json"

type User struct {
	ID           string          `json:"id"`
	CompanyID    string          `json:"company_id"`
	Email        string          `json:"email"`
	PasswordHash string          `json:"password_hash"`
	Salt         string          `json:"salt"`
	Data         json.RawMessage `json:"data"`
	CreatedAt    string          `json:"created_at"`
	Company      *Company        `json:"company"`
}
