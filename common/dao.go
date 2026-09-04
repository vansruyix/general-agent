package common

import "gorm.io/gorm"

type Dao struct {
	DB *gorm.DB
}
