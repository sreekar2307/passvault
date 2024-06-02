package resources

import (
	"gorm.io/gorm"
	"passVault/interfaces"
)

var (
	databaseConn *gorm.DB
	config       configImpl
)

func Database() *gorm.DB {
	return databaseConn
}

func Config() interfaces.Config {
	return config
}
