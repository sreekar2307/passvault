package models

type PasswordGenerationHistory struct {
	Base
	UserID   uint   `json:"-" gorm:"column:user_id;type:bigint;not null"`
	User     User   `json:"-" gorm:"foreignKey:UserID"`
	Password string `json:"password,omitempty" gorm:"type:text;not null"`
}
