package db

import (
	"fmt"

	"github.com/MrCoffey/s3-sync/config"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Object is a S3 Model in database
type Object struct {
	gorm.Model
	ID int `gorm:"AUTO_INCREMENT"`

	Path   string
	Bucket string
}

func MigrateDB(config *config.Values) error {
	db, err := initDB(config)
	if err != nil {
		return err
	}

	db.AutoMigrate(&Object{})
	return nil
}

// ListPaths all the records
func ListPaths(config *config.Values) {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	count := 0
	objects := db.Table("objects").Select("count(distinct(path))").Count(&count)

	fmt.Printf("%v \n\n", objects)
}

// FindPathInDb lookup for a provided path
func FindPathInDb(config *config.Values, bucket string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	if db.Where("bucket = ? AND path = ?", bucket, key).Find(&Object{}).RecordNotFound() {
		fmt.Printf("Failed to find KEY: %s %s \n", bucket, key)
		return false
	}

	return true
}

// UpdateInDb Updates a provided path
func UpdateInDb(config *config.Values, legacyBucket string, destinationBucket string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	fmt.Printf("Updating path %s to: %s/%s \n\n", legacyBucket, destinationBucket, key)

	object := Object{}

	if err := db.Model(&object).Where("bucket = ? AND path = ?", legacyBucket, key).Updates(Object{Bucket: destinationBucket, Path: key}).RecordNotFound(); err {
		fmt.Println("failed to update object with path: \n", legacyBucket)
		return false
	}
	return true
}

// CreateInDb Inserts a provided path
func CreateInDb(config *config.Values, bucketName string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n\n")
	}
	defer db.Close()

	fmt.Printf("Inserting path %s/%s:\n\n", bucketName, key)

	object := Object{Path: key, Bucket: bucketName}
	if err := db.Create(&object).Error; err != nil {
		fmt.Println("failed to create object with path: \n\n", key)
		return false
	}

	return true
}

func BulkCreateRecords(config *config.Values, objects []interface{}) error {
	db, err := initDB(config)
	if err != nil {
		return err
	}

	err = gormbulk.BulkInsert(db, objects, 1000)
	if err != nil {
		return err
	}
	return nil
}

func initDB(config *config.Values) (*gorm.DB, error) {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n\n")
		return db, err
	}
	defer db.Close()

	return db, nil
}
