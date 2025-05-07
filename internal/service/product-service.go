package service

import (
	"gorm.io/gorm"
	"pnm-todo-be/internal/models"
	"pnm-todo-be/pkg"
	"time"
)

type ProductService struct {
	DB *gorm.DB
}

func (service *ProductService) CRUDProduct(data pkg.CRUDProduct) (interface{}, error) {
	result, err := pkg.WithTransaction(service.DB, func(tx *gorm.DB) (interface{}, error) {
		dataDB := models.Product{
			Name:      data.Name,
			ID:        data.ID,
			Jumlah:    data.Jumlah,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if data.ID != 0 {
			if updated := models.UpdateProduct(tx, dataDB); updated != nil {
				return nil, updated
			}
			return &dataDB, nil
		}
		created, err := models.CreateProduct(tx, dataDB)
		if err != nil {
			return nil, err
		}
		dataDB.ID = created.ID
		return &dataDB, nil

	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
