package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	jwt.StandardClaims
	ID         uint   `json:"ID"`
	Name       string `json:"Name"`
	EmployeeID string `json:"EmployeeID"`
	Role       string `json:"Role"`
	Email      string `json:"Email"`
	Department string `json:"Department"`
	Section    string `json:"Section"`
	Position   string `json:"Position"`
	HandSign   string `json:"HandSign"`
	CompanyID  *uint  `json:"CompanyID"`
}
