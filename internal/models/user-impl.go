package models

import "gorm.io/gorm"

func CreateUser(tx *gorm.DB, input User) error {
	return tx.Create(&input).Error
}
