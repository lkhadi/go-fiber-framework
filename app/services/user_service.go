package services

import (
	"fmt"
	"os"
	"p2h-api/app/models"
	"p2h-api/app/repositories"
	"p2h-api/app/requests/reqUser"
	"p2h-api/app/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	User repositories.UserRepositoryInterface
}

func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		User: repository,
	}
}

func (us *UserService) UserAuthentication(email string, password string) (models.Response, error) {
	user, err := us.User.GetUserByEmail(email)

	if err != nil {
		return models.Response{
			Code:    401,
			Status:  false,
			Message: "Akses ditolak!",
		}, err
	}

	if user.Role != "adminsystem" && user.Company.ID == nil {
		return models.Response{
			Code:    401,
			Status:  false,
			Message: "Akses ditolak!",
		}, err
	}

	errPasswordChecking := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if errPasswordChecking != nil {
		return models.Response{
			Code:    401,
			Status:  false,
			Message: "Akses ditolak!",
		}, err
	}

	myHandSign, _ := utils.GeneratePresignedURL("berkas/ttd/" + user.HandSign)
	user.HandSign = myHandSign
	ttl, _ := strconv.Atoi(os.Getenv("JWT_TTL"))
	claims := models.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
		},
		ID:         user.ID,
		Name:       user.Name,
		EmployeeID: user.EmployeeID,
		Role:       user.Role,
		Email:      user.Email,
		Department: user.Department,
		Section:    user.Section,
		Position:   user.Position,
		HandSign:   user.HandSign,
		CompanyID:  user.Company.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		fmt.Println("Error signing the token:", err)
		return models.Response{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
		}, err
	}

	user.AccessToken = tokenString

	if user.Company.Logo != nil {
		logo := *user.Company.Logo
		logo, _ = utils.GeneratePresignedURLLogo("company/logo/" + logo)
		user.Company.Logo = &logo
	}

	return models.Response{
		Code:    200,
		Status:  true,
		Message: "Data user",
		Data:    user,
	}, nil
}

func (us *UserService) Create(role string, request reqUser.Save) models.Response {
	if role != "adminsystem" && request.Role == "adminsystem" {
		return models.Response{
			Code:    403,
			Status:  false,
			Message: "Akses ditolak",
		}
	}

	if _, err := us.User.GetUserByEmail(request.Email); err == nil {
		return models.Response{
			Code:    409,
			Status:  false,
			Message: "Data pengguna sudah ada",
		}
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	request.Password = string(hashed)
	fileName, err := us.User.Create(request)

	if err != nil {
		return models.Response{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
		}
	}

	if request.Tanda_Tangan != nil {
		fileData, err := request.Tanda_Tangan.Open()

		if err != nil {
			return models.Response{
				Code:    500,
				Status:  false,
				Message: err.Error(),
			}
		}

		defer fileData.Close()
		fileName := "berkas/ttd/" + fileName
		utils.UploadToS3(fileName, fileData)
	}

	return models.Response{
		Code:    201,
		Status:  true,
		Message: "Berhasil membuat data pengguna",
	}
}

func (us *UserService) Update(role string, request reqUser.Update) models.Response {
	if role != "adminsystem" && request.Role == "adminsystem" {
		return models.Response{
			Code:    403,
			Status:  false,
			Message: "Akses ditolak",
		}
	}

	if request.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		request.Password = string(hashed)
	}

	fileName, err := us.User.Update(request)

	if err != nil {
		return models.Response{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
		}
	}

	if request.Tanda_Tangan != nil {
		fileData, err := request.Tanda_Tangan.Open()

		if err != nil {
			return models.Response{
				Code:    500,
				Status:  false,
				Message: err.Error(),
			}
		}

		defer fileData.Close()
		fileName := "berkas/ttd/" + fileName
		utils.UploadToS3(fileName, fileData)
	}

	return models.Response{
		Code:    200,
		Status:  true,
		Message: "Berhasil memperbarui data pengguna",
	}
}

func (us *UserService) GetAllUser(user *models.CustomClaims, request reqUser.List) models.PaginationResponse {
	data, totalRow, err := us.User.GetAllUser(user, request)

	if err != nil {
		return models.PaginationResponse{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
			Total:   0,
		}
	}

	return models.PaginationResponse{
		Code:    200,
		Status:  true,
		Message: "Data Unit",
		Data:    data,
		Total:   totalRow,
		Page:    request.Page,
		Limit:   request.Limit,
	}
}

func (us *UserService) Delete(user *models.CustomClaims, uuid string) models.Response {
	err := us.User.Delete(user.CompanyID, user.Role, uuid)

	if err != nil {
		return models.Response{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
		}
	}

	return models.Response{
		Code:    200,
		Status:  true,
		Message: "Berhasil menghapus data pengguna",
	}
}

func (us *UserService) GetUserByID(userId uint) models.Response {
	user, _ := us.User.GetUserByID(userId)
	myHandSign, _ := utils.GeneratePresignedURL("berkas/ttd/" + user.HandSign)
	user.HandSign = myHandSign
	return models.Response{
		Code:    200,
		Status:  true,
		Message: "My Profile",
		Data:    user,
	}
}

func (us *UserService) UpdateProfile(userId uint, request reqUser.UpdateProfile) models.Response {
	if request.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		request.Password = string(hashed)
	}

	fileName, err := us.User.UpdateProfile(userId, request)

	if err != nil {
		return models.Response{
			Code:    500,
			Status:  false,
			Message: "Terjadi kesalahan server",
		}
	}

	if request.Tanda_Tangan != nil {
		fileData, err := request.Tanda_Tangan.Open()

		if err != nil {
			return models.Response{
				Code:    500,
				Status:  false,
				Message: err.Error(),
			}
		}

		defer fileData.Close()
		fileName := "berkas/ttd/" + fileName
		utils.UploadToS3(fileName, fileData)
	}

	user, _ := us.User.GetUserByID(userId)
	return models.Response{
		Code:    200,
		Status:  true,
		Message: "Berhasil memperbarui data profile",
		Data:    user,
	}
}
