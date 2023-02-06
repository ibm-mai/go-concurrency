package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetConnectionString(config Config) string {
	mssqlHost := config.Db.Host
	mssqlUsername := config.Db.Username
	mssqlPassword := config.Db.Password
	mssqlDatabase := config.Db.Database
	mssqlPort := config.Db.Port

	return fmt.Sprintf("sqlserver://%v:%v@%v:%v?database=%s", mssqlUsername, mssqlPassword, mssqlHost, mssqlPort, mssqlDatabase)
}

func getDriverString() string {
	return fmt.Sprint(viper.Get("db.driver"))
}
