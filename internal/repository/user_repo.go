package repository

import (
	"database/sql"
	"hakaton/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers возвращает всех пользователей
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.salt, u.company_id, c.name AS company_name
	FROM users u
	LEFT JOIN companies c ON u.company_id = c.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Salt, &user.CompanyID, &user.CompanyName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByEmail получает пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, salt FROM users WHERE email = $1`
	var user models.User

	// Выполняем запрос в базу данных и сканируем результат в структуру user
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Salt)
	if err != nil {
		// Если ошибка - возвращаем ее
		return nil, err
	}

	// Возвращаем найденного пользователя
	return &user, nil
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(email, passwordHash, salt, companyId string) error {
	query := `INSERT INTO users (email, password_hash, salt, company_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := r.db.Exec(query, email, passwordHash, salt, companyId)
	return err
}

// UpdateUser обновляет информацию о пользователе
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, salt = $3, company_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`

	_, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Salt, user.CompanyID, user.ID)
	return err
}

// DeleteUser удаляет пользователя по его ID
func (r *UserRepository) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, userID)
	return err
}
