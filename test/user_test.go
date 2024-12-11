package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"hakaton/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"hakaton/internal/models"
)

func TestGetAllUsers(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем экземпляр репозитория
	repo := repository.NewUserRepository(db)

	// Подготавливаем ожидаемое поведение для mock
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "salt", "company_id", "company_id", "company_name"}).
		AddRow("1", "test@example.com", "passwordHash", "salt", "company1", "company1", "Test Company")

	// Ожидаем, что будет выполнен запрос и возвращены данные
	mock.ExpectQuery(`SELECT u.id, u.email, u.password_hash, u.salt, u.company_id, c.id AS company_id, c.name AS company_name`).
		WillReturnRows(rows)

	// Вызываем метод GetAllUsers
	users, err := repo.GetAllUsers()

	// Проверяем ошибки
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "test@example.com", users[0].Email)
	assert.Equal(t, "Test Company", users[0].Company.Name)

	// Проверяем, что все ожидания для mock были выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем экземпляр репозитория
	repo := repository.NewUserRepository(db)

	// Подготавливаем ожидаемое поведение для mock
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "salt", "company_id", "company_id", "company_name"}).
		AddRow("1", "test@example.com", "passwordHash", "salt", "company1", "company1", "Test Company")

	// Ожидаем, что будет выполнен запрос и возвращены данные
	mock.ExpectQuery(`SELECT u.id, u.email, u.password_hash, u.salt, u.company_id, c.id AS company_id, c.name AS company_name`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Вызываем метод GetUserByEmail
	user, err := repo.GetUserByEmail("test@example.com")

	// Проверяем ошибки
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test Company", user.Company.Name)

	// Проверяем, что все ожидания для mock были выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем экземпляр репозитория
	repo := repository.NewUserRepository(db)

	// Подготавливаем ожидаемое поведение для mock
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "salt", "company_id"}).
		AddRow("1", "test@example.com", "passwordHash", "salt", "company1")

	// Ожидаем, что будет выполнен запрос и возвращены данные
	mock.ExpectQuery(`INSERT INTO users (id, email, password_hash, salt, company_id, created_at, updated_at)`).
		WithArgs("1", "test@example.com", "passwordHash", "salt", "company1").
		WillReturnRows(rows)

	// Вызываем метод CreateUser
	user, err := repo.CreateUser("1", "test@example.com", "passwordHash", "salt", "company1")

	// Проверяем ошибки
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	// Проверяем, что все ожидания для mock были выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем экземпляр репозитория
	repo := repository.NewUserRepository(db)

	// Ожидаем, что будет выполнен запрос на обновление данных пользователя
	mock.ExpectExec(`UPDATE users SET email = $1, password_hash = $2, salt = $3, company_id = $4, updated_at = CURRENT_TIMESTAMP`).
		WithArgs("new@example.com", "newPasswordHash", "newSalt", "company2", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создаем пользователя
	user := &models.User{
		ID:           "1",
		Email:        "new@example.com",
		PasswordHash: "newPasswordHash",
		Salt:         "newSalt",
		CompanyID:    "company2",
	}

	// Вызываем метод UpdateUser
	err = repo.UpdateUser(user)

	// Проверяем ошибки
	assert.NoError(t, err)

	// Проверяем, что все ожидания для mock были выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем экземпляр репозитория
	repo := repository.NewUserRepository(db)

	// Ожидаем, что будет выполнен запрос на удаление пользователя
	mock.ExpectExec(`DELETE FROM users WHERE id = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем метод DeleteUser
	err = repo.DeleteUser(1)

	// Проверяем ошибки
	assert.NoError(t, err)

	// Проверяем, что все ожидания для mock были выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
