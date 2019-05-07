package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"migrate/config"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql",
		fmt.Sprintf(
			"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
			config.Conf.DBUser, config.Conf.DBPassword, config.Conf.DBName,
		))
	if err != nil {
		panic(err)
	}
}
