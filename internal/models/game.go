package models

type Game struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Name      string `json:"name"`
	Data      struct {
		Store []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Count    int    `json:"count"`
			Duration int    `json:"duration"`
			Level    int    `json:"level"`
			Style    struct {
				Background   string `json:"background"`
				BorderRadius int    `json:"borderRadius"`
				Border       string `json:"border"`
			} `json:"style"`
		} `json:"store"`
		Field struct {
			Count           int    `json:"count"`
			TimerActive     string `json:"timerActive"`
			TimerDuration   int    `json:"timerDuration"`
			StartGameActive bool   `json:"startGameActive"`
		} `json:"field"`
	} `json:"data"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CompanyName string `json:"company_name"`
}
