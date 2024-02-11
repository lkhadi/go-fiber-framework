package repositories

import (
	"fmt"
	"p2h-api/app/models"
	"p2h-api/app/requests/reqCompany"
	"time"

	"gorm.io/gorm"
)

type CompanyRepositoryInterface interface {
	GetAll(request reqCompany.List) ([]models.Company, int64, error)
	GetAllDeleted() ([]models.Company, error)
	GetDetail(uuid string) (models.Company, error)
	Delete(uuid string) error
	Create(reqCompany.Save) (string, error)
	Update(reqCompany.Update) (string, error)
}

type CompanyRepository struct {
	DB *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepositoryInterface {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) GetAll(request reqCompany.List) ([]models.Company, int64, error) {
	var company []models.Company
	var totalRow int64
	query := r.DB.Model(&models.Company{}).Where("deleted_at IS NULL")

	if request.Search != "" {
		query.Where("(nama_perusahaan LIKE ?)", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
	}

	query.Count(&totalRow)

	if request.Limit != -1 {
		offset := (request.Page - 1) * request.Limit
		query.Limit(request.Limit).Offset(offset)
	}

	if len(request.Sort) > 0 {
		orderString := fmt.Sprintf("%s %s", request.Sort[0].Key, request.Sort[0].Order)
		query.Order(orderString)
	}

	if err := query.Find(&company).Error; err != nil {
		return nil, totalRow, err
	}

	return company, totalRow, nil
}

func (r *CompanyRepository) GetAllDeleted() ([]models.Company, error) {
	var company []models.Company
	if err := r.DB.Where("deleted_at IS NOT NULL").Find(&company).Error; err != nil {
		return company, err
	}

	return company, nil
}

func (r *CompanyRepository) GetDetail(uuid string) (models.Company, error) {
	var company models.Company
	if err := r.DB.Where("uuid = ?", uuid).First(&company).Error; err != nil {
		return company, err
	}

	return company, nil
}

func (r *CompanyRepository) Delete(uuid string) error {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")

	if err := r.DB.Model(&models.Company{}).Where("uuid = ?", uuid).Update("deleted_at", currentDateTime).Error; err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) Create(request reqCompany.Save) (string, error) {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := ""
	company := models.Company{
		CompanyName: request.Nama_Perusahaan,
		CreatedAt:   currentDateTime,
		UpdatedAt:   currentDateTime,
	}

	if request.Logo != nil {
		timestamp := time.Now().Format("20060102150405")
		fileName = timestamp + request.Logo.Filename
		company.Logo = &fileName
	}

	if err := r.DB.Create(&company).Error; err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (r *CompanyRepository) Update(request reqCompany.Update) (string, error) {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := ""
	var company models.Company
	companyData := models.Company{
		CompanyName: request.Nama_Perusahaan,
		UpdatedAt:   currentDateTime,
	}

	if request.Logo != nil {
		timestamp := time.Now().Format("20060102150405")
		fileName = timestamp + request.Logo.Filename
		companyData.Logo = &fileName
	}

	if err := r.DB.Model(&company).Where("uuid = ?", request.UUID).Updates(companyData).Error; err != nil {
		return fileName, err
	}

	return fileName, nil
}
