package data

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func openMySql(server, database, username, password string, port int) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func InitDatabase(file, server, database, username, password string, port int, width int, height int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}

	DB.AutoMigrate(&Pixel{})

	var nrOfPixels int64

	DB.Model(&Pixel{}).Count(&nrOfPixels)

	if nrOfPixels == 0 {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				DB.Create(&Pixel{X: x, Y: y, Color: "FFFFFF"})
			}
		}
	}
}
