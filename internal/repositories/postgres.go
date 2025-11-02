package repositories

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/danubiobwm/company-api/internal/models"
)

type DBConfig struct {
	Host, Port, User, Password, DBName, SSLMode string
}

func NewGormDB(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate for development convenience. Remove in prod.
	if err := db.AutoMigrate(&models.Colaborador{}, &models.Departamento{}); err != nil {
		log.Printf("warning: automigrate error: %v", err)
	}

	return db, nil
}
