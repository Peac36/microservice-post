package persist

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BootDatabase(username string, password string, address string, port string, databaseName string) *gorm.DB {
	dsnTemplate := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(dsnTemplate, username, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can not connect to database %s \n", err)
	}

	return db
}
