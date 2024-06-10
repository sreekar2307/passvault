package resources

import (
	"sync"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		if err := initConfig(); err != nil {
			panic(err)
		}
		if err := initDatabaseConn(); err != nil {
			panic(err)
		}
		if err := initLogger(); err != nil {
			panic(err)
		}
	})
}
