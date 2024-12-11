package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hakaton/internal/models"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

// GetAllGames возвращает все игры
func (r *GameRepository) GetAllGames() ([]models.Game, error) {
	query := `SELECT g.id, g.name, g.data, g.company_id, g.created_at, g.updated_at, c.name AS company_name 
	FROM games g
	LEFT JOIN companies c ON g.company_id = c.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(&game.ID, &game.Name, &game.Data, &game.CompanyID, &game.CreatedAt, &game.UpdatedAt, &game.CompanyName)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

// GetGameByID возвращает игру по ID
func (r *GameRepository) GetGameByID(gameID int) (*models.Game, error) {
	query := `SELECT id, name, data, company_id, created_at, updated_at FROM games WHERE id = $1`

	row := r.db.QueryRow(query, gameID)

	var game models.Game
	err := row.Scan(&game.ID, &game.Name, &game.Data, &game.CompanyID, &game.CreatedAt, &game.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// GetGameByName получает игру по названию и ID компании
func (r *GameRepository) GetGameByName(companyID, name string) (*models.Game, error) {
	query := `SELECT id, company_id, name, data FROM games WHERE company_id = $1 AND name = $2`
	var game models.Game

	// Выполняем запрос в базу данных и сканируем результат в структуру game
	err := r.db.QueryRow(query, companyID, name).Scan(&game.ID, &game.CompanyID, &game.Name, &game.Data)
	if err != nil {
		// Если ошибка - возвращаем ее
		return nil, err
	}

	// Возвращаем найденную игру
	return &game, nil
}

// CreateOrUpdateGame обрабатывает создание или обновление игры
func (r *GameRepository) CreateOrUpdateGame(compID, name string, data models.Game) error {
	// Проверяем, существует ли игра
	gameID, err := r.GetGameID(compID, name)
	if err != nil {
		return fmt.Errorf("failed to check game existence: %w", err)
	}

	// Сериализуем структуру данных в JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal game data: %w", err)
	}

	if gameID > 0 {
		// Игра существует — обновляем данные
		updateQuery := `UPDATE games SET data = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		_, err = r.db.Exec(updateQuery, dataJSON, gameID)
		if err != nil {
			return fmt.Errorf("failed to update game: %w", err)
		}
	} else {
		// Игра не существует — создаём новую запись
		insertQuery := `INSERT INTO games (company_id, name, data, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
		_, err = r.db.Exec(insertQuery, compID, name, dataJSON)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}
	}

	return nil
}

func (r *GameRepository) GetGameID(compID, name string) (int, error) {
	var gameID int
	query := `SELECT id FROM games WHERE company_id = $1 AND name = $2`
	err := r.db.QueryRow(query, compID, name).Scan(&gameID)

	if err == sql.ErrNoRows {
		return 0, nil // Игры нет в базе
	}
	if err != nil {
		return 0, err // Произошла ошибка
	}

	return gameID, nil // Возвращаем ID игры
}

//
//// CreateOrUpdateGame создает новую игру или обновляет существующую
//func (r *GameRepository) CreateOrUpdateGame(name, data, compID string) error {
//	// Проверяем, существует ли запись
//	var exists bool
//	checkQuery := `SELECT EXISTS (SELECT 1 FROM games WHERE name = $1 AND company_id = $2)`
//	err := r.db.QueryRow(checkQuery, name, compID).Scan(&exists)
//	if err != nil {
//		return err
//	}
//
//	if exists {
//		// Обновляем существующую запись
//		updateQuery := `UPDATE games SET data = $1, updated_at = CURRENT_TIMESTAMP WHERE name = $2 AND company_id = $3`
//		_, err = r.db.Exec(updateQuery, data, name, compID)
//		return err
//	} else {
//		// Вставляем новую запись
//		insertQuery := `INSERT INTO games (name, data, company_id, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
//		_, err = r.db.Exec(insertQuery, name, data, compID)
//		return err
//	}
//}

// UpdateGame обновляет данные игры по ее ID
func (r *GameRepository) UpdateGame(gameID, data string) error {
	// SQL запрос для обновления данных игры по ее ID
	query := `
		UPDATE games
		SET data = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	// Выполнение запроса с передачей параметров
	_, err := r.db.Exec(query, data, gameID)
	return err
}

// DeleteGame удаляет игру по ID
func (r *GameRepository) DeleteGame(gameID int) error {
	query := `DELETE FROM games WHERE id = $1`

	_, err := r.db.Exec(query, gameID)
	return err
}

// SaveImageForGame сохраняет изображение для игры
func (h *GameRepository) SaveImageForGame(companyID, gameID string, url string) error {
	query := `INSERT INTO images (company_id, game_id, url, created_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

	_, err := h.db.Exec(query, companyID, gameID, url)
	return err
}

// GetImagesByCompanyID возвращает изображения по ID компании
func (h *GameRepository) GetImagesByCompanyID(companyID int) ([]models.Image, error) {
	query := `SELECT id, company_id, game_id, url, created_at FROM images WHERE company_id = $1`

	rows, err := h.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		err := rows.Scan(&image.ID, &image.CompanyID, &image.GameID, &image.URL, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}
