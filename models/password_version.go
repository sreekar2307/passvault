package models

import "database/sql"

type PasswordVersion struct {
	ParentPassword Password       `json:"-" gorm:"foreignKey:PasswordID"`
	PasswordID     uint           `json:"-" gorm:"column:password_id;type:bigint;not null"`
	Password       string         `json:"password,omitempty" gorm:"column:password;type:text;not null"`
	HashedPassword string         `json:"-" gorm:"column:hashed_password;type:text;not null"`
	Username       sql.NullString `json:"username,omitempty" gorm:"column:username;type:text"`
	Notes          sql.NullString `json:"notes,omitempty" gorm:"column:notes;type:text"`
	Website        string         `json:"website,omitempty" gorm:"column:website;type:text;not null"`
	Email          sql.NullString `json:"email,omitempty" gorm:"column:email;type:text"`
	Version        uint64         `json:"versionNumber,omitempty" gorm:"column:version;type:bigint;not null;default:1;check:version >= 1"`
}
