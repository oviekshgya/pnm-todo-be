package models

import "gorm.io/gorm"

func CreateUser(tx *gorm.DB, input User) error {
	return tx.Create(&input).Error
}

func FindEmail(tx *gorm.DB, email string) *User {
	var data User
	tx.Where("email = ?", email).First(&data)
	if data.Email != "" {
		return &data
	}
	return nil
}
