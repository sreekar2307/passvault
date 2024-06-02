package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Password struct {
	Base
	Name           sql.NullString `json:"name,omitempty" gorm:"column:name;type:text;index"`
	UserID         uint           `json:"-" gorm:"column:user_id;type:bigint;not null"`
	User           User           `json:"-" gorm:"foreignKey:UserID"`
	Password       string         `json:"password,omitempty" gorm:"column:password;type:text;not null"`
	HashedPassword string         `json:"-" gorm:"column:hashed_password;type:text;not null"`
	Username       sql.NullString `json:"username,omitempty" gorm:"column:username;type:text"`
	Notes          sql.NullString `json:"notes,omitempty" gorm:"column:notes;type:text"`
	Website        string         `json:"website,omitempty" gorm:"column:website;type:text;not null"`
	Email          sql.NullString `json:"email,omitempty" gorm:"column:email;type:text"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
