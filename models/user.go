package models

type User struct {
	Base
	Password   string   `json:"password,omitempty" gorm:"column:password;type:text;not null"`
	UserSalt   UserSalt `json:"-" gorm:"foreignKey:UserSaltID"`
	UserSaltID uint     `json:"-" gorm:"column:user_salt_id;type:bigint;not null"`
	Name       string   `json:"name,omitempty" gorm:"column:name;type:text;not null"`
	Email      string   `json:"email,omitempty" gorm:"column:email;type:text;not null;unique"`
}
