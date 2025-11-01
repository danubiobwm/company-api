package repositories

import (
	"fmt"

	"github.com/danubiobwm/company-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host, Port, User, Password, DBName, SSLMode string
}

func NewGormDB(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// auto-migrate only for development convenience
	db.AutoMigrate(&models.Colaborador{}, &models.Departamento{})

	return db, nil
}
