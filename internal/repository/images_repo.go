package repository

import (
	"database/sql"
	"hakaton/internal/models"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// SaveImageForGame сохраняет или обновляет изображение для игры
func (r *ImageRepository) SaveImageForGame(companyID, gameID, url string) error {
	// Проверяем, существует ли запись
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM images WHERE company_id = $1 AND game_id = $2)`
	err := r.db.QueryRow(checkQuery, companyID, gameID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Обновляем запись
		updateQuery := `UPDATE images SET url = $1, created_at = CURRENT_TIMESTAMP WHERE company_id = $2 AND game_id = $3`
		_, err := r.db.Exec(updateQuery, url, companyID, gameID)
		return err
	} else {
		// Вставляем новую запись
		insertQuery := `INSERT INTO images (company_id, game_id, url, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`
		_, err := r.db.Exec(insertQuery, companyID, gameID, url)
		return err
	}
}

// GetImagesByCompanyID возвращает изображения по ID компании
func (r *ImageRepository) GetImagesByCompanyID(companyID int) ([]models.Image, error) {
	query := `SELECT id, company_id, game_id, url, created_at FROM images WHERE company_id = $1`

	rows, err := r.db.Query(query, companyID)
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
