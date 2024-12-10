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

// SaveImageForGame сохраняет изображение для игры
func (r *ImageRepository) SaveImageForGame(companyID, gameID string, url string) error {
	query := `INSERT INTO images (company_id, game_id, url, created_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

	_, err := r.db.Exec(query, companyID, gameID, url)
	return err
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
