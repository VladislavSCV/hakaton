package repository

import (
	"database/sql"
	uuid2 "github.com/google/uuid"
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
	query := `SELECT id, name, created_at, updated_at FROM companies`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var company models.Company
		err := rows.Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// GetCompanyByID возвращает компанию по ID
func (r *CompanyRepository) GetCompanyByID(companyID int) (*models.Company, error) {
	query := `SELECT id, name, created_at, updated_at FROM companies WHERE id = $1`

	row := r.db.QueryRow(query, companyID)

	var company models.Company
	err := row.Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

// CreateCompany создает новую компанию
func (r *CompanyRepository) CreateCompany(company *models.Company) (models.Company, error) {
	var comp models.Company
	query := `INSERT INTO companies (id, name, created_at, updated_at) 
	VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, name`
	uuid := uuid2.New().String()
	company.ID = uuid

	//err := r.db.QueryRow(query, company.ID, company.Name).Scan(&comp.ID, comp.Name)
	err := r.db.QueryRow(query, company.ID, company.Name).Scan(&comp.ID, &comp.Name)

	return comp, err
}

// DeleteCompany удаляет компанию по ID
func (r *CompanyRepository) DeleteCompany(companyID int) error {
	query := `DELETE FROM companies WHERE id = $1`

	_, err := r.db.Exec(query, companyID)
	return err
}
