package repository

import (
	"database/sql"
	"hakaton/internal/models"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// GetAllCompanies возвращает список всех компаний
func (r *CompanyRepository) GetAllCompanies() ([]models.Company, error) {
	query := `SELECT id, name, api_key, created_at, updated_at FROM companies`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var company models.Company
		err := rows.Scan(&company.ID, &company.Name, &company.APIKey, &company.CreatedAt, &company.UpdatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// GetCompanyByID возвращает компанию по ID
func (r *CompanyRepository) GetCompanyByID(companyID int) (*models.Company, error) {
	query := `SELECT id, name, api_key, created_at, updated_at FROM companies WHERE id = $1`

	row := r.db.QueryRow(query, companyID)

	var company models.Company
	err := row.Scan(&company.ID, &company.Name, &company.APIKey, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

// CreateCompany создает новую компанию
func (r *CompanyRepository) CreateCompany(company *models.Company) error {
	query := `INSERT INTO companies (name, api_key, created_at, updated_at) 
	VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := r.db.Exec(query, company.Name, company.APIKey)
	return err
}

// DeleteCompany удаляет компанию по ID
func (r *CompanyRepository) DeleteCompany(companyID int) error {
	query := `DELETE FROM companies WHERE id = $1`

	_, err := r.db.Exec(query, companyID)
	return err
}
