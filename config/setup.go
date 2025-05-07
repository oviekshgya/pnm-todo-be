package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"pnm-todo-be/api/routes"
	"pnm-todo-be/db"
)

type SetupDatabase struct {
	DBMain *gorm.DB
}

func SetDatabase() *SetupDatabase {
	configDb := db.DatabaseConfig{
		MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
		Host:         viper.GetString("DB_HOST"),
		Password:     viper.GetString("DB_PASSWORD"),
		Username:     viper.GetString("DB_USER"),
		Driver:       viper.GetString("DB_DRIVER"),
		Port:         viper.GetString("DB_PORT"),
		Dbname:       viper.GetString("DB_NAME"),
		MaxLifetime:  viper.GetInt("DB_MAX_LIFE_TIME"),
		MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNECTIONS"),
	}

	configDb.ConnectionDataBaseMain()

	return &SetupDatabase{
		DBMain: db.ConnDB,
	}
}

func Start() {
	viper.SetConfigFile("./.env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error saat membaca file .env: %v", err)
	}

	SetDatabase()
	routes.Router = fiber.New(fiber.Config{
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         false,
		ServerHeader:          "Fiber",
		AppName:               "Face service",
		DisableStartupMessage: true,
	})

	routes.InitialRoute()
	routes.Route()

	if !fiber.IsChild() {
		log.Printf("INFO: SERVICE RUNNING ON PORT " + viper.GetString("SERVICE_PORT"))
	}

	err2 := routes.Router.Listen(":" + viper.GetString("SERVICE_PORT"))
	if err2 != nil {
		log.Fatalf("ERROR: cannot start server on port " + viper.GetString("SERVICE_PORT"))
		return
	}
}
