package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"pnm-todo-be/internal/models"
	"time"
)

type DatabaseConfig struct {
	Driver       string
	Dbname       string
	Username     string
	Password     string
	Host         string
	Port         string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

var ConnDB *gorm.DB

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Product{})
	if db.Error != nil {
		log.Fatal(db.Error)
	}
}

func (conf DatabaseConfig) ConnectionDataBaseMain() {
	var err error
	var db *gorm.DB

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Dbname,
	)

	for retries := 0; retries < 5; retries++ {
		db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Error connecting to the database: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to the database after retries: %v", err)
	}
	Migrate(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get *sql.DB instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
}
