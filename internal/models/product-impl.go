package models

import "gorm.io/gorm"

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

func SelectProductById() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select("name, jumlah, created_at, updated_at").Order("created_at desc")
	}
}

func CreateProduct(db *gorm.DB, data Product) (*Product, error) {
	if err := db.Create(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func UpdateProduct(db *gorm.DB, data Product) error {
	return db.Save(&data).Error
}

func DeleteProductById(db *gorm.DB, id int) error {
	return db.Where("id = ?", id).Delete(&Product{}).Error
}
