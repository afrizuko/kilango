package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	clientDB *gorm.DB
	err      error
)

func NewConnection() *gorm.DB {
	if clientDB != nil {
		return clientDB
	}

	dsn := "admin:S3cr3t123*@tcp(192.168.1.105:3306)/logify?charset=utf8mb4&parseTime=True&loc=Local"
	clientDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return clientDB
}
