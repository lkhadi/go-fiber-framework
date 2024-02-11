package models

type Company struct {
	ID          *uint   `json:"-" gorm:"primaryKey;column:id"`
	CompanyName string  `json:"nama_perusahaan" gorm:"column:nama_perusahaan"`
	Logo        *string `json:"logo" gorm:"column:file_logo"`
	DeletedAt   *string `json:"-" gorm:"column:deleted_at"`
	CreatedAt   string  `json:"-" gorm:"column:created_at"`
	UpdatedAt   string  `json:"-" gorm:"column:updated_at"`
	UUID        *string `json:"uuid" gorm:"column:uuid"`
}

func (Company) TableName() string {
	return "perusahaan"
}
