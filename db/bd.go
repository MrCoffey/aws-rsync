package db

import (
	"fmt"

	"github.com/MrCoffey/s3-resync/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Object is a S3 Model in database
type Object struct {
	gorm.Model
	Path string
}

// FindPathInDb lookup for a provided path
func FindPathInDb(config *config.Values, path string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	fmt.Printf("Looking for path: %s\n\n", path)

	if db.Where("path = ?", path).Take(Object{}).RecordNotFound() {
		fmt.Println("path not found: \n\n", path)
		return false
	}
	return true
}

// UpdateInDb Updates a provided path
func UpdateInDb(config *config.Values, origin string, destination string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	fmt.Printf("Updating path %s to: %s \n\n", origin, destination)

	if err := db.Model(Object{}).Where("path = ?", origin).Update("path", destination).Error; err != nil {
		fmt.Println("failed to update object with path: \n", origin)
		return false
	}
	return true
}

// CreateInDb Inserts a provided path
func CreateInDb(config *config.Values, path string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n\n")
	}
	defer db.Close()

	fmt.Printf("Inserting path %s:\n\n", path)

	object := Object{Path: path}
	if err := db.Create(&object).Error; err != nil {
		fmt.Println("failed to create object with path: \n\n", path)
		return false
	}

	return true
}
