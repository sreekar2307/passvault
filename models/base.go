package models

import "time"

type Base struct {
	ID        uint      `json:"ID" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;type:timestamp;not null"`
}
