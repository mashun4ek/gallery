package models

import (
	"github.com/jinzhu/gorm"
)

type Gallary struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
