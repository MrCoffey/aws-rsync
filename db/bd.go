package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func FindPathInDb(path string) bool {
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local") // TODO: get from environment variable
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if db.Where("name = ?", path).First(Object{}).RecordNotFound() {
		fmt.Println("path not found")
		return false
	}
	return true
}

func UpdateInDb(path string) bool {
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local") // TODO: get from environment variable
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err := db.Model(Object{}).Where("path = ?", path).Update("path", path).Error; err != nil {
		fmt.Println("failed to update object with path: ", path)
		return false
	}
	return true
}

func CreateInDb(path string) bool {
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local") // TODO: get from environment variable
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	object := Object{Path: path}
	if err := db.Create(&object).Error; err != nil {
		fmt.Println("failed to create object with path: ", path)
		return false
	}
	return true
}
