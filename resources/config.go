package resources

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type configImpl struct {
	viper *viper.Viper
}

func (c configImpl) GetString(s string) string {
	return c.viper.GetString(s)
}

func (c configImpl) Get(s string) any {
	return c.viper.Get(s)
}

func (c configImpl) GetInt(d string) int {
	return c.viper.GetInt(d)
}

func (c configImpl) GetDuration(s string) time.Duration {
	return c.viper.GetDuration(s)
}

func (c configImpl) UnmarshalKey(s string, rawVal any) error {
	return c.viper.UnmarshalKey(s, &rawVal)
}

func initConfig() error {
	viperObj := viper.New()
	viperObj.AutomaticEnv()
	viperObj.SetConfigName("config")
	viperObj.SetConfigName("secrets")
	viperObj.AddConfigPath(".")
	err := viperObj.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	config = configImpl{viper: viperObj}
	return nil
}
