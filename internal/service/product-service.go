package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pnm-todo-be/internal/models"
	"pnm-todo-be/pkg"
	"strconv"
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

func (service *ProductService) GetProduct(id int, c *fiber.Ctx) (interface{}, error) {
	var data []models.Product
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	if id != 0 {
		service.DB.Scopes(models.TableNameGet(models.PRODUCT), models.WHEREProductId(id), models.SelectProduct()).Find(&data)
		return data, nil
	}

	var count int64
	var totalPage int

	switch {
	case pageSize > 100:
		pageSize = pageSize
	case pageSize <= 0:
		pageSize = 10
	}

	if page == 0 {
		service.DB.Scopes(models.TableNameGet(models.PRODUCT), models.SelectProduct()).Find(&data)
		return &data, nil
	}

	service.DB.Scopes(models.Paginate(pageSize, page), models.TableNameGet(models.PRODUCT), models.SEARCHProductByName(c.Query("search")), models.SelectProduct()).Find(&data)
	service.DB.Scopes(models.Paginate(pageSize, page), models.TableNameGet(models.PRODUCT), models.SEARCHProductByName(c.Query("search"))).Count(&count)

	if count < int64(pageSize) {
		totalPage = 1
	} else {
		totalPage = int(count / int64(pageSize))
		if (count % int64(pageSize)) != 0 {
			totalPage = totalPage + 1
		}
	}

	return map[string]interface{}{
		"data":       data,
		"page":       page,
		"pageSize":   pageSize,
		"total":      count,
		"totalPages": totalPage,
	}, nil

}
