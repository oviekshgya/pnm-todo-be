package service

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"pnm-todo-be/internal/models"
	"pnm-todo-be/pkg"
	"time"
)

type UserService struct {
	DB *gorm.DB
}

func (service *UserService) RegisterUser(input pkg.RegisterRequest) (interface{}, error) {
	result, err := pkg.WithTransaction(service.DB, func(tz *gorm.DB) (interface{}, error) {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		if created := models.CreateUser(tz, models.User{
			Password:  string(hash),
			Email:     input.Email,
			Name:      input.Name,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}); created != nil {
			return created, nil
		}
		return map[string]interface{}{
			"email": input.Email,
			"time":  time.Now(),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
