package interfaces

import (
	"gorm.io/gorm"
)

type Resources interface {
	Database() *gorm.DB
	Config() Config
}
