package repository

import (
	"database/sql"
	"project-root/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, password, company_id) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, user.ID, user.Username, user.Password, user.CompanyID)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, company_id FROM users WHERE username=$1`
	row := r.DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CompanyID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
