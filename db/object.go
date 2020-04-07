package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Object struct {
	gorm.Model
	Path         string
	LastModified *time.Time
}
