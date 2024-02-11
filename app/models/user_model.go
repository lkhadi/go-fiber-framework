package models

type User struct {
	ID          uint    `json:"-" gorm:"primaryKey;column:id"`
	Name        string  `json:"nama" gorm:"column:nama"`
	EmployeeID  string  `json:"nip" gorm:"column:nip"`
	Role        string  `json:"role" gorm:"column:role"`
	Email       string  `json:"email" gorm:"column:email"`
	Password    string  `json:"-" gorm:"column:password"`
	Department  string  `json:"department" gorm:"column:department"`
	Section     string  `json:"section" gorm:"column:section"`
	Position    string  `json:"jabatan" gorm:"column:jabatan"`
	HandSign    string  `json:"tanda_tangan" gorm:"column:tanda_tangan"`
	CreatedAt   string  `json:"-" gorm:"column:created_at"`
	UpdatedAt   string  `json:"-" gorm:"column:updated_at"`
	AccessToken string  `json:"access_token" gorm:"-"`
	CompanyID   *uint   `json:"-" gorm:"column:id_perusahaan"`
	Company     Company `json:"perusahaan" gorm:"foreignKey:CompanyID;references:ID"`
	UUID        *string `json:"-" gorm:"column:uuid"`
}

func (User) TableName() string {
	return "user"
}
