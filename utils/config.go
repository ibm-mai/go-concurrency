package utils

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Db          Db `mapstructure:"Db"`
	Concurrence int
}

type Db struct {
	Driver   string
	Source   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
