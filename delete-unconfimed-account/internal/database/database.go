package database

import (
	"delete-unconfirmed-account/internal/configuration"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(database configuration.DatabaseConfig) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", database.Host, database.UserName, database.Password, database.DatabaseName, database.DatabasePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(fmt.Errorf("falhou ao inicializar uma sessao com o banco de dados: %s", err.Error()))
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(fmt.Errorf("falhou ao conectar com o banco de dados: %s", err.Error()))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
