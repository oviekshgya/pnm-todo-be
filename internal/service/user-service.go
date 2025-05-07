package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"pnm-todo-be/internal/models"
	"pnm-todo-be/pkg"
	"time"
)

type UserService struct {
	DB *gorm.DB
}

func (service *UserService) RegisterUser(input pkg.RegisterRequest) (interface{}, error) {
	result, err := pkg.WithTransaction(service.DB, func(tz *gorm.DB) (interface{}, error) {
		getData := models.FindEmail(tz, input.Email)
		if getData != nil {
			return nil, fmt.Errorf("incorrect email already exists")
		}

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

func (service *UserService) LoginUser(input pkg.LoginRequest) (interface{}, error) {
	result, err := pkg.WithTransaction(service.DB, func(tz *gorm.DB) (interface{}, error) {
		reddis := pkg.InitializeRedis()

		defer reddis.Close()
		getData := models.FindEmail(tz, input.Email)
		if getData == nil {
			return nil, fmt.Errorf("incorrect email or password. Please try again")
		}

		var dataAttempt map[string]int
		if get := reddis.GetKey(fmt.Sprintf("%s-login", input.Email), &dataAttempt); get != nil {
			log.Println("reddis.GetKey:", get)
		}

		if dataAttempt["attempt"] > 3 {
			ttl, err := reddis.GetKeyTTL(fmt.Sprintf("%s-login", input.Email))
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("your account has been locked until %v", ttl)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(getData.Password), []byte(input.Password)); err != nil {
			if ser := reddis.SetKey(fmt.Sprintf("%s-login", input.Email), map[string]int{
				"attempt": dataAttempt["attempt"] + 1,
			}, 1*time.Minute); ser != nil {
				return nil, fmt.Errorf("incorrect email or password. Please try again [R2]")
			}
			return nil, errors.New("incorrect password")
		}

		if del := reddis.DeleteKey(fmt.Sprintf("%s-login", input.Email)); del != nil {
			return nil, fmt.Errorf("incorrect email or password. Please try again [R1]")
		}

		tokenGenerate, _ := pkg.CreateJWTToken(pkg.DataTokenJWT{
			Email: input.Email,
			Id:    int(getData.ID),
			Exp:   1,
			Iat:   1,
		})

		return map[string]interface{}{
			"accessToken":  tokenGenerate.AccessToken,
			"refreshToken": tokenGenerate.RefreshToken,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *UserService) CheckEmail(email string) (interface{}, error) {
	getData := models.FindEmail(service.DB, email)
	if getData != nil {
		return nil, fmt.Errorf("incorrect email already exists")
	}
	return map[string]interface{}{
		"emailStatus": true,
	}, nil
}
