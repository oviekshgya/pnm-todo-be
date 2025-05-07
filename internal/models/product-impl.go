package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TableNameGet(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(name)
	}
}

func WHEREProductId(id int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func SEARCHProductByName(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name LIKE ?", "%"+search+"%")
	}
}

func SelectProduct() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select("name, jumlah, created_at, updated_at, id").Order("created_at desc")
	}
}

func CRUDProduct(db *gorm.DB, method string, data Product) (*Product, error) {
	switch method {
	case fiber.MethodPost:
		if err := db.Create(&data).Error; err != nil {
			return nil, err
		}
		return &data, nil

	case fiber.MethodPut:
		if err := db.Save(&data).Error; err != nil {
			return nil, err
		}
		return &data, nil

	case fiber.MethodDelete:
		if err := db.Where("id = ?", data.ID).Delete(&Product{}).Error; err != nil {
			return nil, err
		}
		return &data, nil

	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}
