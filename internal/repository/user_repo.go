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
	query := `SELECT u.id, u.email, u.password_hash, u.salt, u.company_id, c.id AS company_id, c.name AS company_name
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
		var company models.Company

		// Сканируем данные из строки
		err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Salt, &user.CompanyID, &company.ID, &company.Name)
		if err != nil {
			return nil, err
		}

		// Если компания существует, добавляем её в пользователя
		if company.ID != "" {
			user.Company = &company
		}

		users = append(users, user)
	}

	// Проверка на ошибки после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.salt, u.company_id, c.id AS company_id, c.name AS company_name
	          FROM users u 
	          LEFT JOIN companies c ON u.company_id = c.id 
	          WHERE u.email = $1`

	var user models.User
	var company models.Company

	// Выполняем запрос и сканируем результаты
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Salt, &user.CompanyID, &company.ID, &company.Name)
	if err != nil {
		return nil, err
	}

	// Если компания существует, добавляем её в пользователя
	if company.ID != "" {
		user.Company = &company
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(uuid, email, passwordHash, salt, companyId string) (*models.User, error) {
	var user models.User
	var company models.Company

	// Запрос для вставки нового пользователя и возврата всех данных
	query := `INSERT INTO users (id, email, password_hash, salt, company_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
			  RETURNING id, email, password_hash, salt, company_id`

	// Выполняем запрос с возвратом данных
	err := r.db.QueryRow(query, uuid, email, passwordHash, salt, companyId).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Salt, &user.CompanyID)
	if err != nil {
		return nil, err
	}

	// Дополнительно извлекаем информацию о компании (если нужно)
	companyQuery := `SELECT id, name FROM companies WHERE id = $1`
	err = r.db.QueryRow(companyQuery, user.CompanyID).Scan(&company.ID, &company.Name)
	if err != nil {
		return nil, err
	}

	// Связываем данные компании с пользователем
	user.Company = &company

	// Возвращаем полные данные пользователя
	return &user, nil
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
