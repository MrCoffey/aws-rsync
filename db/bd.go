package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Ej: "user:pass@0.0.0.0:3306/dbname?charset=utf8&parseTime=True&loc=Local"
var databaseURL string = os.Getenv("DATABASE_URL")

type Object struct {
	gorm.Model
	Path string
}

func FindPathInDb(path string) bool {
	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n")
	}
	defer db.Close()

	fmt.Printf("Looking for path: %s\n", path)

	if db.Where("path = ?", path).Take(Object{}).RecordNotFound() {
		fmt.Println("path not found: \n", path)
		return false
	}
	return true
}

func UpdateInDb(origin string, destination string) bool {
	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	fmt.Printf("Updating path %s to: %s\n", origin, destination)

	if err := db.Model(Object{}).Where("path = ?", origin).Update("path", destination).Error; err != nil {
		fmt.Println("failed to update object with path: \n", origin)
		return false
	}
	return true
}

func CreateInDb(path string) bool {
	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	fmt.Printf("Inserting path %s:\n", path)

	object := Object{Path: path}
	if err := db.Create(&object).Error; err != nil {
		fmt.Println("failed to create object with path: \n", path)
		return false
	}

	return true
}
