package models

type UserSalt struct {
	Base
	Salt string `json:"salt" gorm:"column:salt;type:text;not null"`
}
