package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pnm-todo-be/internal/models"
	"pnm-todo-be/pkg"
	"time"
)

type ProductService struct {
	DB *gorm.DB
}

func (service *ProductService) CRUDProduct(data pkg.CRUDProduct, c *fiber.Ctx) (interface{}, error) {
	result, err := pkg.WithTransaction(service.DB, func(tx *gorm.DB) (interface{}, error) {
		dataDB := models.Product{
			Name:      data.Name,
			ID:        data.ID,
			Jumlah:    data.Jumlah,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		created, err := models.CRUDProduct(tx, c.Method(), dataDB)
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
