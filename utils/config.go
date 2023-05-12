package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Db      Db      `mapstructure:"Db"`
	GenData GenData `mapstructure:"genData"`
}

type GenData struct {
	Concurrence         int
	InputPath           string
	MsisdnDuplicatePath string
	CommitSize          int
	LogLevel            string
}

type Db struct {
	Driver                   string
	Source                   string
	Host                     string
	Port                     string
	Username                 string
	Password                 string
	Database                 string
	CustomerProfileTableName string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadInConfig()
	if err != nil {
		logrus.Error("Error reading config")
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Error("Error unmarshal config")
		return
	}
	return
}
