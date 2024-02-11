package repositories

import (
	"fmt"
	"p2h-api/app/models"
	"p2h-api/app/requests/reqUser"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserByRole(role string) ([]models.User, error)
	Create(request reqUser.Save) (string, error)
	Update(request reqUser.Update) (string, error)
	GetAllUser(user *models.CustomClaims, request reqUser.List) ([]models.User, int64, error)
	Delete(companyID *uint, role string, uuid string) error
	GetUserByID(userId uint) (models.User, error)
	UpdateProfile(userId uint, request reqUser.UpdateProfile) (string, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) GetUserByID(userId uint) (models.User, error) {
	var user models.User
	if err := ur.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := ur.DB.Model(&models.User{})
	err := query.Preload("Company", func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByRole(role string) ([]models.User, error) {
	var users []models.User

	if err := ur.DB.Where("role = ? AND deleted_at IS NULL", role).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (ur *UserRepository) Create(request reqUser.Save) (string, error) {
	var company models.Company
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := ""

	if err := ur.DB.Where("uuid = ?", request.UUID_Perusahaan).First(&company).Error; err != nil {
		return fileName, err
	}

	user := models.User{
		Name:       request.Nama,
		EmployeeID: request.NIP,
		Role:       request.Role,
		Email:      request.Email,
		Password:   request.Password,
		Department: request.Department,
		Section:    request.Section,
		Position:   request.Jabatan,
		CreatedAt:  currentDateTime,
		UpdatedAt:  currentDateTime,
		CompanyID:  company.ID,
	}

	if request.Tanda_Tangan != nil {
		timestamp := time.Now().Format("20060102150405")
		fileName = timestamp + request.Tanda_Tangan.Filename
		user.HandSign = fileName
	}

	if err := ur.DB.Create(&user).Error; err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (ur *UserRepository) UpdateProfile(userId uint, request reqUser.UpdateProfile) (string, error) {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := ""

	user := models.User{
		Name:       request.Nama,
		EmployeeID: request.NIP,
		Email:      request.Email,
		Department: request.Department,
		Section:    request.Section,
		Position:   request.Jabatan,
		UpdatedAt:  currentDateTime,
	}

	if request.Password != "" {
		user.Password = request.Password
	}

	if request.Tanda_Tangan != nil {
		timestamp := time.Now().Format("20060102150405")
		fileName = timestamp + request.Tanda_Tangan.Filename
		user.HandSign = fileName
	}

	if err := ur.DB.Model(&models.User{}).Where("id = ?", userId).Updates(&user).Error; err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (ur *UserRepository) GetAllUser(user *models.CustomClaims, request reqUser.List) ([]models.User, int64, error) {
	var users []models.User
	var totalRow int64
	query := ur.DB.Model(&models.User{})
	query.Preload("Company").Where("user.id != ?", user.ID).Where("user.deleted_at IS NULL")

	if request.UUID_Perusahaan != "" && user.Role == "adminsystem" {
		query.Joins("JOIN perusahaan p ON p.id=user.id_perusahaan").Where("p.uuid = ?", request.UUID_Perusahaan)
	} else if user.Role != "adminsystem" {
		query.Joins("JOIN perusahaan p ON p.id=user.id_perusahaan").Where("id_perusahaan = ?", user.CompanyID)
	}

	if request.Search != "" {
		query.Where("(nama LIKE ? OR nip LIKE ? OR email LIKE ? OR department LIKE ? OR section LIKE ? OR jabatan LIKE ?)", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
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

	if err := query.Find(&users).Error; err != nil {
		return nil, totalRow, err
	}

	return users, totalRow, nil
}

func (ur *UserRepository) Delete(companyID *uint, role string, uuid string) error {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	query := ur.DB.Model(&models.User{})
	query.Where("uuid = ?", uuid).Where("deleted_at IS NULL")

	if role != "adminsystem" {
		query.Where("id_perusahaan = ?", companyID)
	}

	err := query.Update("deleted_at", currentDateTime).Error
	return err
}

func (ur *UserRepository) Update(request reqUser.Update) (string, error) {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := ""

	user := models.User{
		Name:       request.Nama,
		EmployeeID: request.NIP,
		Role:       request.Role,
		Email:      request.Email,
		Department: request.Department,
		Section:    request.Section,
		Position:   request.Jabatan,
		UpdatedAt:  currentDateTime,
	}

	if request.Password != "" {
		user.Password = request.Password
	}

	if request.Tanda_Tangan != nil {
		timestamp := time.Now().Format("20060102150405")
		fileName = timestamp + request.Tanda_Tangan.Filename
		user.HandSign = fileName
	}

	if err := ur.DB.Model(&models.User{}).Where("uuid = ?", request.UUID).Updates(&user).Error; err != nil {
		return fileName, err
	}

	return fileName, nil
}
