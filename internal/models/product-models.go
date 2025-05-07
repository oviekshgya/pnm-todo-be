package models

import "time"

const PRODUCT = "product"

type Product struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Jumlah    int       `gorm:"type:int;not null" json:"jumlah"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return PRODUCT
}
