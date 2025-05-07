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
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Dbname,
	)
	var err error

	ConnDB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	Migrate(ConnDB)

	sqlDB, err := ConnDB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan instance *sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)

}
