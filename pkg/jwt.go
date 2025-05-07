package pkg

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
)

const SECRETKEY = "pYJLL6mIaWNtx6Q4m0QIcq8Svzuv57Qp"

type Token struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func CreateJWTToken(data DataTokenJWT) (*Token, error) {
	td := &Token{}

	td.AtExpires = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	td.RtExpires = time.Now().Add(time.Hour * time.Duration(48)).Unix()

	atClaims := jwt.MapClaims{
		"id":    data.Id,
		"iat":   time.Now().UTC().Unix(),
		"exp":   td.AtExpires,
		"email": data.Email,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(viper.GetString("SERVICE_SECRET_KEY_ACCESS")))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken

	rtClaims := jwt.MapClaims{
		"id":    data.Id,
		"iat":   time.Now().UTC().Unix(),
		"exp":   td.RtExpires,
		"email": data.Email,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(viper.GetString("SERVICE_SECRET_KEY_REFRESH")))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	return td, nil
}

type DataTokenJWT struct {
	Id    int     `json:"id"`
	Email string  `json:"email"`
	Iat   float64 `json:"iat"`
	Exp   float64 `json:"exp"`
}

func ExtractTokenJWT(tokenString string, c *fiber.Ctx) (*DataTokenJWT, error) {
	if tokenString == "" {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return nil, fiber.ErrBadRequest
		}
		tokenString = authHeader[len(BEARER_SCHEMA):]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(viper.GetString("SERVICE_SECRET_KEY_ACCESS")), nil
	})

	if err != nil {

		fmt.Println("JWT Parsing Error:", err)

		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("token is expired")
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, errors.New("invalid token signature")
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("token has expired")
		}
	}

	// Extract claims
	result := &DataTokenJWT{
		Id:    int(claims["id"].(float64)),
		Iat:   claims["iat"].(float64),
		Email: claims["email"].(string),
		Exp:   claims["exp"].(float64),
	}

	return result, nil
}
