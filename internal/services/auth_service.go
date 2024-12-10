package services

import (
	"errors"
	"hakaton/internal/models"
	"hakaton/internal/repository"
	"hakaton/pkg/utils"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func (s *AuthService) RegisterUser(user *models.User) error {
	// Хэшируем пароль
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Сохраняем пользователя
	return s.UserRepo.CreateUser(user)
}

func (s *AuthService) LoginUser(username, password string) (string, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Генерируем JWT
	return utils.GenerateJWT(user.ID, user.CompanyID)
}
